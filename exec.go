package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/gorilla/websocket"
	"k8s.io/client-go/rest"
)

// ExecOptions describe a execute request args.
type ExecOptions struct {
	Namespace string
	Pod       string
	Container string
	Command   []string
	TTY       bool
	Stdin     bool
}

type RoundTripCallback func(c *websocket.Conn) error

type WebsocketRoundTripper struct {
	TLSConfig *tls.Config
	Callback  RoundTripCallback
}

var cache string

var protocols = []string{
	"v4.channel.k8s.io",
	"v3.channel.k8s.io",
	"v2.channel.k8s.io",
	"channel.k8s.io",
}

const (
	stdin = iota
	stdout
	stderr
)

func WebsocketCallback(c *websocket.Conn) error {
	errChan := make(chan error, 3)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			buf := make([]byte, 1025)
			n, err := os.Stdin.Read(buf[1:])
			if err != nil {
				errChan <- err
				return
			}
			cache = strings.TrimSpace(string(buf[1 : n+1]))
			if err := c.WriteMessage(websocket.BinaryMessage, buf[:n+1]); err != nil {
				errChan <- err
				return
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			_, buf, err := c.ReadMessage()
			if err != nil {
				errChan <- err
				return
			}
			var w io.Writer
			switch buf[0] {
			case stdout:
				w = os.Stdout
			case stderr:
				w = os.Stderr
			}

			if strings.TrimSpace(string(buf[1:])) != cache {
				_, err = w.Write(buf)
				if err != nil {
					errChan <- err
					return
				}
			}
		}
	}()

	wg.Wait()
	close(errChan)
	err := <-errChan
	return err
}

func (wrt *WebsocketRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	dialer := &websocket.Dialer{
		Proxy:           http.ProxyFromEnvironment,
		TLSClientConfig: wrt.TLSConfig,
		Subprotocols:    protocols,
	}

	conn, resp, err := dialer.Dial(r.URL.String(), r.Header)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	return resp, wrt.Callback(conn)
}

func ExecRoundTripper(config *rest.Config, f RoundTripCallback) (http.RoundTripper, error) {
	tlsConfig, err := rest.TLSConfigFor(config)
	if err != nil {
		return nil, err
	}

	rt := &WebsocketRoundTripper{
		Callback:  f,
		TLSConfig: tlsConfig,
	}
	return rest.HTTPWrappersForConfig(config, rt)
}

func ExecRequest(config *rest.Config, opts *ExecOptions) (*http.Request, error) {
	u, err := url.Parse(config.Host)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "https":
		u.Scheme = "wss"
	case "http":
		u.Scheme = "ws"
	default:
		return nil, fmt.Errorf("Unrecognised URL scheme in %v", u)
	}

	u.Path = fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/exec", opts.Namespace, opts.Pod)

	rawQuery := "stdout=true&tty=true"
	for _, c := range opts.Command {
		rawQuery += "&command=" + c
	}

	if opts.Container != "" {
		rawQuery += "&container=" + opts.Container
	}

	if opts.TTY {
		rawQuery += "&tty=true"
	}

	if opts.Stdin {
		rawQuery += "&stdin=true"
	}
	u.RawQuery = rawQuery

	return &http.Request{
		Method: http.MethodGet,
		URL:    u,
	}, nil
}
