#### (一): 监听协程
- 利用watch API,监听/cron/jobs/(保存和删除)和/cron/killer/(强杀)目录的变化;
- 将变化事件通过channel推送给**调度协程**,更新内存中的任务信息.

```go
// 任务管理器
type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}
```
---

#### (二): 调度协程
- 监听任务变更event,更新内存中维护的任务列表;
- 检查任务cron表达式,扫描到期任务,交给**执行协程**运行;
- 监听任务控制event,强制中断正在执行中的子进程;
- 监听任务执行result,更新内存中任务状态,投递执行日志给**日志协程**.

```go
// 任务调度
type Scheduler struct {
	jobEventChan      chan *common.JobEvent              // etcd任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan // 任务调度计划表
	jobExecutingTable map[string]*common.JobExecuteInfo  // 任务执行表
	jobResultChan     chan *common.JobExecuteResult      // 任务结果队列
}
```
---

#### (三): 执行协程
- 在etcd中抢占分布式乐观锁：/cron/lock/任务名;
- 抢占成功则通过Command类执行shell任务;
- 捕获Command输出并等待子进程结束,将执行结果投递给**调度协程**;

```go
// 任务执行器
type Executor struct {
}
```
---

#### (四): 日志协程
- 监听**调度协程**发来的执行日志,放入一个batch中;
- 对新batch执行启动定时器,超时自动提交;
- 若batch被放满，那么立即提交,并取消自动提交定时器.

```go
// mongodb存储日志
type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}
```
---
