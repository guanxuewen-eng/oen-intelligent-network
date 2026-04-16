import api from './index'
import type { ApiResponse, PaginatedResult } from '@/types'

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

export interface RecommendationListParams {
  agent_id?: number
  state?: string
  page?: number
  page_size?: number
}

export function getRecommendations(params?: RecommendationListParams) {
  return api.get<ApiResponse<PaginatedResult<Recommendation>>>('/recommendations', { params })
}

export function getRecommendation(id: number) {
  return api.get<ApiResponse<Recommendation>>(`/recommendations/${id}`)
}

export interface DecisionPayload {
  decision: 'accept' | 'ignore' | 'later'
  decided_by?: string
  note?: string
}

export function makeDecision(id: number, data: DecisionPayload) {
  return api.post<ApiResponse<null>>(`/recommendations/${id}/decision`, data)
}

// --- Candidates ---

export interface Candidate {
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

export interface CandidateListParams {
  state?: string
  source_type?: string
  page?: number
  page_size?: number
}

export function getCandidates(params?: CandidateListParams) {
  return api.get<ApiResponse<PaginatedResult<Candidate>>>('/candidates', { params })
}

export interface ReviewCandidatePayload {
  state: string
  review_notes?: string
  reviewed_by?: string
}

export function reviewCandidate(id: number, data: ReviewCandidatePayload) {
  return api.post<ApiResponse<null>>(`/candidates/${id}/review`, data)
}

// --- Ops ---

export function getOpsStatus() {
  return api.get<ApiResponse<any>>('/ops/status')
}

export function getAuditLogs(params?: { event_type?: string; page?: number; page_size?: number }) {
  return api.get<ApiResponse<PaginatedResult<any>>>('/ops/audit', { params })
}
