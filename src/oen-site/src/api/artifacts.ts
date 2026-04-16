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
