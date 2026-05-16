import axios from 'axios'

const API_BASE = '/api'

const api = axios.create({
  baseURL: API_BASE,
  timeout: 10000,
})

// Add auth token to requests
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('authToken')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

export const apiClient = {
  // System Info
  getSystemInfo: () => api.get('/sysinfo'),
  getSystemIdentity: () => api.get('/sysinfo/identity'),
  getCPUInfo: () => api.get('/sysinfo/cpu'),
  getMemoryInfo: () => api.get('/sysinfo/memory'),
  getDiskInfo: () => api.get('/sysinfo/disk'),
  getNetworkInfo: () => api.get('/sysinfo/network'),
  getSecurityInfo: () => api.get('/sysinfo/security'),
  getProcesses: () => api.get('/sysinfo/processes'),
  
  // Rules & Audit
  getAuditResults: () => api.get('/audit'),
  getAuditByCategory: (category) => api.get(`/audit?category=${category}`),
  
  // Blog/Advisory Content
  getBlogs: () => api.get('/blogs'),
  getBlogById: (id) => api.get(`/blogs/${id}`),
  searchBlogs: (query) => api.get(`/blogs/search?q=${query}`),
  
  // Authentication
  login: (username, password) => api.post('/auth/login', { username, password }),
  logout: () => api.post('/auth/logout'),
  getCurrentUser: () => api.get('/auth/user'),
  
  // Status
  getStatus: () => api.get('/status'),
}

export default api
