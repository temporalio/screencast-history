package main

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

// doSomeWork simulates some useful work
func doSomeWork(ctx workflow.Context) {
	// Normally we'd run activities here to do something useful.
	// In order for the demo to run quickly we just create lots of timers to inflate the history size.
	for i := 0; i < 2000; i++ {
		workflow.NewTimer(ctx, time.Hour)
	}
}

func OverflowWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("* Workflow executing")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for {
		doSomeWork(ctx)

		// Normally we'd sleep for whatever period of time our business logic dictates.
		// In order for the demo to run quickly we'll just have a quick nap and then carry on again.
		workflow.Sleep(ctx, 1)
	}
}
