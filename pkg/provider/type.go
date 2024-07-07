package provider

type RemoteContentProviderService interface {
	SetupWorkingDir() (string, error)
	DownloadContent() error
	RemoveWorkingDir() error
	RemotePath() string
	Path() string
}
