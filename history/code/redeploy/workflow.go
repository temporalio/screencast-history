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

func RedeployWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info("Workflow execution", "name", name)

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, "RedeployActivity", name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	logger.Info("Workflow completed", "result", result)

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) RedeployActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info("Activity execution", "name", name)

	if rand.Intn(2) == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("Simulating a deploy")
	}

	return "Hello from " + name + "!", nil
}