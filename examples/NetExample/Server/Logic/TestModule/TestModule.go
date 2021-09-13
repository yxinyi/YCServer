package TestModule

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
	"github.com/yxinyi/YCServer/examples/NetExample/Msg"
)

type TestInfo struct {
	YModule.BaseInter
}

func NewInfo(node_ *YNode.Info) *TestInfo {
	_info := &TestInfo{}
	_info.Info = YModule.NewInfo(node_)
	
	return _info
}

func (m *TestInfo) GetInfo() *YModule.Info {
	return m.Info
}
func (m *TestInfo) Init() {
	m.Info.Init(m)
}

func (m *TestInfo) Close() {

}

func (m *TestInfo) MSG_C2S_TestMsg(s_ uint64, msg_ Msg.C2S_TestMsg) {
	ylog.Info("TestModule[%v]", msg_)
}

func (m *TestInfo) MSG_C2S_TestMsg_2(s_ uint64, msg_ Msg.C2S_TestMsg_2) {
	ylog.Info("TestModule [%v]", msg_)
}
