package main

import (
	"context"
	"log"

	"github.com/temporalio/screencasts/history/zapadapter"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap/zapcore"
)

var skipLogKeys = []string{
	"Namespace",
	"TaskQueue",
	"WorkflowType",
	"RunID",
}

func main() {
	logger, err := zapadapter.NewZapLogger(zapcore.DebugLevel, skipLogKeys)
	if err != nil {
		log.Fatalln("Unable to create logger", err)
	}

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
		"NonExistantWorkflow",
		"hello?",
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
