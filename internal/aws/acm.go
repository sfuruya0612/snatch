package aws

import (
	"fmt"
	"io"
	"strings"
	"text/tabwriter"

	"github.com/aws/aws-sdk-go/service/acm"
)

// ACM client struct
type ACM struct {
	Client *acm.ACM
}

// NewAcmSess return ACM struct initialized
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

// ListCertificates return Certificates
// input acm.ListCertificatesInput
func (c *ACM) ListCertificates(input *acm.ListCertificatesInput) (Certificates, error) {
	certs, err := c.Client.ListCertificates(input)
	if err != nil {
		return nil, fmt.Errorf("List certificates: %v", err)
	}

	list := Certificates{}
	for _, l := range certs.CertificateSummaryList {

		i := &acm.DescribeCertificateInput{
			CertificateArn: l.CertificateArn,
		}

		output, err := c.Client.DescribeCertificate(i)
		if err != nil {
			return nil, fmt.Errorf("Describe certificate: %v", err)
		}

		cert := output.Certificate

		list = append(list, Certificate{
			DomainName: *cert.DomainName,
			Type:       *cert.Type,
			Status:     *cert.Status,
			NotBefore:  cert.NotBefore.String(),
			NotAfter:   cert.NotAfter.String(),
		})
	}
	if len(list) == 0 {
		return nil, fmt.Errorf("No resources")
	}

	return list, nil
}

func PrintCertificates(wrt io.Writer, resources Certificates) error {
	w := tabwriter.NewWriter(wrt, 0, 8, 1, ' ', 0)
	header := []string{
		"DomainName",
		"Type",
		"Status",
		"NotBefore",
		"NotAfter",
	}

	if _, err := fmt.Fprintln(w, strings.Join(header, "\t")); err != nil {
		return fmt.Errorf("%v", err)
	}

	for _, r := range resources {
		if _, err := fmt.Fprintln(w, r.CertificateTabString()); err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	if err := w.Flush(); err != nil {
		return fmt.Errorf("%v", err)
	}

	return nil
}

func (i *Certificate) CertificateTabString() string {
	fields := []string{
		i.DomainName,
		i.Type,
		i.Status,
		i.NotBefore,
		i.NotAfter,
	}

	return strings.Join(fields, "\t")
}
