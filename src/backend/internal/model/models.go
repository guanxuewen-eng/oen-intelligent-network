package model

import (
	"time"
)

// Agent represents a sub-agent
type Agent struct {
	ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AgentKey      string     `gorm:"uniqueIndex;size:128;not null" json:"agent_key"`
	AgentType     string     `gorm:"size:32;not null;default:generic" json:"agent_type"`
	Name          string     `gorm:"size:256;not null" json:"name"`
	Role          string     `gorm:"size:128" json:"role"`
	State         string     `gorm:"size:64;default:unauthorized" json:"state"`
	RouteMode     string     `gorm:"size:64" json:"route_mode"`
	LastHeartbeat *time.Time `json:"last_heartbeat"`
	Metadata      *string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

func (Agent) TableName() string {
	return "agent"
}

// AgentHeartbeat represents heartbeat records
type AgentHeartbeat struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AgentID      uint       `gorm:"index;not null" json:"agent_id"`
	HeartbeatAt  time.Time  `gorm:"not null" json:"heartbeat_at"`
	Status       string     `gorm:"size:64;not null" json:"status"`
	RouteMode    string     `gorm:"size:64" json:"route_mode"`
	CPUUsage     float64    `json:"cpu_usage"`
	MemoryUsage  float64    `json:"memory_usage"`
	ErrorMessage string     `gorm:"type:text" json:"error_message"`
}

func (AgentHeartbeat) TableName() string {
	return "agent_heartbeat"
}

// Artifact represents the main asset table
type Artifact struct {
	ID                 uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	ArtifactKey        string     `gorm:"uniqueIndex;size:128;not null" json:"artifact_key"`
	ArtifactType       string     `gorm:"size:128;not null" json:"artifact_type"`
	Title              string     `gorm:"size:512;not null" json:"title"`
	Description        string     `gorm:"type:text" json:"description"`
	TargetSystem       string     `gorm:"size:256" json:"target_system"`
	ApplicableVersion  string     `gorm:"size:128" json:"applicable_version"`
	RiskLevel          string     `gorm:"size:64" json:"risk_level"`
	VerificationStatus string     `gorm:"size:64" json:"verification_status"`
	CreatorAgentID     *uint      `json:"creator_agent_id"`
	SourceDraftID      *uint      `json:"source_draft_id"`
	CurrentVersionID   *uint      `json:"current_version_id"`
	Metadata           *string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (Artifact) TableName() string {
	return "artifact"
}

// ArtifactVersion represents artifact versions
type ArtifactVersion struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArtifactID    uint      `gorm:"index;not null" json:"artifact_id"`
	VersionNumber string    `gorm:"size:64;not null" json:"version_number"`
	ChangeSummary string    `gorm:"type:text" json:"change_summary"`
	ContentJSON   *string   `gorm:"type:jsonb" json:"content_json"`
	CreatedAt     time.Time `json:"created_at"`
	CreatedBy     string    `gorm:"size:128" json:"created_by"`
}

func (ArtifactVersion) TableName() string {
	return "artifact_version"
}

// ArtifactView represents three-view of artifacts
type ArtifactView struct {
	ID              uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	ArtifactVersionID uint    `gorm:"index;not null" json:"artifact_version_id"`
	ViewType        string    `gorm:"size:64;not null" json:"view_type"`
	Content         string    `gorm:"type:text;not null" json:"content"`
	CreatedAt       time.Time `json:"created_at"`
}

func (ArtifactView) TableName() string {
	return "artifact_view"
}

// CandidateResource represents candidate resources
type CandidateResource struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	CandidateKey string     `gorm:"uniqueIndex;size:128;not null" json:"candidate_key"`
	AgentID      *uint      `json:"agent_id"`
	SourceType   string     `gorm:"size:64" json:"source_type"`
	Title        string     `gorm:"size:512;not null" json:"title"`
	Summary      string     `gorm:"type:text" json:"summary"`
	RawContent   string     `gorm:"type:text" json:"raw_content"`
	State        string     `gorm:"size:64;default:pending" json:"state"`
	ArtifactID   *uint      `json:"artifact_id"`
	ReviewNotes  string     `gorm:"type:text" json:"review_notes"`
	CreatedAt    time.Time  `json:"created_at"`
	ReviewedAt   *time.Time `json:"reviewed_at"`
	ReviewedBy   string     `gorm:"size:128" json:"reviewed_by"`
}

func (CandidateResource) TableName() string {
	return "candidate_resource"
}

// Recommendation represents recommendations
type Recommendation struct {
	ID                uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	RecommendationKey string     `gorm:"uniqueIndex;size:128;not null" json:"recommendation_key"`
	AgentID           uint       `gorm:"index;not null" json:"agent_id"`
	ArtifactID        uint       `gorm:"index;not null" json:"artifact_id"`
	Title             string     `gorm:"size:512;not null" json:"title"`
	MatchReason       string     `gorm:"type:text" json:"match_reason"`
	EnvironmentHint   string     `gorm:"type:text" json:"environment_hint"`
	RiskSummary       string     `gorm:"type:text" json:"risk_summary"`
	ConfidenceScore   float64    `json:"confidence_score"`
	SuggestedAction   string     `gorm:"type:text" json:"suggested_action"`
	State             string     `gorm:"size:64;default:pending" json:"state"`
	CreatedAt         time.Time  `json:"created_at"`
	DecidedAt         *time.Time `json:"decided_at"`
}

func (Recommendation) TableName() string {
	return "recommendation"
}

// RecommendationDecision represents user decisions on recommendations
type RecommendationDecision struct {
	ID               uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RecommendationID uint      `gorm:"uniqueIndex;not null" json:"recommendation_id"`
	Decision         string    `gorm:"size:64;not null" json:"decision"`
	DecidedBy        string    `gorm:"size:128;not null" json:"decided_by"`
	DecidedAt        time.Time `json:"decided_at"`
	Note             string    `gorm:"type:text" json:"note"`
}

func (RecommendationDecision) TableName() string {
	return "recommendation_decision"
}

// ConsentRecord represents authorization records
type ConsentRecord struct {
	ID          uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	AgentID     uint       `gorm:"index;not null" json:"agent_id"`
	ConsentType string     `gorm:"size:128;not null" json:"consent_type"`
	Status      string     `gorm:"size:64;default:active" json:"status"`
	GrantedAt   time.Time  `json:"granted_at"`
	RevokedAt   *time.Time `json:"revoked_at"`
	GrantedBy   string     `gorm:"size:128" json:"granted_by"`
}

func (ConsentRecord) TableName() string {
	return "consent_record"
}

// AuditLog represents audit logs
type AuditLog struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	EventType    string     `gorm:"size:128;not null" json:"event_type"`
	EventSubtype string     `gorm:"size:128" json:"event_subtype"`
	ActorType    string     `gorm:"size:64;not null" json:"actor_type"`
	ActorID      string     `gorm:"size:128" json:"actor_id"`
	TargetType   string     `gorm:"size:64" json:"target_type"`
	TargetID     string     `gorm:"size:128" json:"target_id"`
	Description  string     `gorm:"type:text" json:"description"`
	Metadata     *string    `gorm:"type:jsonb" json:"metadata"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_log"
}
