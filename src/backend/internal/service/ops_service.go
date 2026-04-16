package service

import (
	"context"
	"time"

	"github.com/xinewang/oen/internal/model"
)

// LogEvent writes an entry to the audit log
func (s *Service) LogEvent(ctx context.Context, eventType, eventSubtype, actorType, actorID, targetType, targetID, description string) {
	log := &model.AuditLog{
		EventType:    eventType,
		EventSubtype: eventSubtype,
		ActorType:    actorType,
		ActorID:      actorID,
		TargetType:   targetType,
		TargetID:     targetID,
		Description:  description,
		CreatedAt:    time.Now(),
	}
	s.Repo.CreateAuditLog(ctx, log)
}

// GetStatusOverview returns a system-wide status overview
func (s *Service) GetStatusOverview(ctx context.Context) (map[string]interface{}, error) {
	// Agent counts by state
	var stateCounts []struct {
		State string
		Count int64
	}
	s.Repo.DB.WithContext(ctx).
		Model(&model.Agent{}).
		Select("state, count(*) as count").
		Group("state").
		Find(&stateCounts)

	var artifactCount, recCount, candidateCount int64
	s.Repo.DB.WithContext(ctx).Model(&model.Artifact{}).Count(&artifactCount)
	s.Repo.DB.WithContext(ctx).Model(&model.Recommendation{}).Where("state = ?", "pending").Count(&recCount)
	s.Repo.DB.WithContext(ctx).Model(&model.CandidateResource{}).Where("state = ?", "pending").Count(&candidateCount)

	return map[string]interface{}{
		"agent_states":       stateCounts,
		"total_artifacts":    artifactCount,
		"pending_recommendations": recCount,
		"pending_candidates": candidateCount,
		"version":            "0.1.0",
		"status":             "running",
	}, nil
}
