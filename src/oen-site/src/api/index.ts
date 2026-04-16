import axios from 'axios'
import type { AxiosInstance } from 'axios'

const instance: AxiosInstance = axios.create({
  baseURL: '/api/v1',
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json',
  },
})

instance.interceptors.request.use(
  (config) => config,
  (error) => Promise.reject(error),
)

instance.interceptors.response.use(
  (response) => response,
  (error) => {
    if (error.response) {
      const { status, data } = error.response
      if (status === 401) {
        console.error('Unauthorized')
      } else if (status === 403) {
        console.error('Forbidden')
      } else if (status >= 500) {
        console.error('Server error:', data?.message || 'Internal server error')
      }
    } else if (error.request) {
      console.error('Network error: no response received')
    } else {
      console.error('Request error:', error.message)
    }
    return Promise.reject(error)
  },
)

export default instance
