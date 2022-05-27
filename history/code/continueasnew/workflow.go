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

func ContinueAsNewWorkflow(ctx workflow.Context, name string, cursor int) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("* Workflow executing", "cursor", cursor)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	for {
		doSomeWork(ctx)

		// Normally we'd sleep for whatever period of time our business logic dictates.
		// In order for the demo to run quickly we'll just have a quick nap and then carry on again.
		workflow.Sleep(ctx, 1)

		// For the purposes of the demo we'll end after 10 iterations
		if cursor >= 10 {
			return "finished", nil
		}

		// Re-run this workflow with a new cursor
		return "", workflow.NewContinueAsNewError(ctx, ContinueAsNewWorkflow, name, cursor+1)
	}
}
