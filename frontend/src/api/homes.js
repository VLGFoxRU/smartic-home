import apiClient from './client'

export const fetchHomes = () => apiClient.get('/homes')
export const createHome = (name) => apiClient.post('/homes', { name })
export const fetchHome = (id) => apiClient.get(`/homes/${id}`)