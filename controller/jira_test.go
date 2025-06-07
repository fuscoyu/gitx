package controller

import (
	"github.com/goeoeo/gitx/repo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestJiraController_AddJira(t *testing.T) {
	cfg := repo.GetConfig("../config.yaml")
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)

	err = jira.Add("dev-tool", "VM-2074", []string{"dev", "qa", "staging"})
	assert.Nil(t, err)
}

func TestJiraController_DelJira(t *testing.T) {
	cfg := repo.GetConfig("../config.yaml")
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)

	err = jira.Del("dev-tool", "VM-2074")
	assert.Nil(t, err)
}

func TestJiraController_Clear(t *testing.T) {
	cfg := repo.GetConfig("")
	jira, err := NewJiraController(cfg)
	assert.Nil(t, err)

	err = jira.Clear()
	assert.Nil(t, err)
}
