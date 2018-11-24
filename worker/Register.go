package worker

import (
	"context"
	"net"
	"os"
	"strings"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/fenggolang/crontab/common"
)

// 注册节点到etcd： /cron/workers/IP地址
type Register struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease

	localIP       string // 本机IP
	localHostname string // 本机Hostname
}

var (
	G_register *Register
)

// 获取本机网卡IP
func getLocalIP() (ipv4 string, err error) {
	var (
		addrs   []net.Addr
		addr    net.Addr
		ipNet   *net.IPNet // IP地址
		isIpNet bool
	)
	// 获取所有网卡
	if addrs, err = net.InterfaceAddrs(); err != nil {
		return
	}
	// 取第一个非lo的网卡IP
	for _, addr = range addrs {
		// 这个网络地址是IP地址: ipv4, ipv6
		if ipNet, isIpNet = addr.(*net.IPNet); isIpNet && !ipNet.IP.IsLoopback() {
			// 跳过IPV6 ,以及排除掉以.0.1结尾的,这种结尾的一般是k8s集群的pod ip网关地址
			if ipNet.IP.To4() != nil && !strings.HasSuffix(ipNet.IP.String(), ".0.1") {
				ipv4 = ipNet.IP.String() // 192.168.1.1
				return
			}
		}
	}

	err = common.ErrNoLocalIpFound
	return
}

// 获取本机主机名称
func getLocalHostname() (hostname string, err error) {
	if hostname, err = os.Hostname(); err != nil {
		return "", err
	}
	return hostname, nil
}

// 注册到/cron/workers/IP, 并自动续租
func (register *Register) keepOnline() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)

	for {
		// 注册路径
		regKey = common.JobWorkerDir + register.localIP

		cancelFunc = nil

		// 创建租约
		if leaseGrantResp, err = register.lease.Grant(context.TODO(), 10); err != nil {
			goto RETRY
		}

		// 自动续租
		if keepAliveChan, err = register.lease.KeepAlive(context.TODO(), leaseGrantResp.ID); err != nil {
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		if _, err = register.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID)); err != nil {
			goto RETRY
		}

		// 处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { // 续租失败
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}

func InitRegister() (err error) {
	var (
		config        clientv3.Config
		client        *clientv3.Client
		kv            clientv3.KV
		lease         clientv3.Lease
		localIp       string
		localHostname string
	)

	// 初始化配置
	config = clientv3.Config{
		Endpoints:   G_config.EtcdEndpoints,                                     // 集群地址
		DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Millisecond, // 连接超时
	}

	// 建立连接
	if client, err = clientv3.New(config); err != nil {
		return
	}

	// 本机IP
	if localIp, err = getLocalIP(); err != nil {
		return
	}

	// 本地主机名
	if localHostname, err = getLocalHostname(); err != nil {
		return
	}
	// 得到KV和Lease的API子集
	kv = clientv3.NewKV(client)
	lease = clientv3.NewLease(client)

	G_register = &Register{
		client:        client,
		kv:            kv,
		lease:         lease,
		localIP:       localIp,
		localHostname: localHostname,
	}

	// 服务注册
	go G_register.keepOnline()

	return
}
