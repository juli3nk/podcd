package handler

import (
	"github.com/juli3nk/podcd/internal/identity"
	"github.com/juli3nk/podcd/internal/ipc"
	"github.com/juli3nk/podcd/internal/reconcile"
	"github.com/juli3nk/podcd/internal/runtime"
)

type Handler struct {
	runtime    runtime.Backend
	reconciler *reconcile.Reconciler
	identity   identity.Manager

	handlers map[string]ipc.HandlerFunc
}

func New(runtimeBackend runtime.Backend, reconciler *reconcile.Reconciler, identity identity.Manager) *Handler {
	h := &Handler{
		runtime:    runtimeBackend,
		reconciler: reconciler,
		identity:   identity,
		handlers:   map[string]ipc.HandlerFunc{},
	}

	// h.handlers["context.list"] = h.contextList

	h.handlers["reqs"] = h.reqs

	h.handlers["pubkey.ssh"] = h.pubkeySSH
	h.handlers["pubkey.age"] = h.pubkeyAge

	// h.handlers["validate"] = h.validate

	// h.handlers["repo.list"] = h.repoList
	// h.handlers["repo.add"] = h.repoAdd
	// h.handlers["repo.delete"] = h.repoDelete

	// h.handlers["status"] = h.status
	// h.handlers["sync"] = h.sync

	h.handlers["version"] = h.version

	return h
}

func (h *Handler) Handle(req ipc.Request) ipc.Response {
	handler, ok := h.handlers[req.Command]
	if !ok {
		return ipc.Response{
			Success: false,
			Error:   "unknown command",
		}
	}

	return handler(req)
}
