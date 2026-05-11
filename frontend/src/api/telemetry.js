import apiClient from './client'

export const sendTelemetry = (deviceId, type, value) =>
  apiClient.post(`/telemetry/${deviceId}`, {
    type,
    value,
    timestamp: new Date().toISOString(),
  })

export const getLatestTelemetry = (deviceId) =>
  apiClient.get(`/telemetry/${deviceId}/latest`)

export const getTelemetryHistory = (deviceId, from, to) =>
  apiClient.get(`/telemetry/${deviceId}`, { params: { from, to } })