package main

import (
	"context"
	"log"

	"github.com/temporalio/screencasts/history/zapadapter"
	"go.temporal.io/sdk/client"
	sdklog "go.temporal.io/sdk/log"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap/zapcore"
)

var skipLogKeys = []string{
	"Namespace",
	"TaskQueue",
	"WorkflowType",
	"RunID",
}

func runWorker(identity string, logger sdklog.Logger) {
	c, err := client.NewClient(
		client.Options{
			Logger:   logger,
			Identity: identity,
		},
	)
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "default", worker.Options{})

	w.RegisterWorkflow(OverflowWorkflow)

	w.Run(nil)
}

func main() {
	logger, err := zapadapter.NewZapLogger(zapcore.InfoLevel, skipLogKeys)
	if err != nil {
		log.Fatalln("Unable to create logger", err)
	}

	go func() {
		runWorker("worker", logger)
	}()

	c, err := client.NewClient(client.Options{Logger: logger})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	wf, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "default",
		},
		OverflowWorkflow,
		"overflow workflow",
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = wf.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}

	logger.Info("Workflow completed", "result", result)
}
