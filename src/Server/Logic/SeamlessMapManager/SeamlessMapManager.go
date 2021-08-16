package SeamlessMapManager

import (
	"YMsg"
	"YNet"
	ylog "YServer/Logic/Log"
	user "YServer/Logic/User"
	"github.com/yxinyi/YEventBus"
	"time"
)

var G_seam_map_manager *SeamlessMapManager

func New() *SeamlessMapManager {
	G_seam_map_manager = &SeamlessMapManager{
	}
	return G_seam_map_manager
}

func userLogin(session_ *YNet.Session) {

}

func userOffline(session_ *YNet.Session) {

}

func GetMapIndexWithPosition(pos_ YMsg.PositionXY) SeamlessMapIndex {
	return SeamlessMapIndex{}
}

func GetAOIArrWithPosition(pos_ YMsg.PositionXY, view_range_ float64) []SeamlessAOIIndex {
	retArr := make([]SeamlessAOIIndex, 0)
	return retArr
}
func GetMapUidWithAOI(SeamlessAOIIndex) uint32 {
	return 0
}
func GetMapUidWithAOIArr(aoi_arr_ []SeamlessAOIIndex) []uint32 {
	retArr := make([]uint32, 0)
	for _, it := range aoi_arr_ {
		retArr = append(retArr, GetMapUidWithAOI(it))
	}

	return retArr
}

func GetMapSession(map_uid_ uint32) *YNet.Session {
	return nil
}

func CreateMapProcessWithCallBack(map_uid_ uint32, cb_ func()) *YNet.Session {
	return nil
}

func userEnterMap(user_ *user.User) {
	//判断当前玩家在哪个地图上,地图为固定大小故很容易算出
	//_map_idx := GetMapIndexWithPosition(user_.M_pos)
	//然后根据同步范围计算出当前玩家需要注册到哪些 AOI 监视列表中(故使用灯塔AOI方案),得到 AOI 列表
	_aoi_arr := GetAOIArrWithPosition(user_.M_pos, user_.M_view_range)
	//AOI 列表中,包含玩家坐标的AOI 为主AOI,玩家所有操作将转发至包含该AOI 的逻辑进程上,其他 AOI 则为拷贝数据,当主进程上的玩家数据变化后需要同步至其他进程
	_sync_map_list := GetMapUidWithAOIArr(_aoi_arr)
	//如果玩家初始坐标的map 还未开启,则新建一个 map 进程并与当前 manager 进行连接,连接完成后进行同步
	for _, it := range _sync_map_list {
		_map_sess := GetMapSession(it)
		_enter_map_func := func() {

		}
		if _map_sess == nil {
			CreateMapProcessWithCallBack(it, _enter_map_func)
		} else {
			_enter_map_func()
		}
	}

	////后续的消息将先发送至本进程,再进行转发,转发至玩家真正所在的进程上
}

func mapCreate() {
	//地图为矩形,通过参数可以得知当前矩形地图的范围

	//通过范围可以计算出当前地图 AOI 灯塔

	//创建 map 进程,并将本机地址传递,记录新进程的信息,建立对应GRPC通道

}

func (b *SeamlessMapManager) Init() error {
	YEventBus.Register("UserLogin", userLogin)
	YEventBus.Register("UserOffline", userOffline)

	YNet.Register(YMsg.MESSAGE_TEST, func(msg_ YMsg.Message, s_ YNet.Session) {
		ylog.Info("MESSAGE_TEST [%v] ", msg_)
	})
	return nil
}

func (b *SeamlessMapManager) Update(time_ time.Time) {
	//ylog.Info("time [%v] ",time_.Unix())
}
