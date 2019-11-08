package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/acm"
	"github.com/sfuruya0612/snatch/internal/util"
)

// ACM client struct
type ACM struct {
	Client *acm.ACM
}

// newAcmSess return ACM struct initialized
func NewAcmSess(profile, region string) *ACM {
	return &ACM{
		Client: acm.New(getSession(profile, region)),
	}
}

// Certificate acm certicate struct
type Certificate struct {
	DomainName string
	Type       string
	Status     string
	NotBefore  string
	NotAfter   string
	// InUseBy    string
}

// Certificates Certificate struct slice
type Certificates []Certificate

func (c *ACM) ListCertificates() error {
	certs, err := c.Client.ListCertificates(nil)
	if err != nil {
		return fmt.Errorf("List certificates: %v", err)
	}

	list := Certificates{}
	for _, l := range certs.CertificateSummaryList {

		input := &acm.DescribeCertificateInput{
			CertificateArn: l.CertificateArn,
		}

		output, err := c.Client.DescribeCertificate(input)
		if err != nil {
			return fmt.Errorf("Describe certificate: %v", err)
		}

		cert := output.Certificate

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
