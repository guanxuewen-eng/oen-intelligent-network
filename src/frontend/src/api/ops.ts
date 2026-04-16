import api from './index'
import type { HealthStatus } from '@/types'

export function getHealthStatus() {
  return api.get<HealthStatus>('/ops/health')
}
