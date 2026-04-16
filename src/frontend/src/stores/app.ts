import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getHealthStatus } from '@/api/ops'
import type { HealthStatus } from '@/types'

export const useAppStore = defineStore('app', () => {
  const healthStatus = ref<HealthStatus | null>(null)
  const loading = ref(false)
  const sidebarCollapsed = ref(false)

  async function fetchHealth() {
    loading.value = true
    try {
      const { data } = await getHealthStatus()
      healthStatus.value = data
    } catch (error) {
      console.error('Failed to fetch health status:', error)
      healthStatus.value = null
    } finally {
      loading.value = false
    }
  }

  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    healthStatus,
    loading,
    sidebarCollapsed,
    fetchHealth,
    toggleSidebar,
  }
})
