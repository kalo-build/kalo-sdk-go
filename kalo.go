// Package kalo provides the Kalo SDK for plugin development.
// Plugin authors use this SDK to access stores configured in kalo.yaml
// without needing to know whether they're backed by WASI filesystem
// or host function callbacks.
package kalo

import (
	"encoding/json"
	"os"
	"sync"
)

// StoreConfig represents the configuration for a single store
// passed from the CLI to the plugin via command line arguments.
type StoreConfig struct {
	ID        uint32 `json:"id"`
	Type      string `json:"type"`      // localFileSystem, cloudSqlDatabase
	MountPath string `json:"mountPath"` // For localFileSystem stores
}

// PluginConfig represents the configuration passed from CLI to plugin.
type PluginConfig struct {
	Stores map[string]StoreConfig `json:"stores"`
	Config map[string]interface{} `json:"config,omitempty"`
}

var (
	globalConfig     *PluginConfig
	globalConfigOnce sync.Once
	configErr        error
)

// Init initializes the SDK by parsing the configuration from command line args.
// This is called automatically on first use of FS() or DB().
func Init() error {
	globalConfigOnce.Do(func() {
		if len(os.Args) < 2 {
			// No config provided, use empty config
			globalConfig = &PluginConfig{
				Stores: make(map[string]StoreConfig),
			}
			return
		}

		configJSON := os.Args[1]
		globalConfig = &PluginConfig{}
		if err := json.Unmarshal([]byte(configJSON), globalConfig); err != nil {
			configErr = err
			return
		}

		if globalConfig.Stores == nil {
			globalConfig.Stores = make(map[string]StoreConfig)
		}
	})
	return configErr
}

// getStoreConfig returns the configuration for a named store.
func getStoreConfig(name string) (StoreConfig, bool) {
	if err := Init(); err != nil {
		return StoreConfig{}, false
	}
	cfg, ok := globalConfig.Stores[name]
	return cfg, ok
}

// FS returns a FileStore for the named store.
// The store must be configured as type "localFileSystem" in kalo.yaml.
func FS(storeName string) FileStore {
	cfg, ok := getStoreConfig(storeName)
	if !ok {
		// Return a store that will error on use
		return &wasiFileStore{storeName: storeName, basePath: "", err: ErrStoreNotFound}
	}

	if cfg.Type != "localFileSystem" {
		return &wasiFileStore{storeName: storeName, basePath: "", err: ErrInvalidStoreType}
	}

	return &wasiFileStore{
		storeName: storeName,
		basePath:  cfg.MountPath,
	}
}

// DB returns a DBStore for the named store.
// The store must be configured as type "cloudSqlDatabase" in kalo.yaml.
func DB(storeName string) DBStore {
	cfg, ok := getStoreConfig(storeName)
	if !ok {
		return &hostDBStore{storeName: storeName, storeID: 0, err: ErrStoreNotFound}
	}

	if cfg.Type != "cloudSqlDatabase" {
		return &hostDBStore{storeName: storeName, storeID: 0, err: ErrInvalidStoreType}
	}

	return &hostDBStore{
		storeName: storeName,
		storeID:   cfg.ID,
	}
}

// GetConfig returns the plugin configuration passed from the CLI.
func GetConfig() map[string]interface{} {
	if err := Init(); err != nil {
		return nil
	}
	return globalConfig.Config
}

