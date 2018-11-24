### 保存一个job
```bash
# postman 发起post请求，发起一个job任务
POST http://localhost:8080/job/save
header Content-Type: application/x-www-form-urlencoded
Body key:job value:{"name": "job1","command":"echo hello","cronExpr":"*/5 * * * * * *"}

# 去etcd中查看保存的job
ETCDCTL_API=3 ./etcdctl get "/cron/jobs" --prefix
```
### 删除一个job
```bash
# postman 发起POST请求
POST http://localhost:8080/job/delete
header Content-Type: application/x-www-form-urlencoded
Body key:name value:job1
```
### 查看job列表
```bash
# postman 发起GET请求
GET http://localhost:8080/job/list
header Content-Type: application/x-www-form-urlencoded
```
### 杀死一个job
```bash
# postman 发起POST请求
POST http://localhost:8080/job/kill
header Content-Type: application/x-www-form-urlencoded
Body key:name value:job1

# 然后你去后台 watch etcd的key的变化
[m8@ansible etcd-v3.3.8-linux-amd64]$ ETCDCTL_API=3 ./etcdctl watch "/cron/killer/" --prefix
```

#### 查看任务执行日志
```bash
# postman 发起GET请求
GET http://localhost:8080/job/log?name=job1
```

#### 查看worker节点列表
```bash
# post 发起GET请求
GET http://localhost:8080/worker/list
```