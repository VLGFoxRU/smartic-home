package handler

import "encoding/json"

type DeviceDTO struct {
    ID     string          `json:"id"`
    Name   string          `json:"name"`
    Type   string          `json:"type"`
    Status string          `json:"status"`
    RoomID string          `json:"room_id,omitempty"`
    State  json.RawMessage `json:"state,omitempty"`
}

type CommandRequest struct {
    Command string                 `json:"command"`
    Params  map[string]interface{} `json:"params"`
}