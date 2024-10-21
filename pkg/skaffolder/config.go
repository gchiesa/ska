package skaffolder

import (
	"fmt"
	"github.com/apex/log"
	cfg "github.com/gchiesa/ska/internal/localconfigservice"
	"github.com/gchiesa/ska/pkg/util"
)

type SkaConfigTask struct {
	BaseURI string
	Log     *log.Entry
}

func NewSkaConfigTask(baseURI string) *SkaConfigTask {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaConfigTask{
		BaseURI: baseURI,
		Log:     logCtx,
	}
}

func (c *SkaConfigTask) ListNamedConfigs() ([]string, error) {
	return cfg.ListNamedConfigs(c.BaseURI)
}

func (c *SkaConfigTask) RenameNamedConfig(name, newName string) error {
	localConfig := cfg.NewLocalConfigService(name)

	if err := localConfig.ReadValidConfig(c.BaseURI); err != nil {
		return err
	}

	err := localConfig.RenameNamedConfig(c.BaseURI, newName)
	if err == nil {
		c.Log.WithFields(log.Fields{"name": name, "newName": newName}).Infof("Renamed config from %s to %s", name, newName)
	}
	return err
}

func (c *SkaConfigTask) DeleteConfig(name string) error {
	localConfig := cfg.NewLocalConfigService(name)

	if err := localConfig.ReadValidConfig(c.BaseURI); err != nil {
		return err
	}

	return localConfig.DeleteConfig(c.BaseURI)
}

func (c *SkaConfigTask) GetNamedConfigJSON(namedConfig string) (string, error) {
	// configservice
	localConfig := cfg.NewLocalConfigService(namedConfig)

	// check if localconfig already exist, if yes we fail
	if !localConfig.ConfigExists(c.BaseURI) {
		return "", fmt.Errorf("unable to find named configuration: %s at the path: %s", namedConfig, c.BaseURI)
	}

	if err := localConfig.ReadValidConfig(c.BaseURI); err != nil {
		return "", err
	}

	jsonData, err := localConfig.ToJSON()
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func (c *SkaConfigTask) QueryNamedConfigJSON(namedConfig, jsonpath string) (string, error) {
	jsonData, err := c.GetNamedConfigJSON(namedConfig)
	if err != nil {
		return "", err
	}
	result, err := util.QueryJSONString(jsonData, jsonpath)
	if err != nil {
		return "", err
	}
	return result, nil
}
