package TestModule

import (
	"YMsg"
	"YServer/Frame/YModule"
	ylog "YServer/Logic/Log"
)

type TestInfo struct {
	*YModule.Info
}

func newInfo() *TestInfo {
	_info := &TestInfo{}
	_info.Info = YModule.NewInfo()
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
	ylog.Info("RPC_Test")
	
}

func (m *TestInfo) RPC_Test_2(val_ uint32) {
	ylog.Info("RPC_Test_2 [%v]", val_)
}

func (m *TestInfo) RPC_Test_3(val_ uint32, str_ string) {
	ylog.Info("RPC_Test_3 [%v] [%v]", val_, str_)
}


func (m *TestInfo) RPC_Test_4(param_ YMsg.TestParam) {
	ylog.Info("RPC_Test_4 [%v]",param_)
}