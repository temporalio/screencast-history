package main

import (
	"context"
	"fmt"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/interceptor"
	"go.temporal.io/sdk/workflow"
)

type workerInterceptor struct {
	interceptor.WorkerInterceptorBase
	options InterceptorOptions
}

type InterceptorOptions struct{}

func NewWorkerInterceptor(options InterceptorOptions) interceptor.WorkerInterceptor {
	return &workerInterceptor{options: options}
}

func (w *workerInterceptor) InterceptActivity(
	ctx context.Context,
	next interceptor.ActivityInboundInterceptor,
) interceptor.ActivityInboundInterceptor {
	i := &activityInboundInterceptor{root: w}
	i.Next = next
	return i
}

type activityInboundInterceptor struct {
	interceptor.ActivityInboundInterceptorBase
	root *workerInterceptor
}

func (a *activityInboundInterceptor) ExecuteActivity(ctx context.Context, in *interceptor.ExecuteActivityInput) (interface{}, error) {
	info := activity.GetInfo(ctx)

	fmt.Printf("-- Activity executing: %s, attempt: %d\n", info.ActivityType.Name, info.Attempt)

	result, err := a.Next.ExecuteActivity(ctx, in)

	if err != nil {
		fmt.Printf("!- Activity failed: %v\n", err)
	} else {
		fmt.Printf("-- Activity completed\n")
	}

	return result, err
}

func (w *workerInterceptor) InterceptWorkflow(
	ctx workflow.Context,
	next interceptor.WorkflowInboundInterceptor,
) interceptor.WorkflowInboundInterceptor {
	i := &workflowInboundInterceptor{root: w}
	i.Next = next
	return i
}

type workflowInboundInterceptor struct {
	interceptor.WorkflowInboundInterceptorBase
	root *workerInterceptor
}

func (w *workflowInboundInterceptor) Init(outbound interceptor.WorkflowOutboundInterceptor) error {
	i := &workflowOutboundInterceptor{root: w.root}
	i.Next = outbound
	return w.Next.Init(i)
}

func (w *workflowInboundInterceptor) ExecuteWorkflow(ctx workflow.Context, in *interceptor.ExecuteWorkflowInput) (interface{}, error) {
	var prefix string
	var action string
	if workflow.IsReplaying(ctx) {
		prefix = ">>"
		action = "replaying"
	} else {
		prefix = " >"
		action = "executing"
	}

	fmt.Printf("%s Workflow %s\n", prefix, action)

	result, err := w.Next.ExecuteWorkflow(ctx, in)

	fmt.Printf("%s Workflow completed\n", prefix)

	return result, err
}

type workflowOutboundInterceptor struct {
	interceptor.WorkflowOutboundInterceptorBase
	root *workerInterceptor
}

func (w *workflowOutboundInterceptor) ExecuteActivity(ctx workflow.Context, activityType string, args ...interface{}) workflow.Future {
	var prefix string
	if workflow.IsReplaying(ctx) {
		prefix = ">>"
	} else {
		prefix = " >"
	}

	fmt.Printf("%s Workflow called SDK's ExecuteActivity function: %s\n", prefix, activityType)
	if workflow.IsReplaying(ctx) {
		fmt.Printf("-- Activity execution for %s skipped, result read from history\n", activityType)
	}

	return w.Next.ExecuteActivity(ctx, activityType, args...)
}
