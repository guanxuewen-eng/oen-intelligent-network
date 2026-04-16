package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xinewang/oen/internal/model"
	"gorm.io/gorm"
)

// Consent handles POST /api/v1/agents/:id/consent
func (h *Handler) Consent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		ConsentType string `json:"consent_type" binding:"required"`
		GrantedBy   string `json:"granted_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if req.GrantedBy == "" {
		req.GrantedBy = "user"
	}

	consent, err := h.svc.Consent(c.Request.Context(), uint(id), req.ConsentType, req.GrantedBy)
	if err != nil {
		InternalError(c, "Consent failed: "+err.Error())
		return
	}

	OK(c, consent)
}

// RevokeConsent handles DELETE /api/v1/agents/:id/consent
func (h *Handler) RevokeConsent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		RevokedBy string `json:"revoked_by"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.RevokedBy == "" {
		req.RevokedBy = "user"
	}

	if err := h.svc.RevokeConsent(c.Request.Context(), uint(id), req.RevokedBy); err != nil {
		InternalError(c, "Revoke consent failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Consent revoked", nil)
}

// RebuildAgent handles POST /api/v1/agents/:id/rebuild
func (h *Handler) RebuildAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		TriggeredBy string `json:"triggered_by"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.TriggeredBy == "" {
		req.TriggeredBy = "user"
	}

	if err := h.svc.RebuildAgent(c.Request.Context(), uint(id), req.TriggeredBy); err != nil {
		InternalError(c, "Rebuild failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Agent rebuilt", nil)
}

// PauseAgent handles POST /api/v1/agents/:id/pause
func (h *Handler) PauseAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		TriggeredBy string `json:"triggered_by"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.TriggeredBy == "" {
		req.TriggeredBy = "user"
	}

	if err := h.svc.PauseAgent(c.Request.Context(), uint(id), req.TriggeredBy); err != nil {
		InternalError(c, "Pause failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Agent paused", nil)
}

// ResumeAgent handles POST /api/v1/agents/:id/resume
func (h *Handler) ResumeAgent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		TriggeredBy string `json:"triggered_by"`
	}
	_ = c.ShouldBindJSON(&req)
	if req.TriggeredBy == "" {
		req.TriggeredBy = "user"
	}

	if err := h.svc.ResumeAgent(c.Request.Context(), uint(id), req.TriggeredBy); err != nil {
		InternalError(c, "Resume failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Agent resumed", nil)
}

