package main

import (
	"context"
	"fmt"
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
	"WorkflowID",
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

	w := worker.New(c, "default", worker.Options{
		EnableLoggingInReplay: true,
	})

	w.RegisterWorkflow(ReplayWorkflow)
	a := Activities{Worker: w}
	w.RegisterActivity(&a)

	w.Run(nil)
}

func main() {
	logger, err := zapadapter.NewZapLogger(zapcore.DebugLevel, skipLogKeys)
	if err != nil {
		log.Fatalln("Unable to create logger", err)
	}

	go func() {
		for i := 0; ; i += 1 {
			runWorker(fmt.Sprintf("worker %d", i), logger)
		}
	}()

	c, err := client.NewClient(client.Options{
		Logger: logger,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	wf, err := c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			TaskQueue: "default",
		},
		ReplayWorkflow,
		"replay workflow",
	)
	if err != nil {
		log.Fatalln("Unable to execute workflow", err)
	}

	var result string
	err = wf.Get(context.Background(), &result)
	if err != nil {
		log.Fatalln("Workflow failed", err)
	}
}
