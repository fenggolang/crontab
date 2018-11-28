package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd    *exec.Cmd
		output []byte
		err    error
	)

	// 生成Cmd结构体对象
	cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe", "-c", "echo 1; echo 2;")

	//执行了命令，捕获了子进程的输出(pipe)
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}
