package domain

import (
	"crypto/tls"
	"crypto/x509"
	"math"
	"net"
	"reflect"
	"strings"
	"time"
)

type errorStr struct {
	s string
}

func (e *errorStr) Error() string {
	return e.s
}

// DomainData holds additional data about a domain, also used for JSON representation
type DomainData struct {
	Name     string    `json:"name"`
	DaysLeft int       `json:"daysLeft"`
	EndTime  time.Time `json:"endTime"`
	Status   string    `json:"status"`
}

// InitDomain inits Domain and get the cert
func InitDomain(domain string) (*Domain, error) {
	obj := Domain{
		Name: domain,
	}

	err := obj.getCertifcate()

	return &obj, err
}

// Domain holds a domain and it's certificate
type Domain struct {
	Name        string
	Certificate *x509.Certificate
	Error       error
}

// getCertificate gets the certificate
func (domain *Domain) getCertifcate() error {
	conn, err := tls.Dial("tcp", domain.Name+":443", nil)

	if err != nil {
		domain.Error = err
		return err
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) > 0 {
		domain.Certificate = certs[0]
	} else {
		err = &errorStr{"no certificate"}
		domain.Error = err
		return err
	}

	return nil
}

// EndTime returns the last date the certificate is valid
func (domain *Domain) EndTime(l *time.Location) time.Time {
	if domain.Certificate != nil {
		return domain.Certificate.NotAfter.In(l)
	}

	return time.Time{}
}

// DaysLeft returns numer of days left (floored)
func (domain *Domain) DaysLeft(l *time.Location) int {
	now := time.Now().In(l)
	past := domain.EndTime(l)
	diff := past.Sub(now)
	return int(math.Floor(diff.Hours() / 24))
}

// Status returns ok or any errors
func (domain *Domain) Status() string {
	if domain.Error == nil {
		return "ok"
	}

	if reflect.TypeOf(domain.Error) == reflect.TypeOf(&net.OpError{}) {
		parts := strings.Split(domain.Error.Error(), ":")
		if len(parts) > 0 {
			return strings.TrimSpace(parts[len(parts)-1])
		}
	}
	return domain.Error.Error()
}

// GetData returns all data as an struct
func (domain *Domain) GetData(l *time.Location) *DomainData {
	obj := &DomainData{
		Name:     domain.Name,
		DaysLeft: domain.DaysLeft(l),
		EndTime:  domain.EndTime(l),
		Status:   domain.Status(),
	}
	return obj
}
