package handler

// DeviceDTO – объект для передачи данных об устройстве клиенту
type DeviceDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Type   string `json:"type"`
	Status string `json:"status"`
}

type CommandRequest struct {
    Command string                 `json:"command"`
    Params  map[string]interface{} `json:"params"`
}