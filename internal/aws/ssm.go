package aws

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Response struct {
	InstanceId string   `json:"instance_id"`
	Status     string   `json:"status"`
	Output     []string `json:"output"`
}

func newSsmSess(profile, region string) *ssm.SSM {
	sess := getSession(profile, region)
	return ssm.New(sess)
}

func StartSession(profile, region string) error {
	client := newSsmSess(profile, region)

	ids, err := listInstanceIds(client)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	// ssm.DescribeInstanceInformation では NameTag が取得できない
	// InstanceId で fileter して ec2.DescribeInstance から取得する
	list, err := getInstancesByInstanceIds(profile, region, ids)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	elements := []string{}
	for _, i := range list {
		item := i.Name + "\t" + i.InstanceId + "\t" + i.LaunchTime
		elements = append(elements, item)
	}

	instance, err := util.Prompt(elements, "Select Instance")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	id := strings.Split(instance, "\t")

	input := &ssm.StartSessionInput{
		Target: aws.String(id[1]),
	}

	sess, endpoint, err := createStartSession(client, input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	sessJson, err := util.Marshal(sess)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	paramsJson, err := util.Marshal(input)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	plug, err := exec.LookPath("session-manager-plugin")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	if err = util.ExecCommand(plug, string(sessJson), region, "StartSession", profile, string(paramsJson), endpoint); err != nil {
		fmt.Println(err)
		err := deleteStartSession(client, *sess.SessionId)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}

func listInstanceIds(client *ssm.SSM) ([]string, error) {
	input := &ssm.DescribeInstanceInformationInput{
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []*string{aws.String("Online")},
			},
		},
	}

	ids := []string{}
	output := func(page *ssm.DescribeInstanceInformationOutput, lastPage bool) bool {
		for _, i := range page.InstanceInformationList {
			ids = append(ids, *i.InstanceId)
		}
		return true
	}

	err := client.DescribeInstanceInformationPages(input, output)
	if err != nil {
		return nil, fmt.Errorf("Describe instance information: %v", err)
	}

	return ids, nil
}

func createStartSession(client *ssm.SSM, input *ssm.StartSessionInput) (*ssm.StartSessionOutput, string, error) {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	sess, err := client.StartSessionWithContext(subctx, input)
	if err != nil {
		return nil, "", err
	}

	return sess, client.Endpoint, nil
}

func deleteStartSession(client *ssm.SSM, sessionId string) error {
	subctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	_, err := client.TerminateSessionWithContext(subctx, &ssm.TerminateSessionInput{SessionId: &sessionId})
	if err != nil {
		return err
	}

	return nil
}

func SendCommand(profile, region, file, id, tag string, args []string) error {
	client := newSsmSess(profile, region)

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

	input := &ssm.SendCommandInput{
		DocumentName:   aws.String("AWS-RunShellScript"),
		MaxConcurrency: aws.String("25%"),
		MaxErrors:      aws.String("0"),
		TimeoutSeconds: aws.Int64(60),
		Parameters:     param,
	}

	if len(id) > 0 {
		input.InstanceIds = []*string{aws.String(id)}
	}

	if len(tag) > 0 {
		spl := strings.Split(tag, ":")
		if len(spl) != 2 {
			return fmt.Errorf("Parse tag=%s", tag)
		}
		input.Targets = []*ssm.Target{
			{
				Key:    aws.String("tag:" + spl[0]),
				Values: []*string{aws.String(spl[1])},
			},
		}
	}

	out, err := client.SendCommand(input)
	if err != nil {
		return fmt.Errorf("Command send: %v", err)
	}

	get := &ssm.ListCommandInvocationsInput{
		CommandId: out.Command.CommandId,
		Details:   aws.Bool(true),
	}

	for {
		got, err := client.ListCommandInvocations(get)
		if err != nil {
			return fmt.Errorf("List command invocation: %v", err)
		}

		if len(got.CommandInvocations) == 0 {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		inprogress := false
		for _, ci := range got.CommandInvocations {
			if *ci.Status == "InProgress" {
				inprogress = true
				break
			}
		}

		if inprogress {
			time.Sleep(100 * time.Millisecond)
			continue
		}

		resp := []Response{}
		for _, ci := range got.CommandInvocations {
			out := *ci.CommandPlugins[0].Output
			spl := strings.Split(out, "\n")
			if len(spl[len(spl)-1]) < 1 {
				spl = spl[:len(spl)-1]
			}

			res := Response{
				InstanceId: *ci.InstanceId,
				Status:     *ci.Status,
				Output:     spl,
			}
			resp = append(resp, res)

			if len(out) < 2500 {
				continue
			}

			res.Output = spl
		}
		json, err := util.JParser(resp)
		if err != nil {
			return fmt.Errorf("Json Marshal: %v", err)
		}

		for r := range json {
			fmt.Printf("\n\x1b[35mInstance_id:\x1b[0m %v \x1b[35mStatus:\x1b[0m %v\n\x1b[35mOutput:\x1b[0m\n", json[r].Instance_id, json[r].Status)

			for _, o := range json[r].Output {
				fmt.Printf("%v\n", o)
			}
		}

		break
	}

	return nil
}
