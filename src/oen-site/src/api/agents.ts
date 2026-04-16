import api from './index'
import type { Agent, PaginatedResult, ApiResponse } from '@/types'

export interface AgentListParams {
  page?: number
  page_size?: number
  agent_type?: string
  state?: string
}

export function getAgents(params?: AgentListParams) {
  return api.get<ApiResponse<PaginatedResult<Agent>>>('/agents', { params })
}
