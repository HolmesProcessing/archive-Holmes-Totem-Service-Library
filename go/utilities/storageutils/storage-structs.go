package storageutils

// A simple interface to interact with a Holmes-Storage server.
// The UserID is the id to be used for submissions to the storage server.
type Storage struct {
	Address string
	UserID  string
}

// Basic struct allowing for specification of either an in-memory or on-disk
// file. (If FileContents is set, FilePath is only used for the multipart forms
// file parameter name)
type StorageSample struct {
	// Local only:
	FilePath     string
	FileContents []byte

	// Remote data:
	Source  string
	Name    string
	Date    string
	Tags    []string
	Comment string
}
