package ecs

import (
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	aws_cloudwatchlogs "github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	aws_ecs "github.com/aws/aws-sdk-go/service/ecs"
	"github.com/whosonfirst/go-whosonfirst-aws/cloudwatch"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	_ "log"
	"strings"
)

type TaskResponse struct {
	Tasks      []string
	TaskOutput *aws_ecs.RunTaskOutput
}

type TaskOptions struct {
	DSN            string
	Task           string
	Container      string
	Cluster        string
	LaunchType     string
	PublicIP       string
	Subnets        []string
	SecurityGroups []string
}

type MonitorTaskResultSet map[string]*MonitorTaskResult

type MonitorTaskResult struct {
	ARN    string
	Errors []error
	Logs   []*aws_cloudwatchlogs.OutputLogEvent
}

type MonitorTaskOptions struct {
	DSN       string
	Container string
	Cluster   string
	WithLogs  bool
	LogsDSN   string
}

func LaunchTask(task_opts *TaskOptions, cmd ...string) (*TaskResponse, error) {

	ecs_sess, err := session.NewSessionWithDSN(task_opts.DSN)

	if err != nil {
		return nil, err
	}

	ecs_svc := aws_ecs.New(ecs_sess)

	cluster := aws.String(task_opts.Cluster)
	task := aws.String(task_opts.Task)

	launch_type := aws.String(task_opts.LaunchType)
	public_ip := aws.String(task_opts.PublicIP)

	subnets := make([]*string, len(task_opts.Subnets))
	security_groups := make([]*string, len(task_opts.SecurityGroups))

	for i, sn := range task_opts.Subnets {
		subnets[i] = aws.String(sn)
	}

	for i, sg := range task_opts.SecurityGroups {
		security_groups[i] = aws.String(sg)
	}

	aws_cmd := make([]*string, len(cmd))

	for i, str := range cmd {
		aws_cmd[i] = aws.String(str)
	}

	network := &aws_ecs.NetworkConfiguration{
		AwsvpcConfiguration: &aws_ecs.AwsVpcConfiguration{
			AssignPublicIp: public_ip,
			SecurityGroups: security_groups,
			Subnets:        subnets,
		},
	}

	process_override := &aws_ecs.ContainerOverride{
		Name:    aws.String(task_opts.Container),
		Command: aws_cmd,
	}

	overrides := &aws_ecs.TaskOverride{
		ContainerOverrides: []*aws_ecs.ContainerOverride{
			process_override,
		},
	}

	input := &aws_ecs.RunTaskInput{
		Cluster:              cluster,
		TaskDefinition:       task,
		LaunchType:           launch_type,
		NetworkConfiguration: network,
		Overrides:            overrides,
	}

	task_output, err := ecs_svc.RunTask(input)

	if err != nil {
		return nil, err
	}

	if len(task_output.Tasks) == 0 {
		return nil, errors.New("run task returned no errors... but no tasks")
	}

	task_arns := make([]string, len(task_output.Tasks))

	for i, t := range task_output.Tasks {
		task_arns[i] = *t.TaskArn
	}

	task_rsp := &TaskResponse{
		Tasks:      task_arns,
		TaskOutput: task_output,
	}

	return task_rsp, nil
}

func MonitorTasks(monitor_opts *MonitorTaskOptions, task_arns ...string) (MonitorTaskResultSet, error) {

	ecs_sess, err := session.NewSessionWithDSN(monitor_opts.DSN)

	if err != nil {
		return nil, err
	}

	ecs_svc := aws_ecs.New(ecs_sess)

	count_tasks := len(task_arns)
	remaining := count_tasks

	ecs_tasks := make([]*string, count_tasks)

	for i, t := range task_arns {
		ecs_tasks[i] = aws.String(t)
	}

	result_set := make(map[string]*MonitorTaskResult)

	for remaining > 0 {

		monitor_req := &aws_ecs.DescribeTasksInput{
			Cluster: aws.String(monitor_opts.Cluster),
			Tasks:   ecs_tasks,
		}

		monitor_rsp, err := ecs_svc.DescribeTasks(monitor_req)

		if err != nil {
			return nil, err
		}

		for _, t := range monitor_rsp.Tasks {

			for _, c := range t.Containers {

				if *c.Name != monitor_opts.Container {
					continue
				}

				if *c.LastStatus != "STOPPED" {
					continue
				}

				task_arn := *t.TaskArn
				task_errors := make([]error, 0)
				task_logs := make([]*aws_cloudwatchlogs.OutputLogEvent, 0)

				if monitor_opts.WithLogs {

					arn := strings.Split(*t.TaskArn, "/")

					cw_group := fmt.Sprintf("/ecs/%s", monitor_opts.Container)
					cw_stream := fmt.Sprintf("ecs/%s/%s", monitor_opts.Container, arn[1])

					events, err := cloudwatch.GetLogEvents(monitor_opts.LogsDSN, cw_group, cw_stream)

					if err != nil {
						task_errors = append(task_errors, err)
					} else {
						task_logs = events
					}
				}

				if *c.ExitCode != 0 {
					msg := fmt.Sprintf("Task failed with exit code %d\n", *c.ExitCode)
					task_errors = append(task_errors, errors.New(msg))
				}

				result := &MonitorTaskResult{
					ARN:    task_arn,
					Errors: task_errors,
					Logs:   task_logs,
				}

				result_set[task_arn] = result
				remaining -= 1
			}
		}
	}

	return result_set, nil
}
