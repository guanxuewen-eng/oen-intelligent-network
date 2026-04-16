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

export function getAgent(id: number) {
  return api.get<ApiResponse<Agent>>(`/agents/${id}`)
}

export interface CreateAgentPayload {
  agent_key: string
  agent_type?: string
  name: string
  role?: string
  state?: string
  route_mode?: string
  metadata?: string
}

export function createAgent(data: CreateAgentPayload) {
  return api.post<ApiResponse<Agent>>('/agents', data)
}

export interface UpdateAgentPayload {
  name?: string
  role?: string
  state?: string
  route_mode?: string
  metadata?: string
}

export function updateAgent(id: number, data: UpdateAgentPayload) {
  return api.put<ApiResponse<Agent>>(`/agents/${id}`, data)
}

export function deleteAgent(id: number) {
  return api.delete<ApiResponse<null>>(`/agents/${id}`)
}

// --- State Machine & Lifecycle ---

export interface ConsentPayload {
  consent_type: string
  granted_by?: string
}

export function consentAgent(id: number, data: ConsentPayload) {
  return api.post<ApiResponse<any>>(`/agents/${id}/consent`, data)
}

export function revokeConsentAgent(id: number, revoked_by?: string) {
  return api.delete<ApiResponse<null>>(`/agents/${id}/consent`, { data: { revoked_by } })
}

export function rebuildAgent(id: number, triggered_by?: string) {
  return api.post<ApiResponse<null>>(`/agents/${id}/rebuild`, { triggered_by })
}

export function pauseAgent(id: number, triggered_by?: string) {
  return api.post<ApiResponse<null>>(`/agents/${id}/pause`, { triggered_by })
}

export function resumeAgent(id: number, triggered_by?: string) {
  return api.post<ApiResponse<null>>(`/agents/${id}/resume`, { triggered_by })
}

// --- Heartbeat ---

export interface HeartbeatPayload {
  status?: string
  route_mode?: string
  cpu_usage?: number
  memory_usage?: number
  error_message?: string
}

export function sendHeartbeat(id: number, data: HeartbeatPayload) {
  return api.post<ApiResponse<null>>(`/agents/${id}/heartbeat`, data)
}

export function getHeartbeats(id: number, limit?: number) {
  return api.get<ApiResponse<any[]>>(`/agents/${id}/heartbeats`, { params: { limit } })
}
