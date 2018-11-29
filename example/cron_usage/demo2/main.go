package main

import (
	"fmt"
	"time"

	"github.com/gorhill/cronexpr"
)

// 定义一个任务结构体
type CronJob struct {
	expr     *cronexpr.Expression // 定时任务表达式
	nextTime time.Time            // 下次调度时间:expr.Next(now)
}

func main() {
	// 需要有一个调度协程，它定时检查所有的Cron定时任务(在调度表里面检查)，谁过期了就执行谁
	var (
		cronJob       *CronJob
		expr          *cronexpr.Expression
		now           time.Time
		scheduleTable map[string]*CronJob // 任务的调度表,key是任务的名字
	)

	scheduleTable = make(map[string]*CronJob)

	// 当前时间
	now = time.Now()

	// 定义第一个cronjob
	expr = cronexpr.MustParse("*/5 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["job1"] = cronJob

	// 定义第二个cronjob
	expr = cronexpr.MustParse("*/6 * * * * * *")
	cronJob = &CronJob{
		expr:     expr,
		nextTime: expr.Next(now),
	}
	// 任务注册到调度表
	scheduleTable["job2"] = cronJob

	// 启动一个调度协程
	go func() {
		var (
			jobName string
			cronJob *CronJob
			now     time.Time
		)
		// 定时检查一下任务调度表
		for {
			now = time.Now()

			for jobName, cronJob = range scheduleTable {
				// 判断是否过期
				if cronJob.nextTime.Before(now) || cronJob.nextTime.Equal(now) {
					// 启动一个协程，执行这个任务
					go func(jobName string) {
						fmt.Println("执行：", jobName)
					}(jobName)

					// 计算下次调度时间(更新调度表信息)
					cronJob.nextTime = cronJob.expr.Next(now)
					fmt.Println(jobName, "下次调度时间：", cronJob.nextTime)
				}
			}

			// 睡眠100毫秒
			select {
			case <-time.NewTimer(100 * time.Millisecond).C: // 将在100毫秒可读，返回
			}
		}
	}()

	time.Sleep(100 * time.Second)
}
