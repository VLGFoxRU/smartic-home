import apiClient from './client'

export const getCurrentUser = () => apiClient.get('/auth/me')
export const changePassword = (oldPassword, newPassword) =>
  apiClient.post('/auth/change-password', { old_password: oldPassword, new_password: newPassword })