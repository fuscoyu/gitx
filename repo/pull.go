package repo

import (
	"errors"
	"github.com/sirupsen/logrus"
	"log"
)

func (rp *RepoPatch) Pull(isDel bool) error {
	for _, tgtBranch := range rp.Patch.GetTgtBranchs() {
		pRepo := NewRepoPull(rp.Repo, rp.Patch, tgtBranch)
		err := pRepo.GitRepo.LsRemote()
		if err != nil {
			logrus.Debugf("git remote connection exception, please check; repo: %s \n", rp.Repo.Path)
			return err
		}

		err = pRepo.pull(isDel)
		if err != nil {
			logrus.Debugf("git repo pull faild: repo: %s, target branch [%s], err: %v \n",
				rp.Repo.Path, tgtBranch, err)
			return err
		}
	}
	return nil
}

type RepoPullPatch struct {
	DevBranch string
	TgtBranch string
}

type RepoPull struct {
	GitRepo       *GitRepo
	RepoPullPatch *RepoPullPatch
}

func NewRepoPull(r *Repo, p *Patch, tgtBranch string) *RepoPull {
	gRepo := NewGitRepo(r.Path, r.Url)
	repoPullPatch := &RepoPullPatch{
		DevBranch: p.DevBranch,
		TgtBranch: tgtBranch,
	}
	return &RepoPull{
		GitRepo:       gRepo,
		RepoPullPatch: repoPullPatch,
	}
}

// pull 本地目标分支, 如 本地 staging 分支。
// 如果有 delete 参数, 先 delete 本地目标分支,
// 如 staging,
// 1. git branch -D staging
// 2. git fetch
// 3. git checkout -b staging origin/staging
func (r *RepoPull) pull(isDel bool) error {

	logrus.Debugf("begin repo pull [%s] branch [%s] ... \n", r.GitRepo.Path, r.RepoPullPatch.TgtBranch)

	master := "master"
	tgtBranch := r.RepoPullPatch.TgtBranch
	devBranch := r.RepoPullPatch.DevBranch

	if isDel {
		err := r.GitRepo.SwitchBranch(master)
		if err != nil {
			logrus.Debugf("switch git branch faild: repo: %s, branch [%s], err: %v \n",
				r.GitRepo.Path, master, err)
			return err
		}

		if tgtBranch == devBranch || tgtBranch == master {
			logrus.Debugf("target branch error: repo: %s, branch [%s]",
				r.GitRepo.Path, tgtBranch)
			return errors.New("target branch error")
		}
		err = r.GitRepo.DelLocalBranch(tgtBranch)
		if err != nil {
			logrus.Debugf("git del local branch faild: repo: %s, branch [%s], err: %v \n",
				r.GitRepo.Path, tgtBranch, err)
			return err
		}

		err = r.GitRepo.Fetch()
		if err != nil {
			logrus.Debugf("git fetch faild: repo: %s, err: %v \n",
				r.GitRepo.Path, err)
			return err
		}

		err = r.GitRepo.NewBranchFromRemote(tgtBranch)
		if err != nil {
			logrus.Debugf("create new branch from remote faild: repo: %s, err: %v \n", r.GitRepo.Path, err)
			return err
		}
	}

	err := r.GitRepo.SwitchBranch(tgtBranch)
	if err != nil {
		logrus.Debugf("switch git branch faild: repo: %s, branch [%s], err: %v \n",
			r.GitRepo.Path, tgtBranch, err)
		return err
	}
	err = r.GitRepo.Pull()
	if err != nil {
		logrus.Debugf("git pull faild: repo: %s, branch [%s], err: %v \n",
			r.GitRepo.Path, tgtBranch, err)
		return err
	}
	log.Println("git pull ok")
	return nil
}
