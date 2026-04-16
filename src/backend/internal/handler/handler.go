package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xinewang/oen/internal/model"
	"github.com/xinewang/oen/internal/service"
	"gorm.io/gorm"
)

type Handler struct {
	svc *service.Service
}

func New(svc *service.Service) *Handler {
	return &Handler{svc: svc}
}

// HealthCheck handles GET /api/v1/ops/health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "oen-intelligent-network",
		"version": "0.1.0",
	})
}

// SystemStatus handles GET /api/v1/ops/status
func (h *Handler) SystemStatus(c *gin.Context) {
	status, err := h.svc.GetStatusOverview(c.Request.Context())
	if err != nil {
		InternalError(c, "Failed to get status: "+err.Error())
		return
	}
	OK(c, status)
}

// ListAuditLogs handles GET /api/v1/ops/audit
func (h *Handler) ListAuditLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	eventType := c.Query("event_type")
	targetType := c.Query("target_type")
	targetID := c.Query("target_id")

	logs, total, err := h.svc.GetAuditLogs(c.Request.Context(), eventType, targetType, targetID, page, pageSize)
	if err != nil {
		InternalError(c, "Failed to list audit logs: "+err.Error())
		return
	}

	OK(c, gin.H{
		"items":     logs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// --- Agent Handlers ---

// ListAgents handles GET /api/v1/agents
func (h *Handler) ListAgents(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	agentType := c.Query("agent_type")
	state := c.Query("state")

	agents, total, err := h.svc.ListAgents(c.Request.Context(), agentType, state, page, pageSize)
	if err != nil {
		InternalError(c, "Failed to list agents: "+err.Error())
		return
	}

	OK(c, gin.H{
		"items":     agents,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetAgent handles GET /api/v1/agents/:id
func (h *Handler) GetAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	agent, err := h.svc.GetAgentByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Agent not found")
			return
		}
		InternalError(c, "Failed to get agent: "+err.Error())
		return
	}

	OK(c, agent)
}

// CreateAgent handles POST /api/v1/agents
func (h *Handler) CreateAgent(c *gin.Context) {
	var req struct {
		AgentKey  string `json:"agent_key" binding:"required"`
		AgentType string `json:"agent_type"`
		Name      string `json:"name" binding:"required"`
		Role      string `json:"role"`
		State     string `json:"state"`
		RouteMode string `json:"route_mode"`
		Metadata  string `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	agent := model.Agent{
		AgentKey:  req.AgentKey,
		AgentType: req.AgentType,
		Name:      req.Name,
		Role:      req.Role,
		State:     req.State,
		RouteMode: req.RouteMode,
	}

	if req.Metadata != "" {
		agent.Metadata = &req.Metadata
	}
	if agent.AgentType == "" {
		agent.AgentType = "generic"
	}
	if agent.State == "" {
		agent.State = "unauthorized"
	}

	if err := h.svc.CreateAgent(c.Request.Context(), &agent); err != nil {
		InternalError(c, "Failed to create agent: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "Agent created",
		Data:    agent,
	})
}

// UpdateAgent handles PUT /api/v1/agents/:id
func (h *Handler) UpdateAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	existing, err := h.svc.GetAgentByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Agent not found")
			return
		}
		InternalError(c, "Failed to get agent: "+err.Error())
		return
	}

	var req struct {
		Name      string `json:"name"`
		Role      string `json:"role"`
		State     string `json:"state"`
		RouteMode string `json:"route_mode"`
		Metadata  string `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.Name != "" {
		existing.Name = req.Name
	}
	if req.Role != "" {
		existing.Role = req.Role
	}
	if req.State != "" {
		existing.State = req.State
	}
	if req.RouteMode != "" {
		existing.RouteMode = req.RouteMode
	}
	if req.Metadata != "" {
		existing.Metadata = &req.Metadata
	}

	if err := h.svc.UpdateAgent(c.Request.Context(), existing); err != nil {
		InternalError(c, "Failed to update agent: "+err.Error())
		return
	}

	OK(c, existing)
}

// DeleteAgent handles DELETE /api/v1/agents/:id
func (h *Handler) DeleteAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	if err := h.svc.DeleteAgent(c.Request.Context(), uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Agent not found")
			return
		}
		InternalError(c, "Failed to delete agent: "+err.Error())
		return
	}

	OKWithMessage(c, "Agent deleted", nil)
}

// --- Artifact Handlers ---

// ListArtifacts handles GET /api/v1/artifacts
func (h *Handler) ListArtifacts(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	artifactType := c.Query("artifact_type")
	riskLevel := c.Query("risk_level")
	verificationStatus := c.Query("verification_status")

	artifacts, total, err := h.svc.ListArtifacts(c.Request.Context(), artifactType, riskLevel, verificationStatus, page, pageSize)
	if err != nil {
		InternalError(c, "Failed to list artifacts: "+err.Error())
		return
	}

	OK(c, gin.H{
		"items":     artifacts,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetArtifact handles GET /api/v1/artifacts/:id
func (h *Handler) GetArtifact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	artifact, err := h.svc.GetArtifactByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Artifact not found")
			return
		}
		InternalError(c, "Failed to get artifact: "+err.Error())
		return
	}

	OK(c, artifact)
}

// CreateArtifact handles POST /api/v1/artifacts
func (h *Handler) CreateArtifact(c *gin.Context) {
	var req struct {
		ArtifactKey       string `json:"artifact_key" binding:"required"`
		ArtifactType      string `json:"artifact_type" binding:"required"`
		Title             string `json:"title" binding:"required"`
		Description       string `json:"description"`
		TargetSystem      string `json:"target_system"`
		ApplicableVersion string `json:"applicable_version"`
		RiskLevel         string `json:"risk_level"`
		VerificationStatus string `json:"verification_status"`
		CreatorAgentID    *uint  `json:"creator_agent_id"`
		SourceDraftID     *uint  `json:"source_draft_id"`
		CurrentVersionID  *uint  `json:"current_version_id"`
		Metadata          string `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	artifact := model.Artifact{
		ArtifactKey:       req.ArtifactKey,
		ArtifactType:      req.ArtifactType,
		Title:             req.Title,
		Description:       req.Description,
		TargetSystem:      req.TargetSystem,
		ApplicableVersion: req.ApplicableVersion,
		RiskLevel:         req.RiskLevel,
		VerificationStatus: req.VerificationStatus,
		CreatorAgentID:    req.CreatorAgentID,
		SourceDraftID:     req.SourceDraftID,
		CurrentVersionID:  req.CurrentVersionID,
	}

	if req.Metadata != "" {
		artifact.Metadata = &req.Metadata
	}

	if err := h.svc.CreateArtifact(c.Request.Context(), &artifact); err != nil {
		InternalError(c, "Failed to create artifact: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "Artifact created",
		Data:    artifact,
	})
}

// UpdateArtifact handles PUT /api/v1/artifacts/:id
func (h *Handler) UpdateArtifact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	existing, err := h.svc.GetArtifactByID(c.Request.Context(), uint(id))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Artifact not found")
			return
		}
		InternalError(c, "Failed to get artifact: "+err.Error())
		return
	}

	var req struct {
		Title             string `json:"title"`
		Description       string `json:"description"`
		TargetSystem      string `json:"target_system"`
		ApplicableVersion string `json:"applicable_version"`
		RiskLevel         string `json:"risk_level"`
		VerificationStatus string `json:"verification_status"`
		Metadata          string `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if req.Title != "" {
		existing.Title = req.Title
	}
	if req.Description != "" {
		existing.Description = req.Description
	}
	if req.TargetSystem != "" {
		existing.TargetSystem = req.TargetSystem
	}
	if req.ApplicableVersion != "" {
		existing.ApplicableVersion = req.ApplicableVersion
	}
	if req.RiskLevel != "" {
		existing.RiskLevel = req.RiskLevel
	}
	if req.VerificationStatus != "" {
		existing.VerificationStatus = req.VerificationStatus
	}
	if req.Metadata != "" {
		existing.Metadata = &req.Metadata
	}

	if err := h.svc.UpdateArtifact(c.Request.Context(), existing); err != nil {
		InternalError(c, "Failed to update artifact: "+err.Error())
		return
	}

	OK(c, existing)
}

// DeleteArtifact handles DELETE /api/v1/artifacts/:id
func (h *Handler) DeleteArtifact(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	if err := h.svc.DeleteArtifact(c.Request.Context(), uint(id)); err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Artifact not found")
			return
		}
		InternalError(c, "Failed to delete artifact: "+err.Error())
		return
	}

	OKWithMessage(c, "Artifact deleted", nil)
}
