package httputils

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

// Used by the Request type.
// If FileContent is not nil, then FilePath is ignored.
// Otherwise content is read from the location specified in FilePath.
// Name must not be empty and is the parameter name used in the multipart form.
type RequestFile struct {
	ParameterName string
	FilePath      string
	FileContent   []byte
}

// Return the file's contents. If the FileContent field is not set, try to read
// from the path specified in the field FilePath.
func (rf *RequestFile) getContents() ([]byte, error) {
	if rf.FileContent == nil {
		bytes, err := ioutil.ReadFile(rf.FilePath)
		if err != nil {
			return nil, err
		}
		rf.FileContent = bytes
	}
	return rf.FileContent, nil
}

// Convenience wrapper for outgoing requests functionality.
type Request struct {
	// The request destination address.
	URL string

	// Should be one of GET, POST, PUT, or DELETE. This is directly passed to
	// the respective functions in net/http, no validity checks are performed.
	Method string

	// The request parameters. If the request is a POST request, then the
	// parameters (together with a file if specified) are put into a multipart
	// form and sent in the request body.
	Parameters map[string]string

	// The requests cookies.
	Cookies []*http.Cookie

	// Potential files for upload using a multipart form in the request body.
	// If nil, ignored. Can also be used to upload custom bodies (e.g.
	// marshalled JSON).
	Files []*RequestFile

	// Alternatively, if the request is not of type POST and no Files are
	// specified, a Body may be specified. It is ignored if either of the
	// stated prerequisites is false.
	Body []byte

	// The content type. If empty it will be omitted. Note that if the Files
	// field is set, the content typte is forced to multipart/form-data.
	ContentType string
}

// Create the request object defined by the given url, method, parameters and
// possibly file.
// Sets the content type if provided and eligible (i.e. not multipart form) and
// also sets the cookies if provided.
func (r *Request) Build() (*http.Request, error) {
	body, contentType, err := r.getBody()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(r.Method, r.getURL(), body)
	if err != nil {
		return nil, err
	}
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}
	if r.Cookies != nil {
		for _, cookie := range r.Cookies {
			request.AddCookie(cookie)
		}
	}
	return request, nil
}

// Helper function that url-encodes the requests' parameters.
func (r *Request) urlEncodedParameters() string {
	if r.Parameters != nil && len(r.Parameters) > 0 {
		data := &url.Values{}
		for key, val := range r.Parameters {
			data.Add(key, val)
		}
		return data.Encode()
	}
	return ""
}

// Helper function to get the request URL (depending on used method).
func (r *Request) getURL() string {
	if r.Method != "POST" {
		params := r.urlEncodedParameters()
		if params != "" {
			sep := "?"
			if strings.Contains(r.URL, "?") {
				sep = "&"
			}
			return r.URL + sep + params
		}
	}
	return r.URL
}

// Helper function to create the request body.
func (r *Request) getBody() (*bytes.Buffer, string, error) {
	var (
		err         error         = nil
		body        *bytes.Buffer = nil
		contentType string        = ""
	)

	// Special treatment for multipart forms.
	// Otherwise distinguish between POST method and other methods.
	// POST allows for either parameters or a custom body.
	if r.Files != nil && len(r.Files) > 0 {
		body, contentType, err = r.getMultipartBody()

	} else if r.Method == "POST" {
		if r.Files == nil || len(r.Files) == 0 {
			params := r.urlEncodedParameters()
			if params != "" {
				body = bytes.NewBuffer([]byte(params))
			}
		}
	}

	if body == nil {
		if r.Body != nil {
			body = bytes.NewBuffer(r.Body)
		} else {
			body = bytes.NewBuffer([]byte{})
		}
	}

	return body, contentType, err
}

// Helper function for creating multipart requests.
func (r *Request) getMultipartBody() (*bytes.Buffer, string, error) {
	// Create the buffer and writer for the multipart body.
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Write parameters first (if any).
	if r.Method == "POST" && r.Parameters != nil {
		for key, val := range r.Parameters {
			err := writer.WriteField(key, val)
			if err != nil {
				return nil, "", err
			}
		}
	}

	// Write files.
	for _, file := range r.Files {
		part, err := writer.CreateFormFile(file.ParameterName, filepath.Base(file.FilePath))
		if err != nil {
			return nil, "", err
		}
		bytes, err := file.getContents()
		if err != nil {
			return nil, "", err
		}
		part.Write(bytes)
	}

	// Close the writer to flush all data to the buffer.
	err := writer.Close()
	if err != nil {
		return nil, "", err
	}

	// Return the created body alongside the specific multipart content type.
	return body, writer.FormDataContentType(), nil
}

// Convenience method. Builds the request and executes it, returning the servers
// response or an error if one occured.
func (r *Request) Run() (*http.Response, error) {
	r.Method = strings.ToUpper(r.Method)
	request, err := r.Build()
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	return client.Do(request)
}
