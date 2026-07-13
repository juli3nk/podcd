package handler

import (
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/juli3nk/podcd/internal/version"
)

func (h *Handler) version(req ipc.Request) ipc.Response {
	return ipc.Response{
		Success: true,
		Data:    version.New(),
	}
}
