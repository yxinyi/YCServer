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
	m.Info.RPCCall("TestModule2", 0, "Test_1", func(val int) {
		ylog.Info("Test 回调 返回值 [%v]", val)
	}).AfterRPC("TestModule2", 0, "Test_2", "测试值", func(val string) {
		ylog.Info("Test_2 回调 返回值 [%v]", val)
	}).AfterRPC("TestModule2", 0, "Test_3",
	).AfterRPC("TestModule2", 0, "Test_4", func() {
		ylog.Info("Test_4 回调 ")
	}).AfterRPC("TestModule2", 0, "Test_5",56)
	
	m.Info.RPCCall("TestModule2", 0, "Test_1", func(val int) {
		ylog.Info("Test_1 取消后续调用链 [%v]", val)
		m.Info.CancelCBList()
	}).AfterRPC("TestModule2", 0, "Test_5", func() {
		ylog.Info("[错误] 没有取消 ")
	})
	
}

func (m *TestInfo) RPC_Test_4(param_ int) {

}
