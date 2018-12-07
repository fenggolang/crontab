package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/fenggolang/crontab/worker"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// worker -config ./worker.json
	// worker -h
	flag.StringVar(&confFile, "config", "./worker.json", "worker.json")
	flag.Parse()
}

// Go 1.5 版本之前，默认使用的是单核心执行。从 Go 1.5 版本开始，默认执行下面语句以便让代码并发执行，最大效率地利用 CPU。(所以go1.5以上不用再指定)
// 初始化线程数量
//func initEnv() {
//	runtime.GOMAXPROCS(runtime.NumCPU())
//}

func main() {
	var (
		err error
	)

	// 初始化命令行参数
	initArgs()

	// 初始化线程
	//initEnv()

	// 加载配置
	if err = worker.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 服务注册
	if err = worker.InitRegister(); err != nil {
		goto ERR
	}

	// 启动日志器：日志协程
	if err = worker.InitLogSink(); err != nil {
		goto ERR
	}

	// 启动执行器：执行协程
	if err = worker.InitExecutor(); err != nil {
		goto ERR
	}

	// 启动调度器：调度协程
	if err = worker.InitScheduler(); err != nil {
		goto ERR
	}

	// 初始化任务管理器：监听协程
	if err = worker.InitJobMgr(); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}

ERR:
	fmt.Println(err)
}
