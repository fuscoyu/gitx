package repo

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRepoPush_newBranchName(t *testing.T) {
	cfg = GetConfig()
	p := &RepoPush{
		config: cfg,
	}
	cfg.Patch.TmpBranchFmt = "{jiraID}_{jiraDesc}_{tgtBranch}"
	s := p.newBranchName("BILLING-888", "x", "staging")
	assert.Equal(t, "BILLING-888_x_staging", s)

	cfg.Patch.TmpBranchFmt = "{tgtBranch}_{jiraDesc}_{jiraID}"
	s = p.newBranchName("BILLING-888", "x", "staging")
	assert.Equal(t, "staging_x_BILLING-888", s)
}

func TestRepoPush_autoMergeBranchHook(t *testing.T) {
	cfg = GetConfig()
	repo := cfg.GetRepo("pitrix-billing")
	p := NewRepoPush(repo, cfg, "dev", nil, false)

	logrus.SetLevel(logrus.DebugLevel)
	p.AutoMergeBranchHook()
}
