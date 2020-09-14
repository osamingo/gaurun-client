package gaurun

import (
	"context"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mercari/gaurun/gaurun"
	"go.uber.org/zap"
)

func TestNewClient(t *testing.T) {

	_, err := NewClient("dummy", nil)
	if err == nil {
		t.Error("failed to test - expected = error, actual = nil")
	}

	// the output of errors from net.Url was changed in go1.14 to quote the
	// URL. Therefore, we'll just strip the quote until n-2 go versions have
	// moved to error quoting.
	errStr := strings.Replace(err.Error(), "\"", "", -1)
	expected := "gaurun: failed to parse url - endpoint = dummy: parse dummy: invalid URI for request"
	if errStr != expected {
		t.Errorf("failed to test - expected = %s, actual = %s", expected, errStr)
	}

	url := "https://api.gaurun.io"
	cli, err := NewClient(url, nil)
	if err != nil {
		t.Error("an error occurred", err)
	}
	if url != cli.Endpoint.String() {
		t.Errorf("failed to test - expected = %s, actual = %s", url, cli.Endpoint.String())
	}
	if cli.HTTPClient == nil {
		t.Errorf("http client should not be nil")
	}
}

func TestClient_PushMulti(t *testing.T) {

	gaurun.LogAccess = zap.NewNop()
	gaurun.LogError = zap.NewNop()
	gaurun.ConfGaurun.Core.NotificationMax = 10

	mux := http.NewServeMux()
	gaurun.RegisterHandlers(mux)
	srv := httptest.NewServer(mux)

	cli, err := NewClient(srv.URL, nil)
	if err != nil {
		t.Error("an error occurred", err)
	}

	p := &Payload{
		Notifications: []*Notification{{}},
	}

	err = cli.PushMulti(context.Background(), p)
	if err != nil {
		t.Error("an error occurred", err)
	}

	p.Notifications = p.Notifications[:0]

	err = cli.PushMulti(context.Background(), p)
	expected := "gaurun: an error occurred - status_code = 400, response = {Message:empty notification}"
	if expected != err.Error() {
		t.Errorf("failed to test - expected = %s, actual = %s", expected, err.Error())
	}

	srv.Close()

	err = cli.PushMulti(context.Background(), p)
	wantPrefix := "gaurun: failed to send http request"
	if !strings.HasPrefix(err.Error(), wantPrefix) {
		t.Errorf("failed to test - want_prefix = %s, actual = %s", wantPrefix, err.Error())
	}
}

func TestClient_do(t *testing.T) {

	cli, err := NewClient("https://api.gaurun.io", nil)
	if err != nil {
		t.Error("an error occurred", err)
	}

	err = cli.do(http.MethodPost, "/push", math.NaN())
	expected := "gaurun: failed to marshal json - body = NaN: json: unsupported value: NaN"
	if expected != err.Error() {
		t.Errorf("failed to test - expected = %s, actual = %s", expected, err.Error())
	}
}
