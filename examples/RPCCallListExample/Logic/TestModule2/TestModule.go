package TestModule2

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

func (m *TestInfo) RPC_Test_1() int {
	ylog.Info("[TestModule2:Test_1]")
	return 10
}

func (m *TestInfo) RPC_Test_2(val_ string) string {
	ylog.Info("[TestModule2:Test_2]")
	return val_
}

func (m *TestInfo) RPC_Test_3() {
	ylog.Info("[TestModule2:Test_3]")
}

func (m *TestInfo) RPC_Test_4() {
	ylog.Info("[TestModule2:Test_4]")
}

func (m *TestInfo) RPC_Test_5(val_ uint32) {
	ylog.Info("[TestModule2:Test_5] [%v]",val_)
}
