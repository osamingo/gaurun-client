package gaurun

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestNewError(t *testing.T) {
	tbl := []struct {
		Response *http.Response
		Expected string
	}{
		{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString("{")),
			},
			Expected: "gaurun: failed to decode http response: unexpected EOF",
		},
		{
			Response: &http.Response{
				StatusCode: http.StatusBadRequest,
				Body:       ioutil.NopCloser(bytes.NewBufferString(`{"message":"empty notification"}`)),
			},
			Expected: "gaurun: an error occurred - status_code = 400, response = {Message:empty notification}",
		},
	}
	for i := range tbl {
		err := NewError(tbl[i].Response)
		if tbl[i].Expected != err.Error() {
			t.Errorf("failed to test - expected = %s, actual = %s", tbl[i].Expected, err.Error())
		}
	}
}
