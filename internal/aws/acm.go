package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/sfuruya0612/snatch/internal/util"
)

type Certificate struct {
	DomainName string
	Type       string
	Status     string
	NotBefore  string
	NotAfter   string
	// InUseBy    string
}

type Certificates []Certificate

func NewAcmSess(profile string, region string) *acm.ACM {
	sess := getSession(profile, region)
	return acm.New(sess)
}

func ListCertificates(profile string, region string) error {
	client := NewAcmSess(profile, region)

	res, err := client.ListCertificates(nil)
	if err != nil {
		return fmt.Errorf("List certificates: %v", err)
	}

	list := Certificates{}
	for _, l := range res.CertificateSummaryList {

		input := &acm.DescribeCertificateInput{
			CertificateArn: l.CertificateArn,
		}

		res, err := client.DescribeCertificate(input)
		if err != nil {
			return fmt.Errorf("Describe certificate: %v", err)
		}

		cert := res.Certificate

		before := cert.NotBefore.String()
		after := cert.NotAfter.String()

		list = append(list, Certificate{
			DomainName: *cert.DomainName,
			Type:       *cert.Type,
			Status:     *cert.Status,
			NotBefore:  before,
			NotAfter:   after,
		})
	}
	f := util.Formatln(
		list.DomainName(),
		list.Type(),
		list.Status(),
		list.NotBefore(),
		list.NotAfter(),
	)

	for _, i := range list {
		fmt.Printf(
			f,
			i.DomainName,
			i.Type,
			i.Status,
			i.NotBefore,
			i.NotAfter,
		)
	}

	return nil
}

func (cert Certificates) DomainName() []string {
	dname := []string{}
	for _, i := range cert {
		dname = append(dname, i.DomainName)
	}
	return dname
}

func (cert Certificates) Type() []string {
	ty := []string{}
	for _, i := range cert {
		ty = append(ty, i.Type)
	}
	return ty
}

func (cert Certificates) Status() []string {
	status := []string{}
	for _, i := range cert {
		status = append(status, i.Status)
	}
	return status
}

func (cert Certificates) NotBefore() []string {
	before := []string{}
	for _, i := range cert {
		before = append(before, i.NotBefore)
	}
	return before
}

func (cert Certificates) NotAfter() []string {
	after := []string{}
	for _, i := range cert {
		after = append(after, i.NotAfter)
	}
	return after
}
