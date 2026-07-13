package handler

import "github.com/juli3nk/podcd/internal/ipc"

func (h *Handler) pubkeySSH(req ipc.Request) ipc.Response {
	key, err := h.identity.SSHPublicKey()
	if err != nil {
		return ipc.Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ipc.Response{
		Success: true,
		Data:    key,
	}
}

func (h *Handler) pubkeyAge(req ipc.Request) ipc.Response {
	key, err := h.identity.AgePublicKey()
	if err != nil {
		return ipc.Response{
			Success: false,
			Error:   err.Error(),
		}
	}

	return ipc.Response{
		Success: true,
		Data:    key,
	}
}
