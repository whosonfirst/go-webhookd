package main

import (
	_ "context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/whosonfirst/go-whosonfirst-aws/ecs"
	"github.com/whosonfirst/go-whosonfirst-cli/flags"
	"log"
	"os"
)

func main() {

	var ecs_dsn = flag.String("dsn", "", "A valid (go-whosonfirst-aws) ECS DSN.")

	var ecs_container = flag.String("container", "", "The name of your AWS ECS container.")
	var ecs_cluster = flag.String("cluster", "", "The name of your AWS ECS cluster.")
	var ecs_task = flag.String("task", "", "The name of your AWS ECS task (inclusive of its version number),")

	var launch_type = flag.String("launch-type", "FARGATE", "...")
	var public_ip = flag.String("public-ip", "ENABLED", "...")

	var monitor = flag.Bool("monitor", false, "...")

	var logs = flag.Bool("logs", false, "...")
	var logs_dsn = flag.String("logs-dsn", "", "A valid (go-whosonfirst-aws) CloudWatchLogs DSN.")

	var subnets flags.MultiString
	flag.Var(&subnets, "subnet", "One or more AWS subnets in which your task will run.")

	var security_groups flags.MultiString
	flag.Var(&security_groups, "security-group", "One of more AWS security groups your task will assume.")

	flag.Parse()

	if *logs == true {
		*monitor = true
	}

	if *logs_dsn == "" {
		*logs_dsn = *ecs_dsn
	}

	task_opts := &ecs.TaskOptions{
		DSN:            *ecs_dsn,
		Task:           *ecs_task,
		Container:      *ecs_container,
		Cluster:        *ecs_cluster,
		Subnets:        subnets,
		SecurityGroups: security_groups,
		LaunchType:     *launch_type,
		PublicIP:       *public_ip,
	}

	cmd := flag.Args()

	task_rsp, err := ecs.LaunchTask(task_opts, cmd...)

	if err != nil {
		log.Fatal(err)
	}

	if *monitor {

		monitor_opts := &ecs.MonitorTaskOptions{
			DSN:       *ecs_dsn,
			Container: *ecs_container,
			Cluster:   *ecs_cluster,
			WithLogs:  *logs,
			LogsDSN:   *logs_dsn,
		}

		monitor_rsp, err := ecs.MonitorTasks(monitor_opts, task_rsp.Tasks...)

		if err != nil {
			log.Fatal(err)
		}

		enc_rsp, err := json.Marshal(monitor_rsp)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(enc_rsp))
	}

	os.Exit(0)
}
