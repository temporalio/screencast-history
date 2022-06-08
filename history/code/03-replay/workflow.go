package main

import (
	"context"
	"fmt"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func ReplayWorkflow(ctx workflow.Context, name string) (string, error) {
	ctx = workflow.WithActivityOptions(ctx, workflow.ActivityOptions{
		StartToCloseTimeout: 1 * time.Second,
	})
	ctx = workflow.WithRetryPolicy(ctx, temporal.RetryPolicy{
		MaximumInterval: 1 * time.Second,
	})

	var result string
	err := workflow.ExecuteActivity(ctx, "FirstActivity", name+" first activity").Get(ctx, &result)
	if err != nil {
		return "", err
	}

	err = workflow.ExecuteActivity(ctx, "SecondActivity", name+" second activity").Get(ctx, &result)
	if err != nil {
		return "", err
	}

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) FirstActivity(ctx context.Context, name string) (string, error) {
	if activity.GetInfo(ctx).Attempt == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("faking a worker crash")
	}

	return "Hello from " + name + "!", nil
}

func (a *Activities) SecondActivity(ctx context.Context, name string) (string, error) {
	return "Hello from " + name + "!", nil
}
