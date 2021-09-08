package TestModule

import (
	ylog "YLog"
	"YModule"
	"YMsg"
	"YNode"
)

type TestInfo struct {
	*YModule.Info
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

func (m *TestInfo) Loop() {
	m.Info.Loop()
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
	m.Info.RPCCallUsingJson("TestModule2",0,"Test_3",3333,"44444")
}


func (m *TestInfo) RPC_Test_4(param_ YMsg.TestParam) {
	ylog.Info("TestModule RPC_Test_4 [%v]",param_)
}