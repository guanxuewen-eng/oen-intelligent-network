package repository

import (
	"context"

	"github.com/xinewang/oen/internal/model"
)

// --- Artifact CRUD ---

// CreateArtifact creates a new artifact record
func (r *Repository) CreateArtifact(ctx context.Context, artifact *model.Artifact) error {
	return r.DB.WithContext(ctx).Create(artifact).Error
}

// GetArtifactByID retrieves an artifact by ID
func (r *Repository) GetArtifactByID(ctx context.Context, id uint) (*model.Artifact, error) {
	var artifact model.Artifact
	err := r.DB.WithContext(ctx).First(&artifact, id).Error
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}

// GetArtifactByKey retrieves an artifact by artifact_key
func (r *Repository) GetArtifactByKey(ctx context.Context, key string) (*model.Artifact, error) {
	var artifact model.Artifact
	err := r.DB.WithContext(ctx).Where("artifact_key = ?", key).First(&artifact).Error
	if err != nil {
		return nil, err
	}
	return &artifact, nil
}

// ListArtifacts returns paginated artifacts with optional filters
func (r *Repository) ListArtifacts(ctx context.Context, artifactType, riskLevel, verificationStatus string, offset, limit int) ([]model.Artifact, int64, error) {
	var artifacts []model.Artifact
	var total int64

	query := r.DB.WithContext(ctx).Model(&model.Artifact{})

	if artifactType != "" {
		query = query.Where("artifact_type = ?", artifactType)
	}
	if riskLevel != "" {
		query = query.Where("risk_level = ?", riskLevel)
	}
	if verificationStatus != "" {
		query = query.Where("verification_status = ?", verificationStatus)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&artifacts).Error; err != nil {
		return nil, 0, err
	}

	return artifacts, total, nil
}

// UpdateArtifact updates modifiable fields of an artifact
func (r *Repository) UpdateArtifact(ctx context.Context, artifact *model.Artifact) error {
	return r.DB.WithContext(ctx).Save(artifact).Error
}

// DeleteArtifact removes an artifact record
func (r *Repository) DeleteArtifact(ctx context.Context, id uint) error {
	return r.DB.WithContext(ctx).Delete(&model.Artifact{}, id).Error
}
