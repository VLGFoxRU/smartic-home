import apiClient from './client'

export const fetchMe = () => apiClient.get('/auth/me')