package main

import (
	"fmt"
	"os/exec"
)

func main() {
	var (
		cmd *exec.Cmd
		err error
	)

	// windows下执行
	//cmd = exec.Command("C:\\cygwin64\\bin\\bash.exe", "-c", "echo 1;echo 2;")
	// linux下执行
	//cmd = exec.Command("/bin/bash", "-c", "sleep 5s;echo 1; echo 2;")
	//cmd = exec.Command("/bin/bash", "-c", "ping -c 4 localhost")
	cmd = exec.Command("/bin/bash", "-c", "ansible all -m ping")

	err = cmd.Run() // 同步执行
	//err = cmd.Start() // 异步执行

	fmt.Println(err)
}
