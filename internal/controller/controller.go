package controller

import (
	"github.com/juli3nk/podcd/internal/reconcile"
)

type Controller struct {
	reconciler *reconcile.Reconciler
	interval   string
}

func New(reconciler *reconcile.Reconciler, interval string) *Controller {
	return &Controller{
		reconciler: reconciler,
		interval:   interval,
	}
}
