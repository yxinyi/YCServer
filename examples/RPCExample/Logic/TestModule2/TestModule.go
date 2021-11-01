package TestModule2

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNode"
)

type TestInfo struct {
	YModule.BaseInter
}

func NewInfo(node_ *YNode.Info, uid_ uint64) YModule.Inter {
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

func (m *TestInfo) RPC_Test() {
	ylog.Info("TestModule2 RPC_Test")
}

func (m *TestInfo) RPC_Test_2(val_ uint32) {
	ylog.Info("TestModule2 RPC_Test_2 [%v]", val_)
}

func (m *TestInfo) RPC_Test_3(val_ uint32, str_ string) float64 {
	ylog.Info("TestModule2 RPC_Test_3 [%v] [%v]", val_, str_)
	m.Info.RPCCall(YMsg.ToAgent("TestModule"), "Test_3", val_+1, "从 Module2 发出")
	return 0.1241423523452342
}
