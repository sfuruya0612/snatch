package aws

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/sfuruya0612/snatch/internal/util"
)

type SsmInstance struct {
	ComputerName    string
	InstanceID      string
	IPAddress       string
	AgentVersion    string
	PlatformName    string
	PlatformVersion string
}

type SsmInstances []SsmInstance

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

	list, err := ListInstances(client)
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	elements := []string{}
	for _, i := range list {

		var item string
		// item = i.InstanceID + "\t" + i.IPAddress
		item = i.InstanceID

		elements = append(elements, item)
	}

	instance, err := util.Prompt(elements, "Select Instance")
	if err != nil {
		return fmt.Errorf("%v", err)
	}

	input := &ssm.StartSessionInput{
		Target: aws.String(instance),
	}

	token, err := client.StartSession(input)
	fmt.Println(token)

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
		MaxErrors:      aws.String("1"),
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

func ListInstances(client *ssm.SSM) (SsmInstances, error) {
	input := &ssm.DescribeInstanceInformationInput{
		MaxResults: aws.Int64(50),
		Filters: []*ssm.InstanceInformationStringFilter{
			{
				Key:    aws.String("PingStatus"),
				Values: []*string{aws.String("Online")},
			},
		},
	}

	instances, err := client.DescribeInstanceInformation(input)
	if err != nil {
		return nil, fmt.Errorf("Describe information: %v", err)
	}

	list := SsmInstances{}
	for _, i := range instances.InstanceInformationList {

		list = append(list, SsmInstance{
			ComputerName:    *i.ComputerName,
			InstanceID:      *i.InstanceId,
			IPAddress:       *i.IPAddress,
			AgentVersion:    *i.AgentVersion,
			PlatformName:    *i.PlatformName,
			PlatformVersion: *i.PlatformVersion,
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].ComputerName < list[j].ComputerName
	})

	return list, nil
}

func (ssmi SsmInstances) ComputerName() []string {
	cname := []string{}
	for _, i := range ssmi {
		cname = append(cname, i.ComputerName)
	}
	return cname
}

func (ssmi SsmInstances) InstanceID() []string {
	id := []string{}
	for _, i := range ssmi {
		id = append(id, i.InstanceID)
	}
	return id
}

func (ssmi SsmInstances) IPAddress() []string {
	ip := []string{}
	for _, i := range ssmi {
		ip = append(ip, i.IPAddress)
	}
	return ip
}

func (ssmi SsmInstances) AgentVersion() []string {
	ver := []string{}
	for _, i := range ssmi {
		ver = append(ver, i.AgentVersion)
	}
	return ver
}

func (ssmi SsmInstances) PlatformName() []string {
	pname := []string{}
	for _, i := range ssmi {
		pname = append(pname, i.PlatformName)
	}
	return pname
}

func (ssmi SsmInstances) PlatformVersion() []string {
	pver := []string{}
	for _, i := range ssmi {
		pver = append(pver, i.PlatformVersion)
	}
	return pver
}
