package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func ReplayWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("* Workflow executing", "Replay", workflow.IsReplaying(ctx))

	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Second,
	})
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{
		MaximumInterval: 1 * time.Second,
	})

	logger.Info("* Requesting ExecuteActivity for the first activity", "Replay", workflow.IsReplaying(ctx))

	var result string
	err := workflow.ExecuteActivity(ctx, "ReplayActivity", name+" first activity").Get(ctx, &result)
	if err != nil {
		return "", err
	}

	logger.Info("* Requesting ExecuteActivity for the second activity", "Replay", workflow.IsReplaying(ctx))

	err = workflow.ExecuteActivity(ctx, "ReplayActivity", name+" second activity").Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) ReplayActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("* Activity executing")

	if rand.Intn(2) == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("Simulating a deploy")
	}

	return "Hello from " + name + "!", nil
}
