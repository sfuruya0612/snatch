package aws

import (
	"fmt"
	"sort"

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
