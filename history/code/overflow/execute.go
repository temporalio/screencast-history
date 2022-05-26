package main

import (
	"context"
	"log"

	"github.com/temporalio/screencasts/history/zapadapter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap/zapcore"
)

var skipLogKeys = []string{
	"Namespace",
	"TaskQueue",
	"WorkerID",
}

func main() {
	logger, err := zapadapter.NewZapLogger(zapcore.InfoLevel, skipLogKeys)
	if err != nil {
		log.Fatalln("Unable to create logger", err)
	}

	c, err := client.NewClient(client.Options{Logger: logger})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, "default", worker.Options{})

	w.RegisterWorkflow(OverflowWorkflow)

	err = w.Start()
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}

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

	log.Println("Workflow result:", result)
}
