package main

import (
	module "YServer/Logic/Module"
	"YServer/Logic/SeamlessMapManager"
)

func MapManagerLogicRegister() {
	module.Register("MapManager", SeamlessMapManager.New())

}
