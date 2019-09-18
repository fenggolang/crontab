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
	//cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe", "-c", "echo 1; echo 2;")
	cmd = exec.Command("/bin/bash", "-c", "sleep 10s;echo 1; echo 2;")

	//执行了命令，捕获了子进程的输出(pipe)
	//cmd.Start() // 异步执行，不等到执行完即返回
	//cmd.Run() // 同步执行，底部调用的是Start(),等待Start的任务执行完
	//cmd.CombinedOutput() // 同步执行？其底层调用的是Run()
	if output, err = cmd.CombinedOutput(); err != nil {
		fmt.Println(err)
		return
	}

	// 打印子进程的输出
	fmt.Println(string(output))
}
