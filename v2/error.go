package gaurun

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/mercari/gaurun/gaurun"
	"github.com/pkg/errors"
)

// An Error is error response from gaurun server.
type Error struct {
	StatusCode int
	Response   gaurun.ResponseGaurun
}

// Error implements error interface.
func (err *Error) Error() string {
	return fmt.Sprintf("gaurun: an error occurred - status_code = %d, response = %+v", err.StatusCode, err.Response)
}

// NewError generates *gaurun.Error from *http.Response.
func NewError(resp *http.Response) error {
	err := &Error{}
	if err := json.NewDecoder(resp.Body).Decode(&err.Response); err != nil {
		return errors.Wrapf(err, "gaurun: failed to decode http response")
	}
	err.StatusCode = resp.StatusCode
	return err
}
