import api from './index'
import type { Artifact, PaginatedResult, ApiResponse } from '@/types'

export interface ArtifactListParams {
  page?: number
  page_size?: number
  artifact_type?: string
  risk_level?: string
}

export function getArtifacts(params?: ArtifactListParams) {
  return api.get<ApiResponse<PaginatedResult<Artifact>>>('/artifacts', { params })
}

export function getArtifact(id: number) {
  return api.get<ApiResponse<Artifact>>(`/artifacts/${id}`)
}

export function getArtifactDetail(id: number) {
  return api.get<ApiResponse<any>>(`/artifacts/${id}/detail`)
}

export interface CreateArtifactPayload {
  artifact_key: string
  artifact_type: string
  title: string
  description?: string
  target_system?: string
  applicable_version?: string
  risk_level?: string
  verification_status?: string
  creator_agent_id?: number | null
  source_draft_id?: number | null
  current_version_id?: number | null
  metadata?: string
}

export function createArtifact(data: CreateArtifactPayload) {
  return api.post<ApiResponse<Artifact>>('/artifacts', data)
}

export interface UpdateArtifactPayload {
  title?: string
  description?: string
  target_system?: string
  applicable_version?: string
  risk_level?: string
  verification_status?: string
  metadata?: string
}

export function updateArtifact(id: number, data: UpdateArtifactPayload) {
  return api.put<ApiResponse<Artifact>>(`/artifacts/${id}`, data)
}

export function deleteArtifact(id: number) {
  return api.delete<ApiResponse<null>>(`/artifacts/${id}`)
}
