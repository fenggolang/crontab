package main

import (
	"flag"
	"fmt"
	"runtime"
	"time"

	"github.com/fenggolang/crontab/master"
)

var (
	confFile string // 配置文件路径
)

// 解析命令行参数
func initArgs() {
	// master -config ./master.json -xxx 123 -yyy ddd
	// master -h
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}

// 配置线程数量和核心数相同，因为go的协程是被调度到线程上，而线程是操作系统的概念，线
// 程的多少决定的并发的好坏，配置线程数量和cpu核心数量相等，发挥最大作用
// 初始化线程数量
// Go 1.5 版本之前，默认使用的是单核心执行。从 Go 1.5 版本开始，默认执行下面语句以便让代码并发执行，最大效率地利用 CPU。(所以go1.5以上不用再指定)
func initEnv() {
	runtime.GOMAXPROCS(runtime.NumCPU())
}

func main() {
	var (
		err error
	)

	// 初始化命令行参数
	initArgs()

	// 初始化线程
	initEnv()

	// 加载配置
	if err = master.InitConfig(confFile); err != nil {
		goto ERR
	}

	// 初始化服务发现模块
	if err = master.InitWorkerMgr(); err != nil {
		goto ERR
	}

	// 日志管理器
	if err = master.InitLogMgr(); err != nil {
		goto ERR
	}

	//  任务管理器
	if err = master.InitJobMgr(); err != nil {
		goto ERR
	}

	// 启动Api HTTP服务
	if err = master.InitApiServer(); err != nil {
		goto ERR
	}

	// 正常退出
	for {
		time.Sleep(1 * time.Second)
	}
ERR:
	fmt.Println(err)
}
