package cert

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"

	"github.com/kenshaw/envcfg"
)

// Config is the struct that used to store the config file
type Config struct {
	FinacleCertPool Cert
}

// Cert is the struct wrapper which contains cert pool and flag to allow skip read the cert or not
type Cert struct {
	AllowSkip bool
	Pool      *x509.CertPool
}

// New init cert config file
func New(c *envcfg.Envcfg) (*Cert, error) {
	//cert config
	cert := Cert{
		AllowSkip: true,
	}

	if c.Env() != "development" {
		certPool := x509.NewCertPool()
		pem, err := ioutil.ReadFile(c.GetKey("finacle.certpath"))
		if err != nil {
			fmt.Printf("error load cert file ERR: %+v\n", err)
			return nil, err
		}
		certPool.AppendCertsFromPEM(pem)

		cert = Cert{
			AllowSkip: false,
			Pool:      certPool,
		}
	}

	return &cert, nil
}
