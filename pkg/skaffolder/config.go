// Package skaffolder provides a programmatic API for managing SKA
// configurations and rendering blueprints. It complements the CLI by enabling
// library-style use in Go projects.
package skaffolder

import (
	"fmt"

	"github.com/apex/log"
	cfg "github.com/gchiesa/ska/internal/localconfigservice"
	"github.com/gchiesa/ska/pkg/util"
)

// SkaConfigTask exposes operations to list, rename, delete and query named
// SKA configurations stored under a project's .ska-config directory.
type SkaConfigTask struct {
	// BaseURI is the path to the project base directory.
	BaseURI string
	// Log is a contextual logger used by the task.
	Log *log.Entry
}

// NewSkaConfigTask constructs a SkaConfigTask bound to a base directory.
func NewSkaConfigTask(baseURI string) *SkaConfigTask {
	logCtx := log.WithFields(log.Fields{
		"pkg": "skaffolder",
	})
	return &SkaConfigTask{
		BaseURI: baseURI,
		Log:     logCtx,
	}
}

// ListNamedConfigs returns the names of all .ska-config configurations stored
// under the BaseURI.
func (c *SkaConfigTask) ListNamedConfigs() ([]string, error) {
	return cfg.ListNamedConfigs(c.BaseURI)
}

// RenameNamedConfig renames an existing named configuration from name to newName.
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

// DeleteConfig deletes the specified named configuration from .ska-config.
func (c *SkaConfigTask) DeleteConfig(name string) error {
	localConfig := cfg.NewLocalConfigService(name)

	if err := localConfig.ReadValidConfig(c.BaseURI); err != nil {
		return err
	}

	return localConfig.DeleteConfig(c.BaseURI)
}

// GetNamedConfigJSON returns the JSON representation of a named configuration.
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

// QueryNamedConfigJSON runs a JSONPath query against a named configuration's
// JSON representation and returns the resulting string.
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
