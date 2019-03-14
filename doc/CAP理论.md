## CAP里面(常用于分布式存储)
[三分钟分布式CAP理论就这么复杂
](https://www.toutiao.com/a6632502751622332942/?tt_from=mobile_qq&utm_campaign=client_share&timestamp=1544267186&app=news_article_lite&utm_source=mobile_qq&iid=53283503090&utm_medium=toutiao_android&group_id=6632502751622332942)
### C(Consistency): 
```markdown
一致性: 写入后立即读到新值

mysql就是强一致性的

```
### A(Availability):
```markdown
可用性: 通常保障最终一致(不会因为单点故障就不可用了),通常说一个服务高可用，就是避免单点故障,mysql就不是一个高可用的，mysql
做了一个折中，向可用性倾斜，支持主从同步，master挂了之后，还可以提供读的可用性，只是写的可用性丢失
```
### P(Partition):
```markdown
分区容错性(partition tolerance)：分布式必须面对网络分区，比如两个节点之间网线被踩断，光缆被挖断就出现分区了，对于分布式系统，网络之间发生调用时非常常用的

所以在CAP里面，P是一定要实现的，你不可能因为说网络不通了就不服务了

mysql就是以CP为主，然后向A(可用性)做了一个倾斜，只是主从同步的异步复制,这样主库down之后，写会受影响，但是读库可以继续服务

```
---
## BASE理论(常用于应用架构)
### 