// Heartbeat handles POST /api/v1/agents/:id/heartbeat
func (h *Handler) Heartbeat(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	var req struct {
		Status      string  `json:"status"`
		RouteMode   string  `json:"route_mode"`
		CPUUsage    float64 `json:"cpu_usage"`
		MemoryUsage float64 `json:"memory_usage"`
		ErrorMessage string `json:"error_message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	if err := h.svc.RecordHeartbeat(c.Request.Context(), uint(id), req.Status, req.RouteMode, req.CPUUsage, req.MemoryUsage, req.ErrorMessage); err != nil {
		InternalError(c, "Heartbeat failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Heartbeat recorded", nil)
}

// GetHeartbeats handles GET /api/v1/agents/:id/heartbeats
func (h *Handler) GetHeartbeats(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid agent ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	heartbeats, err := h.svc.GetRecentHeartbeats(c.Request.Context(), uint(id), limit)
	if err != nil {
		InternalError(c, "Failed to get heartbeats: "+err.Error())
		return
	}

	OK(c, heartbeats)
}

// --- Artifact Version & View Handlers ---

// GetArtifactVersions handles GET /api/v1/artifacts/:id/versions
func (h *Handler) GetArtifactVersions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	versions, err := h.svc.GetArtifactVersions(c.Request.Context(), uint(id))
	if err != nil {
		InternalError(c, "Failed to get versions: "+err.Error())
		return
	}

	OK(c, versions)
}

// GetArtifactView handles GET /api/v1/artifacts/:id/versions/:ver/view/:viewType
func (h *Handler) GetArtifactView(c *gin.Context) {
	artifactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	versionNumber := c.Param("ver")
	viewType := c.Param("viewType")

	// Find the version
	version, err := h.svc.Repo.GetArtifactVersionByNumber(c.Request.Context(), uint(artifactID), versionNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Version not found")
			return
		}
		InternalError(c, "Failed to get version: "+err.Error())
		return
	}

	view, err := h.svc.GetArtifactView(c.Request.Context(), version.ID, viewType)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "View not found")
			return
		}
		InternalError(c, "Failed to get view: "+err.Error())
		return
	}

	OK(c, view)
}

// GetArtifactDetail handles GET /api/v1/artifacts/:id with versions and views
func (h *Handler) GetArtifactDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	result, err := h.svc.GetArtifactWithVersionsAndViews(c.Request.Context(), uint(id))
	if err != nil {
		InternalError(c, "Failed to get artifact detail: "+err.Error())
		return
	}

	OK(c, result)
}

// CreateArtifactVersion handles POST /api/v1/artifacts/:id/versions
func (h *Handler) CreateArtifactVersion(c *gin.Context) {
	artifactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	var req struct {
		VersionNumber string `json:"version_number" binding:"required"`
		ChangeSummary string `json:"change_summary"`
		ContentJSON   string `json:"content_json"`
		CreatedBy     string `json:"created_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if req.CreatedBy == "" {
		req.CreatedBy = "system"
	}

	ver := &model.ArtifactVersion{
		ArtifactID:    uint(artifactID),
		VersionNumber: req.VersionNumber,
		ChangeSummary: req.ChangeSummary,
		CreatedBy:     req.CreatedBy,
	}
	if req.ContentJSON != "" {
		ver.ContentJSON = &req.ContentJSON
	}

	if err := h.svc.CreateArtifactVersion(c.Request.Context(), ver); err != nil {
		InternalError(c, "Failed to create version: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "Version created",
		Data:    ver,
	})
}

// CreateArtifactView handles POST /api/v1/artifacts/:id/versions/:ver/view/:viewType
func (h *Handler) CreateArtifactView(c *gin.Context) {
	artifactID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid artifact ID")
		return
	}

	versionNumber := c.Param("ver")
	viewType := c.Param("viewType")

	version, err := h.svc.Repo.GetArtifactVersionByNumber(c.Request.Context(), uint(artifactID), versionNumber)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			NotFound(c, "Version not found")
			return
		}
		InternalError(c, "Failed to get version: "+err.Error())
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	view := &model.ArtifactView{
		ArtifactVersionID: version.ID,
		ViewType:          viewType,
		Content:           req.Content,
	}
	if err := h.svc.CreateArtifactView(c.Request.Context(), view); err != nil {
		InternalError(c, "Failed to create view: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "View created",
		Data:    view,
	})
}

// --- Recommendation Handlers ---

