package storageutils

import (
	"encoding/json"
	"errors"
	"github.com/HolmesProcessing/Holmes-Totem-Service-Library/go/utilities/httputils"
	"io/ioutil"
	"net/http"
	"strings"
)

// Submit a given sample to Holmes-Storage.
// On success returns nil, otherwise the respective error.
func (s *Storage) SubmitSample(sample *StorageSample) error {
	var (
		r   *http.Response
		err error
		url string = s.Address + "/samples/"
	)

	ps := make(map[string]string)
	ps["user_id"] = s.UserID
	ps["source"] = sample.Source
	ps["name"] = sample.Name
	ps["date"] = sample.Date
	ps["comment"] = sample.Comment

	if sample.Tags != nil && len(sample.Tags) > 0 {
		parts := make([]string, len(sample.Tags))
		for i, tag := range sample.Tags {
			parts[i] = "tags[]=" + tag
		}
		url += "?" + strings.Join(parts, "&")
	}

	// Prefer FileContents over FilePath if set and has a non-zero size.
	if sample.FileContents != nil && len(sample.FileContents) > 0 {
		r, err = httputils.UploadFileBin(url, ps, "sample", sample.FilePath, sample.FileContents)
	} else {
		r, err = httputils.UploadFile(url, ps, "sample", sample.FilePath)
	}
	if err != nil {
		return err
	}

	// Read the entire body.
	bytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	// Definitly not a successful upload.
	if r.StatusCode == 500 {
		return errors.New(r.Status + " - " + string(bytes))
	}

	// Cannot tell whether successful or not by looking at the status code.
	// Unmarshal the body and see what the ResponseCode is.
	var rj struct {
		ResponseCode int
		Failure      string
	}
	err = json.Unmarshal(bytes, &rj)
	if err != nil {
		return errors.New(err.Error() + " -- " + string(bytes))
	}

	// Seems like it is a success after all.
	if rj.ResponseCode == 1 {
		return nil
	}

	// Otherwise it is no success, create a new error to reflect that.
	return errors.New(rj.Failure)
}
