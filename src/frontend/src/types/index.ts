// Common API response types
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
}

export interface ApiError {
  code: number
  message: string
  details: string
}

// Agent types
export interface Agent {
  id: number
  agent_key: string
  agent_type: string
  name: string
  role: string
  state: string
  route_mode: string
  last_heartbeat: string | null
  metadata: string | null
  created_at: string
  updated_at: string
}

export interface AgentHeartbeat {
  id: number
  agent_id: number
  heartbeat_at: string
  status: string
  route_mode: string
  cpu_usage: number
  memory_usage: number
  error_message: string
}

// Artifact types
export interface Artifact {
  id: number
  artifact_key: string
  artifact_type: string
  title: string
  description: string
  target_system: string
  applicable_version: string
  risk_level: string
  verification_status: string
  creator_agent_id: number | null
  source_draft_id: number | null
  current_version_id: number | null
  metadata: string | null
  created_at: string
  updated_at: string
}

export interface ArtifactVersion {
  id: number
  artifact_id: number
  version_number: string
  change_summary: string
  content_json: string
  created_at: string
  created_by: string
}

export interface ArtifactView {
  id: number
  artifact_version_id: number
  view_type: string
  content: string
  created_at: string
}

// Candidate resource types
export interface CandidateResource {
  id: number
  candidate_key: string
  agent_id: number | null
  source_type: string
  title: string
  summary: string
  raw_content: string
  state: string
  artifact_id: number | null
  review_notes: string
  created_at: string
  reviewed_at: string | null
  reviewed_by: string
}

// Recommendation types
export interface Recommendation {
  id: number
  recommendation_key: string
  agent_id: number
  artifact_id: number
  title: string
  match_reason: string
  environment_hint: string
  risk_summary: string
  confidence_score: number
  suggested_action: string
  state: string
  created_at: string
  decided_at: string | null
}

export interface RecommendationDecision {
  id: number
  recommendation_id: number
  decision: string
  decided_by: string
  decided_at: string
  note: string
}

// Consent types
export interface ConsentRecord {
  id: number
  agent_id: number
  consent_type: string
  status: string
  granted_at: string
  revoked_at: string | null
  granted_by: string
}

// Audit log types
export interface AuditLog {
  id: number
  event_type: string
  event_subtype: string
  actor_type: string
  actor_id: string
  target_type: string
  target_id: string
  description: string
  metadata: string
  created_at: string
}

// Health check
export interface HealthStatus {
  status: string
  service: string
  version: string
}

// Pagination
export interface PaginationParams {
  page: number
  page_size: number
}

export interface PaginatedResult<T> {
  items: T[]
  total: number
  page: number
  page_size: number
}
