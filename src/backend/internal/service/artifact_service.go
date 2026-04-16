package service

import (
	"context"
	"fmt"

	"github.com/xinewang/oen/internal/model"
)

// CreateArtifactVersion creates a new version for an artifact
func (s *Service) CreateArtifactVersion(ctx context.Context, ver *model.ArtifactVersion) error {
	if err := s.Repo.CreateArtifactVersion(ctx, ver); err != nil {
		return err
	}

	// Update artifact's current_version_id
	s.Repo.DB.WithContext(ctx).
		Model(&model.Artifact{}).
		Where("id = ?", ver.ArtifactID).
		Update("current_version_id", ver.ID)

	s.LogEvent(ctx, "artifact_version", "created", "system", "oen", "artifact", fmt.Sprintf("%d", ver.ArtifactID),
		fmt.Sprintf("Version %s created for artifact %d", ver.VersionNumber, ver.ArtifactID))

	return nil
}

// GetArtifactVersions returns all versions of an artifact
func (s *Service) GetArtifactVersions(ctx context.Context, artifactID uint) ([]model.ArtifactVersion, error) {
	return s.Repo.GetArtifactVersions(ctx, artifactID)
}

// CreateArtifactView creates a view for an artifact version
func (s *Service) CreateArtifactView(ctx context.Context, view *model.ArtifactView) error {
	return s.Repo.CreateArtifactView(ctx, view)
}

// GetArtifactView returns a specific view of an artifact version
func (s *Service) GetArtifactView(ctx context.Context, artifactVersionID uint, viewType string) (*model.ArtifactView, error) {
	return s.Repo.GetArtifactView(ctx, artifactVersionID, viewType)
}

// GetArtifactWithVersionsAndViews returns an artifact with all its versions and views
func (s *Service) GetArtifactWithVersionsAndViews(ctx context.Context, artifactID uint) (map[string]interface{}, error) {
	artifact, err := s.Repo.GetArtifactByID(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	versions, views, err := s.Repo.GetArtifactVersionsWithViews(ctx, artifactID)
	if err != nil {
		return nil, err
	}

	// Build version -> views map
	viewMap := make(map[uint][]model.ArtifactView)
	for _, v := range views {
		viewMap[v.ArtifactVersionID] = append(viewMap[v.ArtifactVersionID], v)
	}

	result := map[string]interface{}{
		"artifact": artifact,
		"versions": versions,
	}

	// Attach views to each version
	type VersionWithViews struct {
		model.ArtifactVersion
		Views []model.ArtifactView `json:"views"`
	}
	versionsWithViews := make([]VersionWithViews, len(versions))
	for i, v := range versions {
		versionsWithViews[i] = VersionWithViews{
			ArtifactVersion: v,
			Views:           viewMap[v.ID],
		}
	}
	result["versions"] = versionsWithViews

	return result, nil
}
