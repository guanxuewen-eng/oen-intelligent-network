package service

import (
	"context"

	"github.com/xinewang/oen/internal/config"
	"github.com/xinewang/oen/internal/model"
	"github.com/xinewang/oen/internal/repository"
	"gorm.io/gorm"
)

type Service struct {
	DB   *gorm.DB
	Repo *repository.Repository
	Cfg  *config.Config
}

func New(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{
		DB:   db,
		Repo: repository.New(db),
		Cfg:  cfg,
	}
}

// --- Agent Service ---

// CreateAgent creates a new agent
func (s *Service) CreateAgent(ctx context.Context, agent *model.Agent) error {
	return s.Repo.CreateAgent(ctx, agent)
}

// GetAgentByID retrieves an agent by ID
func (s *Service) GetAgentByID(ctx context.Context, id uint) (*model.Agent, error) {
	return s.Repo.GetAgentByID(ctx, id)
}

// ListAgents returns paginated agents
func (s *Service) ListAgents(ctx context.Context, agentType, state string, page, pageSize int) ([]model.Agent, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return s.Repo.ListAgents(ctx, agentType, state, offset, pageSize)
}

// UpdateAgent updates an existing agent
func (s *Service) UpdateAgent(ctx context.Context, agent *model.Agent) error {
	return s.Repo.UpdateAgent(ctx, agent)
}

// DeleteAgent deletes an agent by ID
func (s *Service) DeleteAgent(ctx context.Context, id uint) error {
	return s.Repo.DeleteAgent(ctx, id)
}

// --- Artifact Service ---

// CreateArtifact creates a new artifact
func (s *Service) CreateArtifact(ctx context.Context, artifact *model.Artifact) error {
	return s.Repo.CreateArtifact(ctx, artifact)
}

// GetArtifactByID retrieves an artifact by ID
func (s *Service) GetArtifactByID(ctx context.Context, id uint) (*model.Artifact, error) {
	return s.Repo.GetArtifactByID(ctx, id)
}

// ListArtifacts returns paginated artifacts
func (s *Service) ListArtifacts(ctx context.Context, artifactType, riskLevel, verificationStatus string, page, pageSize int) ([]model.Artifact, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return s.Repo.ListArtifacts(ctx, artifactType, riskLevel, verificationStatus, offset, pageSize)
}

// UpdateArtifact updates an existing artifact
func (s *Service) UpdateArtifact(ctx context.Context, artifact *model.Artifact) error {
	return s.Repo.UpdateArtifact(ctx, artifact)
}

// DeleteArtifact deletes an artifact by ID
func (s *Service) DeleteArtifact(ctx context.Context, id uint) error {
	return s.Repo.DeleteArtifact(ctx, id)
}

// --- Ops Service ---

// GetAuditLogs returns paginated audit logs
func (s *Service) GetAuditLogs(ctx context.Context, eventType, targetType, targetID string, page, pageSize int) ([]model.AuditLog, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize
	return s.Repo.GetAuditLogs(ctx, eventType, targetType, targetID, offset, pageSize)
}

// CountAgents returns the total count of agents
func (s *Service) CountAgents(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).Model(&model.Agent{}).Count(&count).Error
	return count, err
}

// CountArtifacts returns the total count of artifacts
func (s *Service) CountArtifacts(ctx context.Context) (int64, error) {
	var count int64
	err := s.DB.WithContext(ctx).Model(&model.Artifact{}).Count(&count).Error
	return count, err
}
