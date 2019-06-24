package url

import (
	"crypto/tls"
	"net/http"
	"time"
)

const timeout = time.Second * 2

var tlsConf = tls.Config{
	InsecureSkipVerify: true,
}

var transport = http.Transport{
	TLSClientConfig: &tlsConf,
}

// Client erstellt einen neuen http Client
func Client() http.Client {
	return http.Client{
		Transport: &transport,
		Timeout:   timeout,
	}
}
