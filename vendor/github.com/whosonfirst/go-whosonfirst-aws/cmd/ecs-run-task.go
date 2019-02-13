package main

/*

this is the simplest-dumbest tool to run and ECS task and log its output - it works but could use
a lot of finessing... (20190205/thisisaaronland)

./bin/ecs-run-task -monitor -ecs-dsn 'region=us-east-1 credentials=session' -cluster example -container example -task 'example:1' -subnet 'subnet-example' -security-group 'sg-example' curl -s localhost:9200
time passes...
2019/02/05 12:58:31 Task arn:aws:ecs:us-east-1:xxxx:task/example failed with exit code 7

see the way the (task) output isn't included? that's one of those details to finesse...
(20190205/thisisaaronland)

*/

import (
	"errors"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/whosonfirst/go-whosonfirst-aws/session"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"os"
	"strings"
)

func main() {

	var ecs_dsn = flag.String("ecs-dsn", "", "A valid (go-whosonfirst-aws) ECS DSN.")
	var cw_dsn = flag.String("cloudwatch-dsn", "", "A valid (go-whosonfirst-aws) CloudWatch DSN.")

	var container = flag.String("container", "", "The name of your AWS ECS container.")
	var cluster = flag.String("cluster", "", "The name of your AWS ECS cluster.")
	var task = flag.String("task", "", "The name of your AWS ECS task (inclusive of its version number),")

	var launch_type = flag.String("launch-type", "FARGATE", "...")
	var public_ip = flag.String("public-ip", "ENABLED", "...")

	var monitor = flag.Bool("monitor", false, "...")

	var subnets flags.MultiString
	flag.Var(&subnets, "subnet", "One or more AWS subnets in which your task will run.")

	var security_groups flags.MultiString
	flag.Var(&security_groups, "security-group", "One of more AWS security groups your task will assume.")

	flag.Parse()

	ecs_sess, err := session.NewSessionWithDSN(*ecs_dsn)

	if err != nil {
		log.Fatal(err)
	}

	if *cw_dsn == "" {
		*cw_dsn = *ecs_dsn
	}

	cw_sess, err := session.NewSessionWithDSN(*cw_dsn)

	if err != nil {
		log.Fatal(err)
	}

	ecs_svc := ecs.New(ecs_sess)
	cw_svc := cloudwatchlogs.New(cw_sess)

	ecs_cluster := aws.String(*cluster)
	ecs_task := aws.String(*task)

	ecs_launch_type := aws.String(*launch_type)
	ecs_public_ip := aws.String(*public_ip)

	ecs_cmd := make([]*string, len(flag.Args()))

	for i, fl := range flag.Args() {
		ecs_cmd[i] = aws.String(fl)
	}

	// either this doesn't work or I am doing it wrong...
	// ecs_cmd = append(ecs_cmd, aws.String("2>&1"))
	
	ecs_subnets := make([]*string, len(subnets))
	ecs_security_groups := make([]*string, len(security_groups))

	for i, sn := range subnets {
		ecs_subnets[i] = aws.String(sn)
	}

	for i, sg := range security_groups {
		ecs_security_groups[i] = aws.String(sg)
	}

	ecs_network := &ecs.NetworkConfiguration{
		AwsvpcConfiguration: &ecs.AwsVpcConfiguration{
			AssignPublicIp: ecs_public_ip,
			SecurityGroups: ecs_security_groups,
			Subnets:        ecs_subnets,
		},
	}

	ecs_process_override := &ecs.ContainerOverride{
		Name:    aws.String(*container),
		Command: ecs_cmd,
	}

	ecs_overrides := &ecs.TaskOverride{
		ContainerOverrides: []*ecs.ContainerOverride{
			ecs_process_override,
		},
	}

	req := &ecs.RunTaskInput{
		Cluster:              ecs_cluster,
		TaskDefinition:       ecs_task,
		LaunchType:           ecs_launch_type,
		NetworkConfiguration: ecs_network,
		Overrides:            ecs_overrides,
	}

	rsp, err := ecs_svc.RunTask(req)

	if err != nil {
		log.Fatal(err)
	}

	if !*monitor {

		for _, t := range rsp.Tasks {
			fmt.Println(*t.TaskArn)
		}

		os.Exit(0)
	}

	count_tasks := len(rsp.Tasks)
	remaining := count_tasks

	ecs_tasks := make([]*string, count_tasks)

	for i, t := range rsp.Tasks {
		ecs_tasks[i] = t.TaskArn
	}

	task_errors := make([]error, 0)

	for remaining > 0 {

		monitor_req := &ecs.DescribeTasksInput{
			Cluster: aws.String(*cluster),
			Tasks:   ecs_tasks,
		}

		monitor_rsp, err := ecs_svc.DescribeTasks(monitor_req)

		if err != nil {
			log.Fatal(err)
		}

		for _, t := range monitor_rsp.Tasks {

			for _, c := range t.Containers {

				if *c.Name != *container {
					continue
				}

				if *c.LastStatus != "STOPPED" {
					continue
				}

				// start of generic code to put in a function
				// TO DO: what if the logs haven't reached CW yet... ?

				arn := strings.Split(*t.TaskArn, "/")

				cw_group := fmt.Sprintf("/ecs/%s", *container)
				cw_stream := fmt.Sprintf("ecs/%s/%s", *container, arn[1])

				cw_req := &cloudwatchlogs.GetLogEventsInput{
					LogGroupName:  aws.String(cw_group),
					LogStreamName: aws.String(cw_stream),
					StartFromHead: aws.Bool(true),
				}

				cw_rsp, err := cw_svc.GetLogEvents(cw_req)

				if err == nil {

					for _, e := range cw_rsp.Events {
						log.Printf("[%s][%d] %s\n", *t.TaskArn, *e.Timestamp, *e.Message)
					}
				}

				// TODO: paginated logs...
				// end of generic code to put in a function

				if *c.ExitCode != 0 {
					msg := fmt.Sprintf("Task %s failed with exit code %d\n", *t.TaskArn, *c.ExitCode)
					err := errors.New(msg)
					task_errors = append(task_errors, err)
				}

				remaining -= 1
			}
		}
	}

	if len(task_errors) > 0 {

		for _, e := range task_errors {
			log.Println(e)
		}

		os.Exit(1)
	}

	os.Exit(0)
}
