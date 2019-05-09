package gaurun

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"runtime"

	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const version = "1.0"

var userAgent = fmt.Sprintf("GaurunGoClient/%s (%s)", version, runtime.Version())

// A Client for gaurun server.
type Client struct {
	Endpoint   *url.URL
	HTTPClient *http.Client
}

// NewClient generates a client for gaurun server.
func NewClient(endpoint string, cli *http.Client) (*Client, error) {
	u, err := url.ParseRequestURI(endpoint)
	if err != nil {
		return nil, errors.Wrapf(err, "gaurun: failed to parse url - endpoint = %s", endpoint)
	}
	if cli == nil {
		cli = http.DefaultClient
	}
	return &Client{
		Endpoint:   u,
		HTTPClient: cli,
	}, nil
}

// PushMulti sends payloads to gaurun server.
func (cli *Client) PushMulti(c context.Context, ps ...*Payload) error {
	eg, ctx := errgroup.WithContext(c)
	for _, p := range ps {
		eg.Go(func() error {
			return cli.Push(ctx, p)
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

// Push sends a payload to gaurun server.
func (cli *Client) Push(c context.Context, p *Payload) error {
	return cli.do(http.MethodPost, "/push", p)
}

func (cli *Client) do(method, spath string, body interface{}) error {
	bin, err := json.Marshal(body)
	if err != nil {
		return errors.Wrapf(err, "gaurun: failed to marshal json - body = %+v", body)
	}
	u := *cli.Endpoint
	u.Path = path.Join(cli.Endpoint.Path, spath)
	req, err := http.NewRequest(method, u.String(), bytes.NewReader(bin))
	if err != nil {
		return errors.Wrapf(err, "gaurun: failed to generate new http request - method = %s, path = %s, body = %+v", method, spath, body)
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-type", "application/json")
	resp, err := cli.HTTPClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "gaurun: failed to send http request - request = %+v", req)
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return NewError(resp)
}
