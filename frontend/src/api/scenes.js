import apiClient from './client'

export const fetchScenes = (homeId) =>
  apiClient.get('/scenes', { params: { home_id: homeId } })

export const createScene = (sceneData) =>
  apiClient.post('/scenes', sceneData)

export const activateScene = (sceneId) =>
  apiClient.patch(`/scenes/${sceneId}/activate`)

export const pauseScene = (sceneId) =>
  apiClient.patch(`/scenes/${sceneId}/pause`)

export const deleteScene = (sceneId) =>
  apiClient.delete(`/scenes/${sceneId}`)