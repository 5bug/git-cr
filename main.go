package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

func gitBranch() (branch string, err error) {
	cmd := exec.Command("/bin/bash", "-c", "git rev-parse --abbrev-ref HEAD")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "请在git仓库目录下运行！！！")
		return
	}
	branch = strings.TrimSpace(string(out))
	return
}

func gitUser() (username string, err error) {
	cmd := exec.Command("/bin/bash", "-c", "git config user.name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrap(err, "请在git仓库目录下运行！！！")
		return
	}
	username = strings.TrimSpace(string(out))
	return
}

func gitPull(branch string) error {
	fmt.Printf("开始同步分支 %s 远端代码到本地...\n", branch)
	shell := `
#!/bin/bash

set -e

export TARGET_BRANCH="%s"
RES=$(git ls-remote --heads origin refs/heads/${TARGET_BRANCH})
if [[ "RES" == "" ]]; then
  echo "远端仓库里不存在分支:${TARGET_BRANCH},请先创建好该分支"
  exit 1
fi

git pull origin ${TARGET_BRANCH}
	`
	args := []string{"-c", fmt.Sprintf(shell, branch)}
	cmd := exec.Command("/bin/sh", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrapf(err, "执行失败，错误信息如下：\n%s", string(out))
		return err
	}
	fmt.Println("同步远端代码到本地成功.")
	return nil
}

func gitPush(username, branch string) error {
	fmt.Println("推送代码到远端仓库并发起评审任务...")
	shell := `
#!/bin/bash

set -e

export USER_NAME="%s"
export TARGET_BRANCH="%s"
export SOURCE_BRANCH="cr/${USER_NAME}/${TARGET_BRANCH}"

git push origin HEAD:${SOURCE_BRANCH} \
-o merge_request.create \
-o merge_request.title="%s" \
-o merge_request.target=${TARGET_BRANCH} \
-o merge_request.source=${SOURCE_BRANCH} \
-o merge_request.remove_source_branch=true
	`
	args := []string{"-c", fmt.Sprintf(shell, username, branch, fmt.Sprintf("%s merge code into %s", username, branch))}

	fmt.Println(fmt.Sprintf(shell, username, branch, fmt.Sprintf("%s merge code into %s", username, branch)))
	cmd := exec.Command("/bin/sh", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		err = errors.Wrapf(err, "执行失败，错误信息如下：\n%s", string(out))
		return err
	}
	fmt.Println("推送代码并创建评审任务成功，详情如下：")
	fmt.Println(string(out))
	return nil
}

func main() {
	branch, err := gitBranch()
	if err != nil {
		fmt.Println(err)
		return
	}
	username, err := gitUser()
	if err != nil {
		fmt.Println(err)
		return
	}
	if err = gitPull(branch); err != nil {
		fmt.Println(err)
		return
	}
	if err = gitPush(username, branch); err != nil {
		fmt.Println(err)
		return
	}
}
