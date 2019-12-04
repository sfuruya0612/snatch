package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/ssm"
	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli"
)

func StartSession(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	ssmclient := saws.NewSsmSess(profile, region)

	input := &ssm.DescribeInstanceInformationInput{
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []*string{aws.String("Online")},
			},
		},
	}

	ids, err := ssmclient.DescribeInstanceInformation(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ec2client := saws.NewEc2Sess(profile, region)

	ec2input := &ec2.DescribeInstancesInput{
		InstanceIds: aws.StringSlice(ids),
	}

	// ssm.DescribeInstanceInformation では NameTag が取得できない
	// InstanceId で fileter して ec2.DescribeInstances から取得する
	list, err := ec2client.DescribeInstances(ec2input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ec2list := []string{}
	for _, i := range list {
		item := i.Name + "\t" + i.InstanceId + "\t" + i.LaunchTime
		ec2list = append(ec2list, item)
	}

	instance, err := util.Prompt(ec2list, "Select Instance")
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	s := strings.Split(instance, "\t")

	si := &ssm.StartSessionInput{
		Target: aws.String(s[1]),
	}

	sess, endpoint, err := ssmclient.CreateStartSession(si)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	sessJson, err := util.Marshal(sess)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	paramsJson, err := util.Marshal(si)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	plug, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ti := &ssm.TerminateSessionInput{
		SessionId: aws.String(*sess.SessionId),
	}

	if err = util.ExecCommand(plug, string(sessJson), region, "StartSession", profile, string(paramsJson), endpoint); err != nil {
		fmt.Println(err)
		if err := ssmclient.DeleteStartSession(ti); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}

func GetSsmHist(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")
	flag := c.Bool("active")

	state := "History"
	if flag {
		state = "Active"
	}

	input := &ssm.DescribeSessionsInput{
		State: aws.String(state),
	}

	client := saws.NewSsmSess(profile, region)

	hist, err := client.DescribeSessions(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintSessHist(os.Stdout, hist); err != nil {
		return fmt.Errorf("Failed to print resources")
	}

	return nil
}

func SendCommand(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	args := c.Args()
	file := c.String("file")
	if len(args) == 0 && len(file) == 0 {
		return fmt.Errorf("Args or file is required")
	}

	id := c.String("instanceid")
	tag := c.String("tag")
	if len(id) == 0 && len(tag) == 0 {
		return fmt.Errorf("Instance id or tag is required")
	}

	param := make(map[string][]*string)

	if len(args) > 0 {
		command := []*string{
			aws.String(args[0]),
		}
		param["commands"] = command
	}

	if len(file) > 0 {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("open file %s: %v", file, err)
		}
		defer f.Close()

		command := []*string{}
		s := bufio.NewScanner(f)
		for s.Scan() {
			command = append(command, aws.String(s.Text()))
		}
		param["commands"] = command
	}

	ci := &ssm.SendCommandInput{
		DocumentName:   aws.String("AWS-RunShellScript"),
		MaxConcurrency: aws.String("25%"),
		MaxErrors:      aws.String("0"),
		TimeoutSeconds: aws.Int64(60),
		Parameters:     param,
	}

	if len(id) > 0 {
		ci.InstanceIds = []*string{aws.String(id)}
	}

	if len(tag) > 0 {
		spl := strings.Split(tag, ":")
		if len(spl) != 2 {
			return fmt.Errorf("Parse tag=%s", tag)
		}
		ci.Targets = []*ssm.Target{
			{
				Key:    aws.String("tag:" + spl[0]),
				Values: []*string{aws.String(spl[1])},
			},
		}
	}

	client := saws.NewSsmSess(profile, region)

	comm, err := client.SendCommand(ci)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ii := &ssm.ListCommandInvocationsInput{
		CommandId: comm.Command.CommandId,
		Details:   aws.Bool(true),
	}

	invo, err := client.ListCommandInvocations(ii)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	json, err := util.JParser(invo)
	if err != nil {
		return fmt.Errorf("Json Marshal: %v", err)
	}

	for r := range json {
		fmt.Printf("\n\x1b[35mInstance_id:\x1b[0m %v \x1b[35mStatus:\x1b[0m %v\n\x1b[35mOutput:\x1b[0m\n", json[r].Instance_id, json[r].Status)

		for _, o := range json[r].Output {
			fmt.Printf("%v\n", o)
		}
	}

	return nil
}

func GetCmdLog(c *cli.Context) error {
	profile := c.GlobalString("profile")
	region := c.GlobalString("region")

	client := saws.NewSsmSess(profile, region)

	input := &ssm.ListCommandsInput{
		MaxResults: aws.Int64(30),
	}

	logs, err := client.ListCommands(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintCmdLogs(os.Stdout, logs); err != nil {
		return fmt.Errorf("Failed to print command logs")
	}

	return nil
}
