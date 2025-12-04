package repo

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
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
	sdo := stdout.String()
	sde := stderr.String()
	if isPrint {
		logrus.Debugf("out:\n[%s]\n", sdo)
		if sde != "" {
			if err != nil {
				// 命令失败时，stderr 作为错误输出
				logrus.Errorf("err:%s", sde)
			} else {
				// 命令成功时，stderr 作为调试信息（git 的信息性消息常输出到 stderr）
				logrus.Debugf("stderr:%s", sde)
			}
		}
	}
	if err != nil {
		logrus.Debugf("failed to call Run(): %v \n", err)
	}
	cmdRet := &CmdRet{
		Out:    sdo,
		ErrStr: sde,
	}
	return cmdRet, err
}
