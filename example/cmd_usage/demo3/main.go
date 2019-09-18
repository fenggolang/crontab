package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

type result struct {
	err    error
	output []byte
}

func main() {
	// 执行1个cmd,让它在一个协程里去执行，让它执行2秒: sleep 2; echo hello;
	// 1秒的时候，我们杀死cmd
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		cmd        *exec.Cmd
		resultChan chan *result // 把命令执行结果写入resultChan 这个channel中
		res        *result      // 读取出resultChan中的数据到res中
	)

	// 创建了一个结果队列
	resultChan = make(chan *result, 1000)

	// context: chan byte
	// cancelFunc: close(chan byte)
	// ctx用于感知chan被关闭，cancelFunc用于关闭chan
	ctx, cancelFunc = context.WithCancel(context.TODO())

	// 开启一个协程执行一个bash任务
	go func() {
		var (
			output []byte
			err    error
		)

		//cmd = exec.CommandContext(ctx, "C:\\cygwin64\\bin\\bash.exe", "-c", "sleep 2; echo hello;")
		//cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "sleep 2; echo hello;")
		//cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "docker pull hello-world:latest")
		cmd = exec.CommandContext(ctx, "/bin/bash", "-c", "ansible all -m ping")

		// 执行任务，捕获输出
		output, err = cmd.CombinedOutput()

		// 把任务输出结果，传给main协程
		resultChan <- &result{
			err:    err,
			output: output,
		}
	}()

	// 继续往下走
	//time.Sleep(1 * time.Second) // go协程无法执行完，程序最终输出
	time.Sleep(3 * time.Second) // go协程可以执行完,程序最终输出<nil> hello

	// 取消上下文,程序会杀死与之关联的协程中的子进程
	cancelFunc()

	// 在main协程里，等待子进程的退出，并打印任务执行结果
	res = <-resultChan

	// 打印任务执行结果
	fmt.Println(res.err, string(res.output))
}
