package main

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/workflow"
)

func BasicWorkflow(ctx workflow.Context, name string) (string, error) {
	logger := workflow.GetLogger(ctx)

	logger.Info(" * Workflow executing")

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var result string
	err := workflow.ExecuteActivity(ctx, BasicActivity, name).Get(ctx, &result)
	if err != nil {
		logger.Error("Activity failed.", "Error", err)
		return "", err
	}

	return result, nil
}

func BasicActivity(ctx context.Context, name string) (string, error) {
	logger := activity.GetLogger(ctx)

	logger.Info(" * Activity executing")

	return "Hello from " + name + "!", nil
}
