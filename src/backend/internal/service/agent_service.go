package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xinewang/oen/internal/model"
)

// Valid agent states
const (
	AgentUnauthorized = "unauthorized"
	AgentAuthorized   = "authorized"
	AgentCreating     = "creating"
	AgentRunning      = "running"
	AgentDegraded     = "degraded"
	AgentPaused       = "paused"
	AgentRevoked      = "revoked"
)

// Valid state transitions: from → {allowed to states}
var validTransitions = map[string][]string{
	AgentUnauthorized: {AgentAuthorized},
	AgentAuthorized:   {AgentCreating, AgentRevoked},
	AgentCreating:     {AgentRunning, AgentDegraded},
	AgentRunning:      {AgentDegraded, AgentPaused, AgentRevoked, AgentRunning}, // running→running for heartbeat refresh
	AgentDegraded:     {AgentRunning, AgentPaused}, // rebuild → running
	AgentPaused:       {AgentRunning},
	AgentRevoked:      {}, // terminal state
}

// ErrInvalidTransition is returned when a state transition is not allowed
var ErrInvalidTransition = errors.New("invalid state transition")

// TransitionAgent validates and executes a state transition
func (s *Service) TransitionAgent(ctx context.Context, id uint, from, to string) error {
	allowed, ok := validTransitions[from]
	if !ok {
		return fmt.Errorf("unknown state: %s", from)
	}

	found := false
	for _, a := range allowed {
		if a == to {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("%w: %s → %s", ErrInvalidTransition, from, to)
	}

	// Optimistic state check via WHERE state = from
	agent, err := s.Repo.GetAgentByID(ctx, id)
	if err != nil {
		return err
	}
	if agent.State != from {
		return fmt.Errorf("state mismatch: expected %s, got %s", from, agent.State)
	}

	return s.Repo.UpdateAgentState(ctx, id, from, to)
}

// Consent handles agent authorization
func (s *Service) Consent(ctx context.Context, agentID uint, consentType, grantedBy string) (*model.ConsentRecord, error) {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	// Transition: unauthorized → authorized
	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentAuthorized); err != nil {
		return nil, fmt.Errorf("consent failed: %w", err)
	}

	consent := &model.ConsentRecord{
		AgentID:     agentID,
		ConsentType: consentType,
		Status:      "active",
		GrantedBy:   grantedBy,
	}
	if err := s.Repo.CreateConsent(ctx, consent); err != nil {
		return nil, err
	}

	// Audit
	s.LogEvent(ctx, "consent", "granted", "user", grantedBy, "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Agent %s authorized, type=%s", agent.Name, consentType))

	return consent, nil
}

// RevokeConsent handles authorization revocation
func (s *Service) RevokeConsent(ctx context.Context, agentID uint, revokedBy string) error {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return err
	}

	// Transition to revoked (from authorized, running, creating, degraded, or paused)
	validFrom := []string{AgentAuthorized, AgentCreating, AgentRunning, AgentDegraded, AgentPaused}
	validFromState := false
	for _, s := range validFrom {
		if agent.State == s {
			validFromState = true
			break
		}
	}
	if !validFromState {
		return fmt.Errorf("cannot revoke consent from state: %s", agent.State)
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentRevoked); err != nil {
		return fmt.Errorf("revoke failed: %w", err)
	}

	consent, err := s.Repo.GetActiveConsent(ctx, agentID)
	if err == nil && consent != nil {
		s.Repo.RevokeConsent(ctx, consent.ID)
	}

	s.LogEvent(ctx, "consent", "revoked", "user", revokedBy, "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Agent %s consent revoked", agent.Name))

	return nil
}

// CreateSubAgent transitions authorized → creating (sub-agent creation)
func (s *Service) CreateSubAgent(ctx context.Context, agentID uint) (*model.Agent, error) {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return nil, err
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentCreating); err != nil {
		return nil, fmt.Errorf("cannot create sub-agent: %w", err)
	}

	s.LogEvent(ctx, "agent", "creating", "system", "oen", "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Sub-agent creation triggered for %s", agent.Name))

	return agent, nil
}

