package controller

import (
	"context"
	"fmt"
	"time"
)

func (c *Controller) Run(ctx context.Context) error {
	interval, err := time.ParseDuration(c.interval)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
			if err := c.reconciler.Reconcile(); err != nil {
				fmt.Printf("reconciliation error: %+v\n", err)
				continue
			}
		}
	}
}
