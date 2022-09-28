package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"

	saws "github.com/sfuruya0612/snatch/internal/aws"
	"github.com/sfuruya0612/snatch/internal/util"
	"github.com/urfave/cli/v2"
)

var Ssm = &cli.Command{
	Name:  "ssm",
	Usage: "Start a session on your instances by launching bash or shell terminal",
	Action: func(c *cli.Context) error {
		return startSession(c.String("profile"), c.String("region"))
	},
	Subcommands: []*cli.Command{
		{
			Name:      "command",
			Aliases:   []string{"cmd"},
			Usage:     "Runs shell script to target instances",
			ArgsUsage: "[ --tag | -t ] <Key:Value> [ --id | -i ] <InstanceId> [ --file | -f ] <ScriptFile>",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "tag",
					Aliases: []string{"t"},
					Usage:   "Set Key-Value of the tag (e.g. -t Name:test-ec2)",
				},
				&cli.StringFlag{
					Name:    "id",
					Aliases: []string{"i"},
					Usage:   "Set EC2 instance id",
				},
				&cli.StringFlag{
					Name:    "file",
					Aliases: []string{"f"},
					Usage:   "Set execute file",
				},
			},
			Action: func(c *cli.Context) error {
				return sendCommand(c.String("profile"), c.String("region"), c.String("tag"), c.String("id"), c.String("file"), c.Args())
			},
		},
		{
			Name:    "parameter",
			Aliases: []string{"param"},
			Usage:   "Get parameter store",
			Action: func(c *cli.Context) error {
				return startSession(c.String("profile"), c.String("region"))
			},
		},
	},
}

func startSession(profile, region string) error {
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
		InstanceIds: ids,
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

func sendCommand(profile, region, tag, id, file string, args cli.Args) error {
	if args.Len() == 0 && len(file) == 0 {
		return fmt.Errorf("args or file is required")
	}

	if len(id) == 0 && len(tag) == 0 {
		return fmt.Errorf("instance id or tag is required")
	}

	param := make(map[string][]*string)

	if args.Len() > 0 {
		param["commands"] = []*string{
			aws.String(args.Get(0)),
		}
	}

	if len(file) > 0 {
		f, err := os.Open(file)
		if err != nil {
			return fmt.Errorf("open file %s: %v", file, err)
		}
		// defer f.Close()

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
			return fmt.Errorf("parse tag=%s", tag)
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
	client := saws.NewSsmSess(profile, region)

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
