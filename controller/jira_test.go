package controller

import (
	"github.com/goeoeo/gitx/repo"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// getTestConfigPath returns the path to the test configuration file
func getTestConfigPath() string {
	// Try to find test_config.yaml in current directory first
	if _, err := os.Stat("test_config.yaml"); err == nil {
		return "test_config.yaml"
	}
	
	// Try to find testdata/config.yaml
	wd, _ := os.Getwd()
	for {
		testConfigPath := filepath.Join(wd, "testdata", "config.yaml")
		if _, err := os.Stat(testConfigPath); err == nil {
			return testConfigPath
		}
		
		// Try test_config.yaml in the project root
		testConfigPath = filepath.Join(wd, "test_config.yaml")
		if _, err := os.Stat(testConfigPath); err == nil {
			return testConfigPath
		}
		
		// Move up one directory
		parent := filepath.Dir(wd)
		if parent == wd {
			break
		}
		wd = parent
	}
	
	// Fallback to config.yaml in current directory
	return "config.yaml"
}

func TestJiraController_AddJira(t *testing.T) {
	// Skip test if config file is not available
	configPath := getTestConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("Skipping test: config file not found: %s", configPath)
	}
	
	cfg := repo.GetConfig(configPath)
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)

	err = jira.Add("dev-tool", "VM-2074", []string{"dev", "qa", "staging"})
	assert.Nil(t, err)
}

func TestJiraController_DelJira(t *testing.T) {
	// Skip test if config file is not available
	configPath := getTestConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("Skipping test: config file not found: %s", configPath)
	}
	
	cfg := repo.GetConfig(configPath)
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)

	err = jira.Del("dev-tool", "VM-2074")
	assert.Nil(t, err)
}

func TestJiraController_Clear(t *testing.T) {
	jira := getJiraController(t)
	err := jira.Clear("production", true)
	assert.Nil(t, err)
}

func TestJiraController_PrintJira(t *testing.T) {
	jira := getJiraController(t)
	err := jira.Print("production", "BILLING-3037")
	assert.Nil(t, err)
}

func TestJiraController_Detach(t *testing.T) {
	jira := getJiraController(t)

	err := jira.Detach("production", "BILLING-3383", "v6.2")
	assert.Nil(t, err)
}

func TestJiraController_CheckBranchMerged(t *testing.T) {
	jira := getJiraController(t)
	err := jira.CheckBranchMerged("production", "BILLING-3037")
	assert.Nil(t, err)
}

func getJiraController(t *testing.T) *JiraController {
	// Skip test if config file is not available
	configPath := getTestConfigPath()
	if _, err := os.Stat(configPath); err != nil {
		t.Skipf("Skipping test: config file not found: %s", configPath)
	}
	
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{})
	cfg := repo.GetConfig(configPath)
	cfg.DisableInitLog = true
	cfg.Init()
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)
	return jira
}
