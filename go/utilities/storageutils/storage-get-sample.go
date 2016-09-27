package storageutils

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

// Submit a given sample to Holmes-Storage.
// On success returns nil, otherwise the respective error.
func (s *Storage) GetSample(sha256 string) ([]byte, error) {
	r, err := http.Get(s.Address + "/samples/" + sha256)
	if err != nil {
		return nil, err
	}

	// Read the entire body. Might contain useful error information.
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Not a successful download.
	if r.StatusCode != 200 {
		return nil, errors.New(r.Status + " - " + string(bytes))
	}

	// Check if content-type is set and if it is application/octet-stream
	if t := r.Header.Get("Content-Type"); t != "application/octet-stream" {
		// In this case we probably got a JSON error message
		var x struct {
			Failure string
		}
		err := json.Unmarshal(bytes, &x)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(x.Failure)
	}

	// Success
	return bytes, nil
}
