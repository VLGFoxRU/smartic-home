import apiClient from './client'

export const fetchRooms = (homeId) => apiClient.get(`/homes/${homeId}/rooms`)
export const createRoom = (homeId, name) => apiClient.post(`/homes/${homeId}/rooms`, { name })