// CompleteSubAgentCreation transitions creating → running
func (s *Service) CompleteSubAgentCreation(ctx context.Context, agentID uint, routeMode string) error {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return err
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentRunning); err != nil {
		return fmt.Errorf("cannot complete sub-agent creation: %w", err)
	}

	if routeMode != "" {
		s.Repo.UpdateAgentRouteMode(ctx, agentID, routeMode)
	}

	s.LogEvent(ctx, "agent", "running", "system", "oen", "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Sub-agent %s is now running, route_mode=%s", agent.Name, routeMode))

	return nil
}

// RecordHeartbeat records a heartbeat and auto-triggers degraded detection
func (s *Service) RecordHeartbeat(ctx context.Context, agentID uint, status, routeMode string, cpuUsage, memoryUsage float64, errMsg string) error {
	now := time.Now()
	hb := &model.AgentHeartbeat{
		AgentID:      agentID,
		HeartbeatAt:  now,
		Status:       status,
		RouteMode:    routeMode,
		CPUUsage:     cpuUsage,
		MemoryUsage:  memoryUsage,
		ErrorMessage: errMsg,
	}

	if err := s.Repo.RecordHeartbeat(ctx, hb); err != nil {
		return err
	}

	// Auto-degrade if heartbeat reports error
	if errMsg != "" {
		agent, err := s.Repo.GetAgentByID(ctx, agentID)
		if err == nil && agent.State == AgentRunning {
			s.Repo.UpdateAgentState(ctx, agentID, AgentRunning, AgentDegraded)
			s.LogEvent(ctx, "agent", "degraded", "system", "oen", "agent", fmt.Sprintf("%d", agentID),
				fmt.Sprintf("Agent %s auto-degraded: %s", agent.Name, errMsg))
		}
	}

	return nil
}

// GetRecentHeartbeats returns recent heartbeat history for an agent
func (s *Service) GetRecentHeartbeats(ctx context.Context, agentID uint, limit int) ([]model.AgentHeartbeat, error) {
	if limit <= 0 {
		limit = 20
	}
	return s.Repo.GetAgentHeartbeats(ctx, agentID, limit)
}

// RebuildAgent triggers agent rebuild (degraded/paused → running)
func (s *Service) RebuildAgent(ctx context.Context, agentID uint, triggeredBy string) error {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return err
	}

	validFrom := []string{AgentDegraded, AgentPaused}
	validFromState := false
	for _, s := range validFrom {
		if agent.State == s {
			validFromState = true
			break
		}
	}
	if !validFromState {
		return fmt.Errorf("can only rebuild from degraded/paused, current state: %s", agent.State)
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentRunning); err != nil {
		return err
	}

	s.LogEvent(ctx, "agent", "rebuild", "user", triggeredBy, "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Agent %s rebuilt from %s", agent.Name, agent.State))

	return nil
}

// PauseAgent pauses a running agent
func (s *Service) PauseAgent(ctx context.Context, agentID uint, triggeredBy string) error {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return err
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentPaused); err != nil {
		return err
	}

	s.LogEvent(ctx, "agent", "paused", "user", triggeredBy, "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Agent %s paused", agent.Name))

	return nil
}

// ResumeAgent resumes a paused agent
func (s *Service) ResumeAgent(ctx context.Context, agentID uint, triggeredBy string) error {
	agent, err := s.Repo.GetAgentByID(ctx, agentID)
	if err != nil {
		return err
	}

	if err := s.TransitionAgent(ctx, agentID, agent.State, AgentRunning); err != nil {
		return err
	}

	s.LogEvent(ctx, "agent", "resumed", "user", triggeredBy, "agent", fmt.Sprintf("%d", agentID),
		fmt.Sprintf("Agent %s resumed", agent.Name))

	return nil
}
