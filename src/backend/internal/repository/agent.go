package repository

import (
	"context"
	"time"

	"github.com/xinewang/oen/internal/model"
	"gorm.io/gorm"
)

// --- Agent CRUD ---

// CreateAgent creates a new agent record
func (r *Repository) CreateAgent(ctx context.Context, agent *model.Agent) error {
	return r.DB.WithContext(ctx).Create(agent).Error
}

// GetAgentByID retrieves an agent by ID
func (r *Repository) GetAgentByID(ctx context.Context, id uint) (*model.Agent, error) {
	var agent model.Agent
	err := r.DB.WithContext(ctx).First(&agent, id).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// GetAgentByKey retrieves an agent by agent_key
func (r *Repository) GetAgentByKey(ctx context.Context, key string) (*model.Agent, error) {
	var agent model.Agent
	err := r.DB.WithContext(ctx).Where("agent_key = ?", key).First(&agent).Error
	if err != nil {
		return nil, err
	}
	return &agent, nil
}

// ListAgents returns paginated agents with optional filters
func (r *Repository) ListAgents(ctx context.Context, agentType, state string, offset, limit int) ([]model.Agent, int64, error) {
	var agents []model.Agent
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Agent{})

	if agentType != "" {
		query = query.Where("agent_type = ?", agentType)
	}
	if state != "" {
		query = query.Where("state = ?", state)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&agents).Error; err != nil {
		return nil, 0, err
	}

	return agents, total, nil
}

// UpdateAgentState updates only the state field with optimistic check
func (r *Repository) UpdateAgentState(ctx context.Context, id uint, oldState, newState string) error {
	return r.DB.WithContext(ctx).
		Model(&model.Agent{}).
		Where("id = ? AND state = ?", id, oldState).
		Update("state", newState).Error
}

// UpdateAgent updates modifiable fields of an agent
func (r *Repository) UpdateAgent(ctx context.Context, agent *model.Agent) error {
	return r.DB.WithContext(ctx).Save(agent).Error
}

// UpdateAgentRouteMode updates the route_mode of an agent
func (r *Repository) UpdateAgentRouteMode(ctx context.Context, id uint, routeMode string) error {
	return r.DB.WithContext(ctx).
		Model(&model.Agent{}).
		Where("id = ?", id).
		Update("route_mode", routeMode).Error
}

// DeleteAgent removes an agent record
func (r *Repository) DeleteAgent(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&model.Agent{}, id).Error
}

// --- Consent ---

// CreateConsent creates a consent record
func (r *Repository) CreateConsent(ctx context.Context, consent *model.ConsentRecord) error {
	consent.GrantedAt = time.Now()
	return r.DB.WithContext(ctx).Create(consent).Error
}

// GetActiveConsent returns the active consent for an agent
func (r *Repository) GetActiveConsent(ctx context.Context, agentID uint) (*model.ConsentRecord, error) {
	var consent model.ConsentRecord
	err := r.DB.WithContext(ctx).
		Where("agent_id = ? AND status = ?", agentID, "active").
		Order("granted_at DESC").
		First(&consent).Error
	if err != nil {
		return nil, err
	}
	return &consent, nil
}

// RevokeConsent revokes an active consent
func (r *Repository) RevokeConsent(ctx context.Context, consentID uint) error {
	now := time.Now()
	return r.DB.WithContext(ctx).
		Model(&model.ConsentRecord{}).
		Where("id = ? AND status = ?", consentID, "active").
		Updates(map[string]interface{}{
			"status":     "revoked",
			"revoked_at": now,
		}).Error
}

// --- Artifact Versions & Views ---

// CreateArtifactVersion creates a new artifact version
func (r *Repository) CreateArtifactVersion(ctx context.Context, ver *model.ArtifactVersion) error {
	return r.DB.WithContext(ctx).Create(ver).Error
}

// GetArtifactVersions returns versions for an artifact
func (r *Repository) GetArtifactVersions(ctx context.Context, artifactID uint) ([]model.ArtifactVersion, error) {
	var versions []model.ArtifactVersion
	err := r.DB.WithContext(ctx).
		Where("artifact_id = ?", artifactID).
		Order("created_at DESC").
		Find(&versions).Error
	return versions, err
}

// GetArtifactVersionByNumber returns a specific version of an artifact
func (r *Repository) GetArtifactVersionByNumber(ctx context.Context, artifactID uint, versionNumber string) (*model.ArtifactVersion, error) {
	var ver model.ArtifactVersion
	err := r.DB.WithContext(ctx).
		Where("artifact_id = ? AND version_number = ?", artifactID, versionNumber).
		First(&ver).Error
	if err != nil {
		return nil, err
	}
	return &ver, nil
}

// CreateArtifactView creates an artifact view
func (r *Repository) CreateArtifactView(ctx context.Context, view *model.ArtifactView) error {
	return r.DB.WithContext(ctx).Create(view).Error
}

// GetArtifactView returns a specific view for an artifact version
func (r *Repository) GetArtifactView(ctx context.Context, artifactVersionID uint, viewType string) (*model.ArtifactView, error) {
	var view model.ArtifactView
	err := r.DB.WithContext(ctx).
		Where("artifact_version_id = ? AND view_type = ?", artifactVersionID, viewType).
		First(&view).Error
	if err != nil {
		return nil, err
	}
	return &view, nil
}

