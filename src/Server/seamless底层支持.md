### 底层支持

无缝地图与普通mmo地图的差异为在地图边缘处,处于两个不同进程的玩家需要看到其他进程在视野范围内的行为,以及在边缘处进行攻击交互等操作.

简单的处理方案为主从对象设计,每个对象会有一个主对象,所有的修改操作都将通过 rpc 进行执行,查询行为只需查阅当前从对象即可,保证进程间最少的消息通信交互.

基于此方案可以简单得出,所有地图元素,地图道具,怪物,玩家等都需满足这个需求,从设计的角度出发,我们可以设定一个基类,后续所有对象都进行继承即可,也就是将所有对象都抽象成一个 entity 对象,后续基于 entity 的基础上进行细节逻辑填充.

### 需要rpc以支持服务调用

#### 服务注册发现:

zookeeper ,redis等都可以实现

服务注册目的为当新的服务器上线之后,当前运行的所有服务器不用进行重启操作即可直接与新服务器进行交互,

逻辑为

新服务器上线后将自己注册到服务器注册模块

然后服务器注册模块通知当前在线的服务器

在线的服务器与新连接服务器进行连接

将新服务器上注册的服务进行拉取,将自身的服务进行推送

完成连接

#### 单点服务:

意味着当前服务唯一,注册消息时候无需进行路由操作,且在服务器发现注册时也可以迅速判断是否已经注册

#### 集群服务:

集群服务可以分为 无状态集群, 有状态集群

无状态集群可以在调用时进行负载均衡,如权重,轮询,负载等逻辑

有状态集群需要在调用时将每次请求分发到正确的目标服务器,所以每个服务需要一个唯一标识,如果服务上绑定有entity对象,则要在entity上将服务标识进行绑定

#### rpc接口定义:

广播接口

特点是不用对广播结果进行确认

如玩家在地图边缘时需要将当前修改信息同步给附近的地图

```
service.notify([]RoundMapSession,"ServiceName",&UserDiffInfo{})
```

指向性远程调用

不关心结果

如进程A的玩家攻击进程B的玩家导致血量变化

```
service.AsyncRPC([]UserBSession,"Attack",&UserDiffInfo{})
```

远程确定性修改

如进程A的玩家攻击进程B的玩家导致血量变化后需要根据回调进行自身BUFF处理

```
service.AsyncCBRPC([]UserBSession,"Attack",&UserDiffInfo{},fun(success error){
	//callback do add some buffer to self
})
```

### ServiceMananger 结构

```
type ServiceRPCFunc struct{
	M_rpc_name string
	M_func reflect.Value
    M_param reflect.Type
}
type ServiceManager struct{
	*Session
	M_func_pool map[string]*ServiceRPCFunc
    M_entity_pool map[uint64]*interface{}
}
func (mgr *ServiceManager)GetServiceList()[]string{}
```

### ConnectManager

```

```























