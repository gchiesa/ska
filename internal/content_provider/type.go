package content_provider

type RemoteContentProvider interface {
	// DownloadContent Download the remote content to a local working directory
	DownloadContent() error
	// RemoveWorkingDir Remove the local working directory
	RemoveWorkingDir() error
	// WorkingDir Get the local working directory
	WorkingDir() string
	// RemoteURI Get the remote URI
	RemoteURI() string
}

const (
	workingDirPrefix = "ska-content-provider-wd-"
)