// CreateRecommendation handles POST /api/v1/recommendations
func (h *Handler) CreateRecommendation(c *gin.Context) {
	var req struct {
		RecommendationKey string  `json:"recommendation_key" binding:"required"`
		AgentID           uint    `json:"agent_id" binding:"required"`
		ArtifactID        uint    `json:"artifact_id" binding:"required"`
		Title             string  `json:"title" binding:"required"`
		MatchReason       string  `json:"match_reason"`
		EnvironmentHint   string  `json:"environment_hint"`
		RiskSummary       string  `json:"risk_summary"`
		ConfidenceScore   float64 `json:"confidence_score"`
		SuggestedAction   string  `json:"suggested_action"`
		State             string  `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	rec := &model.Recommendation{
		RecommendationKey: req.RecommendationKey,
		AgentID:           req.AgentID,
		ArtifactID:        req.ArtifactID,
		Title:             req.Title,
		MatchReason:       req.MatchReason,
		EnvironmentHint:   req.EnvironmentHint,
		RiskSummary:       req.RiskSummary,
		ConfidenceScore:   req.ConfidenceScore,
		SuggestedAction:   req.SuggestedAction,
		State:             req.State,
	}
	if rec.State == "" {
		rec.State = "pending"
	}

	if err := h.svc.CreateRecommendation(c.Request.Context(), rec); err != nil {
		InternalError(c, "Failed to create recommendation: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "Recommendation created",
		Data:    rec,
	})
}

// ListRecommendations handles GET /api/v1/recommendations
func (h *Handler) ListRecommendations(c *gin.Context) {
	agentID, _ := strconv.ParseUint(c.DefaultQuery("agent_id", "0"), 10, 64)
	state := c.Query("state")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	recs, total, err := h.svc.ListRecommendations(c.Request.Context(), uint(agentID), state, page, pageSize)
	if err != nil {
		InternalError(c, "Failed to list recommendations: "+err.Error())
		return
	}

	OK(c, gin.H{
		"items":     recs,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// GetRecommendation handles GET /api/v1/recommendations/:id
func (h *Handler) GetRecommendation(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid recommendation ID")
		return
	}

	rec, err := h.svc.GetRecommendation(c.Request.Context(), uint(id))
	if err != nil {
		InternalError(c, "Failed to get recommendation: "+err.Error())
		return
	}

	OK(c, rec)
}

// MakeDecision handles POST /api/v1/recommendations/:id/decision
func (h *Handler) MakeDecision(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid recommendation ID")
		return
	}

	var req struct {
		Decision  string `json:"decision" binding:"required"` // accept/ignore/later
		DecidedBy string `json:"decided_by"`
		Note      string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if req.DecidedBy == "" {
		req.DecidedBy = "user"
	}

	if err := h.svc.MakeDecision(c.Request.Context(), uint(id), req.Decision, req.DecidedBy, req.Note); err != nil {
		InternalError(c, "Decision failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Decision recorded", nil)
}

// --- Candidate Handlers ---

// CreateCandidate handles POST /api/v1/candidates
func (h *Handler) CreateCandidate(c *gin.Context) {
	var req struct {
		CandidateKey string `json:"candidate_key" binding:"required"`
		AgentID      *uint  `json:"agent_id"`
		SourceType   string `json:"source_type"`
		Title        string `json:"title" binding:"required"`
		Summary      string `json:"summary"`
		RawContent   string `json:"raw_content"`
		State        string `json:"state"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}

	candidate := &model.CandidateResource{
		CandidateKey: req.CandidateKey,
		AgentID:      req.AgentID,
		SourceType:   req.SourceType,
		Title:        req.Title,
		Summary:      req.Summary,
		RawContent:   req.RawContent,
		State:        req.State,
	}
	if candidate.State == "" {
		candidate.State = "pending"
	}

	if err := h.svc.CreateCandidate(c.Request.Context(), candidate); err != nil {
		InternalError(c, "Failed to create candidate: "+err.Error())
		return
	}

	c.JSON(http.StatusCreated, SuccessResponse{
		Code:    0,
		Message: "Candidate created",
		Data:    candidate,
	})
}

// ListCandidates handles GET /api/v1/candidates
func (h *Handler) ListCandidates(c *gin.Context) {
	state := c.Query("state")
	sourceType := c.Query("source_type")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	candidates, total, err := h.svc.ListCandidates(c.Request.Context(), state, sourceType, page, pageSize)
	if err != nil {
		InternalError(c, "Failed to list candidates: "+err.Error())
		return
	}

	OK(c, gin.H{
		"items":     candidates,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// ReviewCandidate handles POST /api/v1/candidates/:id/review
func (h *Handler) ReviewCandidate(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		BadRequest(c, "Invalid candidate ID")
		return
	}

	var req struct {
		State       string `json:"state" binding:"required"`
		ReviewNotes string `json:"review_notes"`
		ReviewedBy  string `json:"reviewed_by"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		BadRequest(c, "Invalid request: "+err.Error())
		return
	}
	if req.ReviewedBy == "" {
		req.ReviewedBy = "user"
	}

	if err := h.svc.ReviewCandidate(c.Request.Context(), uint(id), req.State, req.ReviewNotes, req.ReviewedBy, nil); err != nil {
		InternalError(c, "Review failed: "+err.Error())
		return
	}

	OKWithMessage(c, "Candidate reviewed", nil)
}
