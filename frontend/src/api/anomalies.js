import apiClient from './client'

export const fetchAnomalies = (params = {}) =>
  apiClient.get('/anomalies', { params })

export const updateAnomalyStatus = (id, status) =>
  apiClient.patch(`/anomalies/${id}`, { status })