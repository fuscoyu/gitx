package repo

import (
	"bytes"
	"context"
	"github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type CmdRet struct {
	Out    string
	ErrStr string
}

func ExecCmdCtx(ctx context.Context, dir string, name string, arg ...string) (*CmdRet, error) {
	logrus.Debugf("cmd: dir: [%s], name [%s], arg: %v \n", dir, name, arg)
	cmd := exec.CommandContext(ctx, name, arg...)
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	return execC(ctx, cmd)
}

func ExecCmd(dir string, name string, arg ...string) (*CmdRet, error) {
	logrus.Debugf("cmd: dir: [%s], cmd:%s %s \n", dir, name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	if len(dir) > 0 {
		cmd.Dir = dir
	}
	return execC(context.Background(), cmd)
}

func execC(ctx context.Context, cmd *exec.Cmd) (*CmdRet, error) {
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	isPrint, ok := ctx.Value("print").(bool)
	if !ok {
		isPrint = true
	}
	err := cmd.Run()
	if err != nil {
		logrus.Debugf("failed to call Run(): %v \n", err)
	}
	sdo := stdout.String()
	sde := stderr.String()
	if isPrint {
		logrus.Debugf("out:\n[%s]\n", sdo)
		if sde != "" {
			logrus.Errorf("err:%s", sde)
		}
	}
	cmdRet := &CmdRet{
		Out:    sdo,
		ErrStr: sde,
	}
	return cmdRet, err
}
