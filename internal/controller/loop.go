package controller

func (c *Controller) Run() {
    for {
        changes := c.Source.Diff()

        plan := c.Reconciler.Plan(changes)

        c.Reconciler.Apply(plan)

        time.Sleep(30 * time.Second)
    }
}
