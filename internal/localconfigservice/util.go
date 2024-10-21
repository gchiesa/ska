package localconfigservice

func ListNamedConfigs(dirPath string) ([]string, error) {
	return configEntries(makeConfigPath(dirPath))
}
