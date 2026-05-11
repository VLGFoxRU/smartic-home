import apiClient from './client'

export const fetchDevices = () => 
  apiClient.get('/devices')

export const createDevice = (name, type, roomId) =>
  apiClient.post('/devices', { name, type, room_id: roomId })

export const getDeviceState = (deviceId) =>
  apiClient.get(`/devices/${deviceId}`)

export const sendCommand = (deviceId, command, params) =>
  apiClient.post(`/devices/${deviceId}/command`, { command, params })

export const deleteDevice = (deviceId) =>
  apiClient.delete(`/devices/${deviceId}`)