// GetArtifactVersionsWithViews returns artifact versions with their views
func (r *Repository) GetArtifactVersionsWithViews(ctx context.Context, artifactID uint) ([]model.ArtifactVersion, []model.ArtifactView, error) {
	var versions []model.ArtifactVersion
	var views []model.ArtifactView

	if err := r.DB.WithContext(ctx).
		Where("artifact_id = ?", artifactID).
		Order("created_at DESC").
		Find(&versions).Error; err != nil {
		return nil, nil, err
	}

	if len(versions) == 0 {
		return versions, views, nil
	}

	versionIDs := make([]uint, len(versions))
	for i, v := range versions {
		versionIDs[i] = v.ID
	}

	if err := r.DB.WithContext(ctx).
		Where("artifact_version_id IN ?", versionIDs).
		Find(&views).Error; err != nil {
		return nil, nil, err
	}

	return versions, views, nil
}

// --- Candidate Resource ---

// CreateCandidate creates a candidate resource
func (r *Repository) CreateCandidate(ctx context.Context, candidate *model.CandidateResource) error {
	return r.DB.WithContext(ctx).Create(candidate).Error
}

// GetCandidate returns a candidate by ID
func (r *Repository) GetCandidate(ctx context.Context, id uint) (*model.CandidateResource, error) {
	var c model.CandidateResource
	err := r.DB.WithContext(ctx).First(&c, id).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// ListCandidates returns paginated candidates
func (r *Repository) ListCandidates(ctx context.Context, state, sourceType string, page, pageSize int) ([]model.CandidateResource, int64, error) {
	var candidates []model.CandidateResource
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.CandidateResource{})
	if state != "" {
		query = query.Where("state = ?", state)
	}
	if sourceType != "" {
		query = query.Where("source_type = ?", sourceType)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&candidates).Error; err != nil {
		return nil, 0, err
	}
	return candidates, total, nil
}

// ReviewCandidate updates candidate state after review
func (r *Repository) ReviewCandidate(ctx context.Context, id uint, state, reviewNotes string, artifactID *uint) error {
	now := time.Now()
	updates := map[string]interface{}{
		"state":         state,
		"review_notes":  reviewNotes,
		"reviewed_at":   now,
	}
	if artifactID != nil {
		updates["artifact_id"] = *artifactID
	}
	return r.DB.WithContext(ctx).
		Model(&model.CandidateResource{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// PromoteCandidateToArtifact converts a candidate to an artifact
func (r *Repository) PromoteCandidateToArtifact(ctx context.Context, candidate *model.CandidateResource, artifact *model.Artifact) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(artifact).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&model.CandidateResource{}).
			Where("id = ?", candidate.ID).
			Updates(map[string]interface{}{
				"state":       "promoted",
				"artifact_id": artifact.ID,
				"reviewed_at": now,
			}).Error
	})
}

// --- Recommendation ---

// CreateRecommendation creates a recommendation
func (r *Repository) CreateRecommendation(ctx context.Context, rec *model.Recommendation) error {
	return r.DB.WithContext(ctx).Create(rec).Error
}

// ListRecommendations returns paginated recommendations for an agent
func (r *Repository) ListRecommendations(ctx context.Context, agentID uint, state string, page, pageSize int) ([]model.Recommendation, int64, error) {
	var recs []model.Recommendation
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Recommendation{})
	if agentID > 0 {
		query = query.Where("agent_id = ?", agentID)
	}
	if state != "" {
		query = query.Where("state = ?", state)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&recs).Error; err != nil {
		return nil, 0, err
	}
	return recs, total, nil
}

// GetRecommendation returns a recommendation by ID
func (r *Repository) GetRecommendation(ctx context.Context, id uint) (*model.Recommendation, error) {
	var rec model.Recommendation
	err := r.DB.WithContext(ctx).First(&rec, id).Error
	if err != nil {
		return nil, err
	}
	return &rec, nil
}

// MakeDecision records a decision on a recommendation
func (r *Repository) MakeDecision(ctx context.Context, decision *model.RecommendationDecision, recID uint, newState string) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(decision).Error; err != nil {
			return err
		}
		return tx.Model(&model.Recommendation{}).
			Where("id = ?", recID).
			Updates(map[string]interface{}{
				"state":      newState,
				"decided_at": time.Now(),
			}).Error
	})
}

// --- Agent Heartbeat ---

// RecordHeartbeat inserts a heartbeat record and updates agent's last_heartbeat
func (r *Repository) RecordHeartbeat(ctx context.Context, hb *model.AgentHeartbeat) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(hb).Error; err != nil {
			return err
		}
		now := time.Now()
		return tx.Model(&model.Agent{}).
			Where("id = ?", hb.AgentID).
			Update("last_heartbeat", now).Error
	})
}

// GetAgentHeartbeats returns recent heartbeats for an agent
func (r *Repository) GetAgentHeartbeats(ctx context.Context, agentID uint, limit int) ([]model.AgentHeartbeat, error) {
	var heartbeats []model.AgentHeartbeat
	err := r.DB.WithContext(ctx).
		Where("agent_id = ?", agentID).
		Order("heartbeat_at DESC").
		Limit(limit).
		Find(&heartbeats).Error
	return heartbeats, err
}

// --- Audit Log ---

// CreateAuditLog inserts an audit log record
func (r *Repository) CreateAuditLog(ctx context.Context, log *model.AuditLog) error {
	log.CreatedAt = time.Now()
	return r.DB.WithContext(ctx).Create(log).Error
}

// GetAuditLogs returns paginated audit logs with optional filters
func (r *Repository) GetAuditLogs(ctx context.Context, eventType, targetType, targetID string, offset, limit int) ([]model.AuditLog, int64, error) {
	var logs []model.AuditLog
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.AuditLog{})

	if eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if targetType != "" {
		query = query.Where("target_type = ?", targetType)
	}
	if targetID != "" {
		query = query.Where("target_id = ?", targetID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}
