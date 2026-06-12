package controller

import (
	"context"
	"log"
	"time"
)

func (c *Controller) Run(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if err := c.reconciler.Reconcile(); err != nil {
				log.Printf("reconcile error: %v", err)
			}
		}
	}
}
