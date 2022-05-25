package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func ReplayWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info(" * Workflow executing", "replay", workflow.IsReplaying(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	logger.Info(" * Execute first activity", "Replay", workflow.IsReplaying(ctx))

	var result string
	err := workflow.ExecuteActivity(ctx, "ReplayActivity", name+" first activity").Get(ctx, &result)
	if err != nil {
		log.Println("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info(" * Execute second activity", "Replay", workflow.IsReplaying(ctx))

	err = workflow.ExecuteActivity(ctx, "ReplayActivity", name+" second activity").Get(ctx, &result)
	if err != nil {
		log.Println("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info(" * Workflow completed", "Replay", workflow.IsReplaying(ctx))

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) ReplayActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info(" * Activity executing")

	if rand.Intn(2) == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("Simulating a deploy")
	}

	return "Hello from " + name + "!", nil
}
