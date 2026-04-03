import axios from 'axios'

const instance = axios.create({
  baseURL: '/api',
  timeout: 30000,
})

instance.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }
  return config
})

instance.interceptors.response.use(
  res => {
    if (res.data && res.data.code === -1) {
      return Promise.reject(new Error(res.data.error || 'Request failed'))
    }
    return res.data
  },
  err => {
    if (err.response?.status === 401) {
      localStorage.removeItem('token')
      window.location.href = '/admin/login'
    }
    return Promise.reject(err)
  }
)

const api = {
  // Auth
  login: (password) => instance.post('/auth/login', { password }),
  changePassword: (data) => instance.post('/v1/auth/change-password', data),

  // Stats
  overview: () => instance.get('/v1/stats/overview'),
  getClicks: (params) => instance.get('/v1/stats/clicks', { params }),

  // Templates
  listTemplates: () => instance.get('/v1/templates'),
  createTemplate: (data) => instance.post('/v1/templates', data),
  updateTemplate: (id, data) => instance.put(`/v1/templates/${id}`, data),
  deleteTemplate: (id) => instance.delete(`/v1/templates/${id}`),

  // Domains
  listDomains: (params) => instance.get('/v1/domains', { params }),
  createDomain: (data) => instance.post('/v1/domains', data),
  getDomain: (id) => instance.get(`/v1/domains/${id}`),
  updateDomain: (id, data) => instance.put(`/v1/domains/${id}`, data),
  deleteDomain: (id) => instance.delete(`/v1/domains/${id}`),
  bindTemplate: (id, templateId) => instance.post(`/v1/domains/${id}/bind-template`, { template_id: templateId }),
  batchDomainOp: (data) => instance.post('/v1/domains/batch', data),

  // Content
  getContent: (domainId) => instance.get(`/v1/content/${domainId}`),
  saveContent: (domainId, data) => instance.put(`/v1/content/${domainId}`, data),

  // Keywords
  listKeywords: (params) => instance.get('/v1/keywords', { params }),
  keywordCategories: () => instance.get('/v1/keywords/categories'),
  importKeywords: (formData) => instance.post('/v1/keywords/import', formData, { headers: { 'Content-Type': 'multipart/form-data' } }),
  assignKeyword: (id, domainId) => instance.post(`/v1/keywords/${id}/assign`, { domain_id: domainId }),
  deleteKeyword: (id) => instance.delete(`/v1/keywords/${id}`),

  // Build / Deploy
  buildSite: (id) => instance.post(`/v1/build/${id}`),
  deploySite: (id) => instance.post(`/v1/build/${id}/deploy`),
  buildFull: (id) => instance.post(`/v1/build/${id}/full`),
  getBuildStatus: (id) => instance.get(`/v1/build/${id}/status`),
  getBuildLog: (id) => instance.get(`/v1/build/${id}/log`),
  batchBuild: (data) => instance.post('/v1/build/batch', data),

  // Settings
  getSettings: () => instance.get('/v1/settings'),
  saveSettings: (data) => instance.put('/v1/settings', data),

  // V4: IndexNow
  submitIndexNow: (id, urls) => instance.post(`/v1/indexnow/${id}/submit`, { urls }),
  getIndexNowRecords: (id) => instance.get(`/v1/indexnow/${id}/records`),

  // V4: City Matrix
  getCityMatrix: (id) => instance.get(`/v1/city-matrix/${id}`),
  saveCityMatrix: (id, data) => instance.put(`/v1/city-matrix/${id}`, data),

  // V4: Title Pool
  listTitlePool: (params) => instance.get('/v1/title-pool', { params }),
  createTitleVariant: (data) => instance.post('/v1/title-pool', data),
  deleteTitleVariant: (id) => instance.delete(`/v1/title-pool/${id}`),

  // V4: Clusters
  listClusters: () => instance.get('/v1/clusters'),
  createCluster: (data) => instance.post('/v1/clusters', data),
  deleteCluster: (id) => instance.delete(`/v1/clusters/${id}`),
  addClusterMember: (id, data) => instance.post(`/v1/clusters/${id}/members`, data),
  removeClusterMember: (id, domainId) => instance.delete(`/v1/clusters/${id}/members`, { data: { domain_id: domainId } }),

  // V4: Content Refresh
  listRefreshSchedule: () => instance.get('/v1/refresh-schedule'),
  saveRefreshSchedule: (data) => instance.post('/v1/refresh-schedule', data),

  // V4: Health
  checkSiteHealth: (id) => instance.get(`/v1/health-check/${id}`),
  getHealthAlerts: () => instance.get('/v1/health-check/alerts'),

  // V4: IP/UA Rules
  listIPRules: () => instance.get('/v1/ip-rules'),
  addIPRule: (data) => instance.post('/v1/ip-rules', data),
  deleteIPRule: (id) => instance.delete(`/v1/ip-rules/${id}`),
  listUARules: () => instance.get('/v1/ua-rules'),
  addUARule: (data) => instance.post('/v1/ua-rules', data),
  deleteUARule: (id) => instance.delete(`/v1/ua-rules/${id}`),

  // V4: Export
  exportDomain: (id) => {
    const token = localStorage.getItem('token')
    window.open(`/api/v1/export/${id}?token=${token}`)
  },

  // generic
  get: (url, params) => instance.get('/v1' + url, { params }),
  post: (url, data) => instance.post('/v1' + url, data),
}

export default api
