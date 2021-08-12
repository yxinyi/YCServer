package SeamlessMapManager

import (
	"YMsg"
	"YNet"
	ylog "YServer/Logic/Log"
	module "YServer/Logic/Module"
	user "YServer/Logic/User"
	"github.com/yxinyi/YEventBus"
	"time"
)

var G_seam_map_manager *SeamlessMapManager
type SeamlessMapManager struct {
	module.ModuleBase
	M_user_manager map[uint64]*user.User
}

func New() *SeamlessMapManager {
	G_seam_map_manager = &SeamlessMapManager{
	}
	return G_seam_map_manager
}

func userLogin(session_ *YNet.Session) {

}

func userOffline(session_ *YNet.Session) {

}


func userEnterMap(session_ *YNet.Session) {
	//判断当前玩家在哪个地图上,地图为固定大小故很容易算出

	//然后根据同步范围计算出当前玩家需要注册到哪些 AOI 监视列表中(故使用灯塔AOI方案),得到 AOI 列表

	//AOI 列表中,包含玩家坐标的AOI 为主AOI,玩家所有操作将转发至包含该AOI 的逻辑进程上,其他 AOI 则为拷贝数据,当主进程上的玩家数据变化后需要同步至其他进程

	//如果玩家初始坐标的map 还未开启,则新建一个 map 进程并与当前 manager 进行连接,连接完成后进行同步

	//后续的消息将先发送至本进程,再进行转发,转发至玩家真正所在的进程上
}

func mapCreate() {
	//地图为矩形,通过参数可以得知当前矩形地图的范围

	//通过范围可以计算出当前地图 AOI 灯塔

	//创建 map 进程,并将本机地址传递,记录新进程的信息,建立对应GRPC通道

	
}

func (b *SeamlessMapManager) Init() error {
	YEventBus.Register("UserLogin", userLogin)
	YEventBus.Register("UserOffline", userOffline)

	YNet.Register(YMsg.MESSAGE_TEST,func(msg_ YMsg.Message,s_ YNet.Session){
		ylog.Info("MESSAGE_TEST [%v] ", msg_)
	})
	return nil
}

func (b *SeamlessMapManager)Update(time_ time.Time)  {
	//ylog.Info("time [%v] ",time_.Unix())
}























