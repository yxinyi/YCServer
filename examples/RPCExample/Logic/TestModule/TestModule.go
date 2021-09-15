package TestModule

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YNode"
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

func (m *TestInfo) RPC_Test() {
	ylog.Info("TestModule RPC_Test")
	
}

func (m *TestInfo) RPC_Test_2(val_ uint32) {
	ylog.Info("TestModule RPC_Test_2 [%v]", val_)
}

func (m *TestInfo) RPC_Test_3(val_ uint32, str_ string) {
	ylog.Info("TestModule RPC_Test_3 [%v] [%v]", val_, str_)
	var _func func()
	_func = func() {
		m.Info.RPCCall("TestModule2", 0, "Test", func() {
			ylog.Info("Test 回调")
			_func()
		})
	}
	_func()
}

func (m *TestInfo) RPC_Test_4(param_ int) {

}




