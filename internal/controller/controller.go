package controller

import (
	"github.com/juli3nk/podcd/internal/reconcile"
)

type Controller struct {
	reconciler *reconcile.Reconciler
}

func New(reconciler *reconcile.Reconciler) *Controller {
	return &Controller{
		reconciler: reconciler,
	}
}
