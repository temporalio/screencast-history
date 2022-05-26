package main

import (
	"context"
	"fmt"
	"log"

	"github.com/temporalio/screencasts/history/zapadapter"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var skipLogKeys = []string{
	"Namespace",
	"TaskQueue",
	"WorkflowID",
	"RunID",
	"WorkflowType",
}

func main() {
	encoderConfig := zap.NewDevelopmentEncoderConfig()
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("15:04:05")
	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "console",
		EncoderConfig:     encoderConfig,
		OutputPaths:       []string{"stdout"},
		ErrorOutputPaths:  []string{"stderr"},
	}
	logger, err := config.Build()

	go func() {
		for i := 0; ; i += 1 {
			c, err := client.NewClient(
				client.Options{
					Identity: fmt.Sprintf("worker %d", i),
					Logger:   zapadapter.NewZapAdapter(logger, skipLogKeys),
				},
			)
			if err != nil {
				log.Fatalln("Unable to create client", err)
			}

			w := worker.New(c, "default", worker.Options{
				EnableLoggingInReplay: true,
			})

			w.RegisterWorkflow(ReplayWorkflow)
			a := Activities{Worker: w}
			w.RegisterActivity(&a)

			err = w.Run(nil)
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
		}
	}()

	c, err := client.NewClient(client.Options{
		Logger: zapadapter.NewZapAdapter(logger, skipLogKeys),
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
