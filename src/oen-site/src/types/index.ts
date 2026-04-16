// Common API response types
export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
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
