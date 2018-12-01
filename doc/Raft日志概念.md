### raft协议原理
- 每一个任期内只能有一个leader
- leader只能追加日志条目，不能重写或者删除日志条目
- 如果两个日志条目的index和term都相同，则两个如果日志中，两个条目及它们之前的日志条目也完全相同
- 如果一条日志被commited过，那么大于该日志条目任期的日志都应该包含这个点
- 如果一个server将某个特定index的日志条目交由状态机处理了，那么对于其他server，交由状态及处理的log中相同index的日志条目应该相同\
---
#### 如何保证选举安全
raft中，只要candidate获得多数投票，就可以成为领导人。follower在投票的时候遵循两个条件：

先到先得
cadidate的term大于follower的term，或者两个term相等，但cadidate的index大于follower的index
　　对于选举结果：

如果票被瓜分，产生两个leader，这次选举失效，进行下一轮的选举
只有一个leader的情况是被允许的
　　这里重要的一点是：如何保证在有限的时间内确定出一个候选人，而不会总是出现票被瓜分的情况？raft使用了一个比较优雅的实现方式，
随机选举超时(randomize election timeouts)。这就使得每个server的timeout不一样，发起新一轮选举的时候，有些server还不是
voter;或者一些符合条件的candidate还没有参加下一轮。这种做法使得单个leader会很快被选举出来。
- 选举leader需要半数以上节点参与
- 节点commit日志最多的允许选举为leader
- commit日志同样多，则term,index越大的允许选举成为leader

#### 如何保证日志匹配
　　Leader在进行AppendEntry RPCs的时候，这个消息中会携带preLogIndex和preLogTerm这两个信息，follower收到消息的时候，首先
判断它最新日志条目的index和term是否和rpc中的一样，如果一样，才会append.

　　这就保证了新加日志和其前一条日志一定是一样的。从第一条日志起就开始遵循这个原理，很自然地可以作出这样的推断。

### raft二阶段协议
- 第一阶段是日志复制，一旦复制到大多数follower成功，leader就可以本地提交(告诉客户端提交或者写入数据成功)
- 第二阶段是集群剩余没有从leader得到日志复制的follower的日志复制
- 俗称抽屉理论
- raft是一致性协议，是用来保障servers上副本一致性的一种算法。
---
### raft日志概念
- replication: 日志在leader生成，向follower复制，达到各个节点的日志序列最终一致
- term: 任期，重新选举产生的leader,其term单调递增
- log index: 日志行在日志序列的下标