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

export function getOpsStatus() {
  return api.get<ApiResponse<any>>('/ops/status')
}
