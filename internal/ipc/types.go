package ipc

type Request struct {
	Command string         `json:"command"`
	Params  map[string]any `json:"params,omitempty"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}
