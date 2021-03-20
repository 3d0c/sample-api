package rpc

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

const (
	Default_Connection_Timeout = 1 * time.Second
	Default_RW_Timeout         = 10 * time.Second
)

type Config struct {
	ConnectTimeout time.Duration
	RWTimeout      time.Duration
	Headers        http.Header
	CertPEM        []byte
	KeyPEM         []byte
	CaCertPEM      []byte
}

func dialer(config *Config) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		if config == nil || (config.ConnectTimeout == 0 && config.RWTimeout == 0) {
			config = &Config{
				ConnectTimeout: Default_Connection_Timeout,
				RWTimeout:      Default_RW_Timeout,
			}
		}

		conn, err := net.DialTimeout(netw, addr, config.ConnectTimeout)
		if err != nil {
			return nil, err
		}

		conn.SetDeadline(time.Now().Add(config.RWTimeout))

		return conn, nil
	}
}

func tlsConfig(config *Config) *tls.Config {
	if config == nil {
		return nil
	}

	if config.CertPEM == nil || config.KeyPEM == nil && config.CaCertPEM == nil {
		return nil
	}

	clientTLSCert, err := tls.X509KeyPair(config.CertPEM, config.KeyPEM)
	if err != nil {
		// Here, probably, should be an error
		return nil
	}

	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(config.CaCertPEM)

	return &tls.Config{
		RootCAs:      certPool,
		Certificates: []tls.Certificate{clientTLSCert},
	}
}

func client(config *Config) (*http.Client, error) {
	return &http.Client{
		Transport: &http.Transport{
			Dial:            dialer(config),
			TLSClientConfig: tlsConfig(config),
		},
	}, nil
}

func Post(url string, payload []byte, config *Config) ([]byte, error) {
	return DefaultResponseHandler(Request("POST", url, payload, config))
}

func Get(url string, config *Config) ([]byte, error) {
	return DefaultResponseHandler(Request("GET", url, nil, config))
}

func Request(method string, url string, payload []byte, config *Config) (*http.Response, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, &Error{0, err}
	}

	if config != nil && config.Headers != nil {
		req.Header = config.Headers
	} else {
		req.Header.Set("Content-Type", "application/json")
	}

	req.Close = true

	c, err := client(config)
	if err != nil {
		return nil, &Error{0, err}
	}

	resp, err := c.Do(req)
	if err != nil {
		return nil, &Error{0, err}
	}

	return resp, nil
}

func DefaultResponseHandler(resp *http.Response, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 226 {
		return nil, &Error{resp.StatusCode, errors.New("Non 2XX response")}
	}

	result, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return result, &Error{resp.StatusCode, err}
	}

	return result, nil
}

type Error struct {
	Code int
	msg  error
}

func (e *Error) Error() string {
	return e.msg.Error()
}
