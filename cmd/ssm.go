package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	ec2Types "github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	ssmTypes "github.com/aws/aws-sdk-go-v2/service/ssm/types"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli/v2"
)

var Ssm = &cli.Command{
	Name:  "ssm",
	Usage: "Use Systems Manager services",
	Subcommands: []*cli.Command{
		{
			Name:    "parameter",
			Aliases: []string{"p"},
			Usage:   "Get parameter store",
			Action: func(c *cli.Context) error {
				return getParameter(c.String("profile"), c.String("region"))
			},
		},
	},
}

func startSession(profile, region string) error {
	ssmclient := saws.NewSsmClient(profile, region)

	input := &ssm.DescribeInstanceInformationInput{
		Filters: []ssmTypes.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []string{"Online"},
			},
		},
	}

	ids, err := ssmclient.DescribeInstanceInformation(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ec2client := saws.NewEc2Client(profile, region)

	ec2input := &ec2.DescribeInstancesInput{
		InstanceIds: ids,
		Filters: []ec2Types.Filter{
			{
				Name: aws.String("instance-state-name"),
				Values: []string{
					"running",
				},
			},
		},
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

	sess, err := ssmclient.StartSession(si)
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

	if err = util.ExecCommand(plug, string(sessJson), region, "StartSession", profile, string(paramsJson), fmt.Sprintf("https://ssm.%s.amazonaws.com", region)); err != nil {
		fmt.Println(err)
		if err := ssmclient.DeleteSession(ti); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}

func sendCommand(profile, region, tag, id, file string, args cli.Args) error {
	if args.Len() == 0 && len(file) == 0 {
		return fmt.Errorf("args or file is required")
	}

	if len(id) == 0 && len(tag) == 0 {
		return fmt.Errorf("instance id or tag is required")
	}

	param := make(map[string][]string)

	if args.Len() > 0 {
		param["commands"] = []string{
			args.Get(0),
		}
	}

	if len(file) > 0 {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("open file %s: %v", file, err)
		}
		// defer f.Close()

		command := []string{}
		s := bufio.NewScanner(f)
		for s.Scan() {
			command = append(command, s.Text())
		}
		param["commands"] = command
	}

	ci := &ssm.SendCommandInput{
		DocumentName:   aws.String("AWS-RunShellScript"),
		MaxConcurrency: aws.String("25%"),
		MaxErrors:      aws.String("0"),
		TimeoutSeconds: aws.Int32(60),
		Parameters:     param,
	}

	if len(id) > 0 {
		ci.InstanceIds = []string{id}
	}

	if len(tag) > 0 {
		spl := strings.Split(tag, ":")
		if len(spl) != 2 {
			return fmt.Errorf("parse tag=%s", tag)
		}
		ci.Targets = []ssmTypes.Target{
			{
				Key:    aws.String("tag:" + spl[0]),
				Values: []string{spl[1]},
			},
		}
	}

	client := saws.NewSsmClient(profile, region)

	comm, err := client.SendCommand(ci)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	ii := &ssm.ListCommandInvocationsInput{
		CommandId: comm.Command.CommandId,
		Details:   true,
	}

	invo, err := client.ListCommandInvocations(ii)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	json, err := util.JParser(invo)
	if err != nil {
		return fmt.Errorf("json marshal: %v", err)
	}

	for r := range json {
		fmt.Printf("\n\x1b[35mInstance_id:\x1b[0m %v \x1b[35mStatus:\x1b[0m %v\n\x1b[35mOutput:\x1b[0m\n", json[r].Instance_id, json[r].Status)

		for _, o := range json[r].Output {
			fmt.Printf("%v\n", o)
		}
	}

	return nil
}

func getParameter(profile, region string) error {
	client := saws.NewSsmClient(profile, region)

	params, err := client.DescribeParameters(&ssm.DescribeParametersInput{})
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	param, err := client.GetParameter(params)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err := saws.PrintParameters(os.Stdout, param); err != nil {
		return fmt.Errorf("failed to print parameters")
	}

	return nil
}
