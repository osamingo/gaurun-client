package gaurun

import (
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mercari/gaurun/gaurun"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

func TestNewClient(t *testing.T) {

	expected := "gaurun: failed to parse url - endpoint = dummy: parse dummy: invalid URI for request"
	_, err := NewClient("dummy", nil)
	if expected != err.Error() {
		t.Errorf("failed to test - expected = %s, actual = %s", expected, err.Error())
	}

	url := "https://api.gaurun.io"
	cli, err := NewClient(url, nil)
	if err != nil {
		t.Error("an error occured", err)
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
		t.Error("an error occured", err)
	}

	p := &Payload{
		Notifications: []*Notification{{}},
	}

	err = cli.PushMulti(context.Background(), p)
	if err != nil {
		t.Error("an error occured", err)
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
		t.Error("an error occured", err)
	}

	err = cli.do(http.MethodPost, "/push", math.NaN())
	expected := "gaurun: failed to marshal json - body = NaN: json: unsupported value: NaN"
	if expected != err.Error() {
		t.Errorf("failed to test - expected = %s, actual = %s", expected, err.Error())
	}
}
