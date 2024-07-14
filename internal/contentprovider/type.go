package contentprovider

type RemoteContentProvider interface {
	// DownloadContent Download the remote content to a local working directory
	DownloadContent() error
	// Cleanup Perform the cleanup task for the provider
	Cleanup() error
	// WorkingDir Get the local working directory
	WorkingDir() string
	// RemoteURI Get the remote URI
	RemoteURI() string
}

const (
	workingDirPrefix = "ska-content-provider-wd-"
)
