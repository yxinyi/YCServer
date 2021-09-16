package YNode

import (
	"github.com/yxinyi/YCServer/engine/YModule"
)


func RegisterToFactory(name_ string, func_ func(er YModule.RemoteNodeER, uid_ uint64) YModule.Inter) {
	obj.registerToFactory(name_,func_)
}

func (n *Info) registerToFactory(name_ string, func_ func(er YModule.RemoteNodeER, uid_ uint64) YModule.Inter) {
	n.m_moduele_factory[name_] = func_
}

func (n *Info) RPC_ModuleRegister(create_func_name string, uid_ uint64) {
	
	_new_module := n.m_moduele_factory[create_func_name](obj,uid_)
	obj.register(_new_module)
	go obj.startModule(_new_module)
}
