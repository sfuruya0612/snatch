package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
)

var ec2Cmd = &cobra.Command{
	Use:   "ec2",
	Short: "Manage EC2",
}

var ec2ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "List EC2 instances",
	Run:   listEC2Instances,
}

var ec2StartSessionCmd = &cobra.Command{
	Use:     "session",
	Aliases: []string{"s"},
	Short:   "Start a session to an EC2 instance",
	Run:     startSession,
}

var ec2Columns = []util.Column{
	{Header: "Name", Width: 50},
	{Header: "InstanceID", Width: 20},
	{Header: "InstanceType", Width: 12},
	{Header: "Lifecycle", Width: 9},
	{Header: "PrivateIP", Width: 12},
	{Header: "PublicIP", Width: 14},
	{Header: "State", Width: 10},
	{Header: "KeyName", Width: 20},
	{Header: "AZ", Width: 10},
	{Header: "LaunchTime", Width: 30},
}

func listEC2Instances(cmd *cobra.Command, args []string) {
	client := aws.NewEC2Client(cmd.Flag("profile").Value.String(), cmd.Flag("region").Value.String())

	opts := &aws.Ec2Opts{
		Running: cmd.Flag("running").Value.String() == "true",
	}

	input, err := aws.GenerateDescribeInstancesInput(opts)
	if err != nil {
		fmt.Printf("Unable to generate describe instances input: %v\n", err)
		return
	}

	instances, err := aws.DescribeInstances(client, input)
	if err != nil {
		fmt.Printf("Unable to describe instances: %v\n", err)
		return
	}

	if len(instances) == 0 {
		fmt.Println("No EC2 instances found")
		return
	}

	formatter := util.NewTableFormatter(ec2Columns, cmd.Flag("output").Value.String())
	formatter.PrintHeader()
	formatter.PrintRows(instances)
}

func startSession(cmd *cobra.Command, args []string) {
	profile := cmd.Flag("profile").Value.String()
	region := cmd.Flag("region").Value.String()

	ssmClient := aws.NewSsmClient(profile, region)

	ssmOpts := &aws.SsmOpts{
		PingStatus:   "Online",
		ResourceType: "EC2Instance",
		InstanceId:   cmd.Flag("instance-id").Value.String(),
	}

	if cmd.Flag("instance-id").Value.String() == "" {
		ssmInput := aws.GenerateDescribeInstanceInformationInput(ssmOpts)

		instanceIds, err := aws.DescribeInstanceInformation(ssmClient, ssmInput)
		if err != nil {
			fmt.Printf("Unable to describe instance information: %v\n", err)
			return
		}

		instances, err := ssmTargetInstance(profile, region, instanceIds)
		if err != nil {
			fmt.Printf("Unable to get target instance: %v\n", err)
			return
		}

		selected, err := util.Select(instances, "Select an EC2 instance:")
		if err != nil {
			fmt.Printf("Error selecting instance: %v\n", err)
			return
		}

		ssmOpts.InstanceId = selected.ID()
	}

	startSessionInput := aws.GenerateStartSessionInput(ssmOpts)

	session, err := aws.StartSession(ssmClient, startSessionInput)
	if err != nil {
		fmt.Printf("Unable to start session: %v\n", err)
		return
	}

	sessJson, err := parser(session)
	if err != nil {
		fmt.Printf("Unable to marshal session: %v\n", err)
		return
	}

	paramsJson, err := parser(startSessionInput)
	if err != nil {
		fmt.Printf("Unable to marshal start session input: %v\n", err)
		return
	}

	plug, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		fmt.Println("Unable to find session-manager-plugin")
		return
	}

	ssmOpts.SessionId = *session.SessionId

	terminateSessionInput := aws.GenerateTerminateSessionInput(ssmOpts)

	if err = execCommand(plug, string(sessJson), region, "StartSession", profile, string(paramsJson), fmt.Sprintf("https://ssm.%s.amazonaws.com", region)); err != nil {
		fmt.Println(err)
		_, err := aws.TerminateSession(ssmClient, terminateSessionInput)
		if err != nil {
			fmt.Printf("Unable to terminate session: %v\n", err)
			return
		}
	}
}

func ssmTargetInstance(profile, region string, instanceIds []string) ([]util.Item, error) {
	ec2Client := aws.NewEC2Client(profile, region)

	list := [][]string{}
	for _, i := range instanceIds {
		ec2Opts := &aws.Ec2Opts{
			Running:    true,
			InstanceId: i,
		}

		ec2Input, err := aws.GenerateDescribeInstancesInput(ec2Opts)
		if err != nil {
			return nil, fmt.Errorf("Unable to generate describe instances input: %v", err)
		}

		instances, err := aws.DescribeInstances(ec2Client, ec2Input)
		if err != nil {
			return nil, fmt.Errorf("Unable to describe instances: %v", err)
		}

		list = append(list, instances[0])
	}

	if len(list) == 0 {
		return nil, fmt.Errorf("No EC2 instances found")
	}

	var instances []util.Item
	for _, inst := range list {
		instances = append(instances, aws.EC2Instance{
			Name:       inst[0],
			InstanceID: inst[1],
		})
	}

	return instances, nil
}

func parser(i interface{}) ([]byte, error) {
	bytes, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("json Marshal error: %v", err)
	}

	return bytes, nil
}

func execCommand(process string, args ...string) error {
	call := exec.Command(process, args...)
	call.Stderr = os.Stderr
	call.Stdout = os.Stdout
	call.Stdin = os.Stdin

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	done := make(chan bool, 1)
	go func() {
		for {
			select {
			case <-sigs:
			case <-done:
				// break
			}
		}
	}()
	defer close(done)

	if err := call.Run(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}
