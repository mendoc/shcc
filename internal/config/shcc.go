package config

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/mendoc/shcc/internal/system"
)

type ShccLocalConfig struct {
	Name string `json:"name"`
}

// GetShccConfigDir renvoie ~/.shcc/
func GetShccConfigDir() string {
	home := system.GetHomeDir()
	return filepath.Join(home, ".shcc")
}

// ReadShccConfig lit la config locale shcc
func ReadShccConfig() (*ShccLocalConfig, error) {
	path := filepath.Join(GetShccConfigDir(), "config.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &ShccLocalConfig{}, nil
		}
		return nil, err
	}
	var cfg ShccLocalConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveShccConfig enregistre la config locale shcc
func SaveShccConfig(cfg *ShccLocalConfig) error {
	dir := GetShccConfigDir()
	os.MkdirAll(dir, 0700)
	path := filepath.Join(dir, "config.json")
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0600)
}
