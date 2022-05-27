package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func WorkerCrashWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("* Workflow executing")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, "WorkerCrashActivity", name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) WorkerCrashActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("* Activity executing")

	if rand.Intn(2) == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("Simulating a worker crash")
	}

	return "Hello from " + name + "!", nil
}
