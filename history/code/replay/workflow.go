package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

func ReplayWorkflow(ctx workflow.Context, name string) (string, error) {
	if workflow.IsReplaying(ctx) {
		log.Println(" * REPLAY: Workflow code executing")
	} else {
		log.Println(" * Workflow code executing")
	}

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	if workflow.IsReplaying(ctx) {
		log.Println(" * REPLAY: Execute first activity")
	} else {
		log.Println(" * Execute first activity")
	}

	var result string
	err := workflow.ExecuteActivity(ctx, "ReplayActivity", name+" first activity").Get(ctx, &result)
	if err != nil {
		log.Println("Activity failed.", "Error", err)
		return "", err
	}

	if workflow.IsReplaying(ctx) {
		log.Println(" * REPLAY: Execute second activity")
	} else {
		log.Println(" * Execute second activity")
	}

	err = workflow.ExecuteActivity(ctx, "ReplayActivity", name+" second activity").Get(ctx, &result)
	if err != nil {
		log.Println("Activity failed.", "Error", err)
		return "", err
	}

	log.Println(" * Workflow completed", "result", result)

	return result, nil
}

type Activities struct {
	Worker worker.Worker
}

func (a *Activities) ReplayActivity(ctx context.Context, name string) (string, error) {
	if rand.Intn(2) == 1 {
		a.Worker.Stop()
		return "", fmt.Errorf("Simulating a deploy")
	}

	return "Hello from " + name + "!", nil
}
