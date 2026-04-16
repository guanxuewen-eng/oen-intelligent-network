package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xinewang/oen/internal/model"
)

// CreateCandidate creates a candidate resource (from online search or local learning)
func (s *Service) CreateCandidate(ctx context.Context, candidate *model.CandidateResource) error {
	if candidate.State == "" {
		candidate.State = "pending"
	}
	if err := s.Repo.CreateCandidate(ctx, candidate); err != nil {
		return err
	}

	s.LogEvent(ctx, "candidate", "created", "system", "oen", "candidate", fmt.Sprintf("%d", candidate.ID),
		fmt.Sprintf("Candidate resource created: %s", candidate.Title))

	return nil
}

// ListCandidates returns paginated candidates
func (s *Service) ListCandidates(ctx context.Context, state, sourceType string, page, pageSize int) ([]model.CandidateResource, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.Repo.ListCandidates(ctx, state, sourceType, page, pageSize)
}

// GetCandidate returns a candidate by ID
func (s *Service) GetCandidate(ctx context.Context, id uint) (*model.CandidateResource, error) {
	return s.Repo.GetCandidate(ctx, id)
}

// ReviewCandidate reviews a candidate and optionally promotes it to an artifact
func (s *Service) ReviewCandidate(ctx context.Context, candidateID uint, state, reviewNotes, reviewedBy string, promoteToArtifact *model.Artifact) error {
	candidate, err := s.Repo.GetCandidate(ctx, candidateID)
	if err != nil {
		return err
	}

	if promoteToArtifact != nil {
		return s.Repo.PromoteCandidateToArtifact(ctx, candidate, promoteToArtifact)
	}

	return s.Repo.ReviewCandidate(ctx, candidateID, state, reviewNotes, nil)
}

// CreateRecommendation creates a recommendation for an agent
func (s *Service) CreateRecommendation(ctx context.Context, rec *model.Recommendation) error {
	if rec.State == "" {
		rec.State = "pending"
	}
	if err := s.Repo.CreateRecommendation(ctx, rec); err != nil {
		return err
	}

	s.LogEvent(ctx, "recommendation", "created", "system", "oen", "recommendation", fmt.Sprintf("%d", rec.ID),
		fmt.Sprintf("Recommendation created for agent %d: %s (confidence=%.2f)", rec.AgentID, rec.Title, rec.ConfidenceScore))

	return nil
}

// ListRecommendations returns paginated recommendations
func (s *Service) ListRecommendations(ctx context.Context, agentID uint, state string, page, pageSize int) ([]model.Recommendation, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return s.Repo.ListRecommendations(ctx, agentID, state, page, pageSize)
}

// GetRecommendation returns a recommendation by ID
func (s *Service) GetRecommendation(ctx context.Context, id uint) (*model.Recommendation, error) {
	return s.Repo.GetRecommendation(ctx, id)
}

// MakeDecision records a user decision on a recommendation
func (s *Service) MakeDecision(ctx context.Context, recID uint, decision, decidedBy, note string) error {
	rec, err := s.Repo.GetRecommendation(ctx, recID)
	if err != nil {
		return err
	}
	if rec.State != "pending" {
		return fmt.Errorf("recommendation already decided (state=%s)", rec.State)
	}

	// Map decision to recommendation state
	recState := "ignored"
	switch decision {
	case "accept":
		recState = "accepted"
	case "ignore":
		recState = "ignored"
	case "later":
		recState = "deferred"
	default:
		return fmt.Errorf("invalid decision: %s (must be accept/ignore/later)", decision)
	}

	d := &model.RecommendationDecision{
		RecommendationID: recID,
		Decision:         decision,
		DecidedBy:        decidedBy,
		DecidedAt:        time.Now(),
		Note:             note,
	}

	if err := s.Repo.MakeDecision(ctx, d, recID, recState); err != nil {
		return err
	}

	s.LogEvent(ctx, "recommendation", "decision", "user", decidedBy, "recommendation", fmt.Sprintf("%d", recID),
		fmt.Sprintf("Decision made on recommendation: %s (note: %s)", decision, note))

	return nil
}

// GenerateRecommendationsFromCandidates creates recommendations from approved candidates
func (s *Service) GenerateRecommendationsFromCandidates(ctx context.Context, artifactID uint, agentID uint) error {
	artifact, err := s.Repo.GetArtifactByID(ctx, artifactID)
	if err != nil {
		return err
	}

	rec := &model.Recommendation{
		RecommendationKey: fmt.Sprintf("rec_artifact_%d_agent_%d", artifactID, agentID),
		AgentID:           agentID,
		ArtifactID:        artifactID,
		Title:             artifact.Title,
		MatchReason:       fmt.Sprintf("Artifact %s matches agent %d's context", artifact.Title, agentID),
		RiskSummary:       fmt.Sprintf("Risk level: %s", artifact.RiskLevel),
		ConfidenceScore:   0.75, // base confidence, can be tuned
		State:             "pending",
	}

	return s.CreateRecommendation(ctx, rec)
}
