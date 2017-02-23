package uaac

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	"github.build.ge.com/adoption/predix-cli/cli/global"
	"github.com/PredixDev/go-uaa-lib"
)

func NewSimpleClient(id string, secret string) *lib.Client {
	return &lib.Client{
		ID:          id,
		Secret:      secret,
		Scopes:      []string{"uaa.none", "openid"},
		GrantTypes:  []string{"authorization_code", "client_credentials", "refresh_token", "password"},
		Authorities: []string{"openid", "uaa.none", "uaa.resource"},
		AutoApprove: []string{"openid"},
	}
}

func NewHTTPClient(tlsConfig *tls.Config) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
			Proxy:           http.ProxyFromEnvironment,
		},
	}
}

func NewTLSConfig(skipSslValidation bool, caCertFile string) *tls.Config {
	config := &tls.Config{
		InsecureSkipVerify: skipSslValidation,
	}
	if !skipSslValidation && caCertFile != "" {
		caCert, err := ioutil.ReadFile(caCertFile)
		if err != nil {
			global.UI.Failed(err.Error())
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)
		config.RootCAs = caCertPool
	}
	return config
}
