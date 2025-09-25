package controller

import (
	"github.com/goeoeo/gitx/repo"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestConfigLoading(t *testing.T) {
	configPath := getTestConfigPath()
	
	// Check if config file exists
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("Skipping test: config file not found: %s", configPath)
	}
	
	// Test config loading
	cfg := repo.GetConfig(configPath)
	assert.NotNil(t, cfg)
	
	// Test config initialization
	cfg.DisableInitLog = true
	cfg.Init()
	
	// Verify config has expected values
	assert.NotNil(t, cfg.Patch)
	assert.NotNil(t, cfg.GitLabConfigs)
	assert.NotNil(t, cfg.Repo)
}
