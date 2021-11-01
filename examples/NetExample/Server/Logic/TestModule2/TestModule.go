package TestModule2

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/NetExample/Msg"
)

type Info struct {
	YModule.BaseInter
}

func NewInfo(node_ *YNode.Info, uid_ uint64) YModule.Inter {
	_info := &Info{}
	_info.Info = YModule.NewInfo(node_)
	return _info
}

func (m *Info) Init() {
	m.Info.Init(m)
}

func (m *Info) Loop() {
	m.Info.Loop_Msg()
}
func (m *Info) Close() {
}

func (m *Info) MSG_C2S_TestMsg(s_ uint64, msg_ Msg.C2S_TestMsg) {
	ylog.Info("TestModule2[%v]", msg_)
}

func (m *Info) MSG_C2S_TestMsg_2(s_ uint64, msg_ Msg.C2S_TestMsg_2) {
	ylog.Info("TestModule2 [%v]", msg_)
}
