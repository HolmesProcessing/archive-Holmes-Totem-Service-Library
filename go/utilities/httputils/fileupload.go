package httputils

import (
	"net/http"
)

// Convenience function for uploading a file to the provided URL using a GET
// request. Attempts to read the file from the given path.
func UploadFile(url string, extraParameters map[string]string, parameterName, filePath string) (*http.Response, error) {
	return UploadFileBin(url, extraParameters, parameterName, filePath, nil)
}

// Convenience function for uploading a file to the provided URL using a GET
// request. If fileData is not provided (nil), attempts to read from the given
// path.
func UploadFileBin(url string, extraParameters map[string]string, parameterName, filePath string, fileContent []byte) (*http.Response, error) {
	request := createFileUploadRequest(url, extraParameters, parameterName, filePath, fileContent)
	return request.Run()
}

// Helper function. Exists to make testing possible.
func createFileUploadRequest(url string, extraParameters map[string]string, parameterName, filePath string, fileContent []byte) *Request {
	file := &RequestFile{
		ParameterName: parameterName,
		FilePath:      filePath,
		FileContent:   fileContent,
	}
	request := &Request{
		URL:        url,
		Method:     "PUT",
		Parameters: extraParameters,
		Files: []*RequestFile{
			file,
		},
	}
	return request
}
