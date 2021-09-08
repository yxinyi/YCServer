package SeamlessMapManager

import (
	module "YServer/Logic/Module"
	user "YServer/Logic/User"
)

type SeamlessMapManager struct {
	module.ModuleBase
	M_user_manager map[uint64]*user.User
}
type SeamlessMapIndex struct {
	M_x int
	M_y int
}


type SeamlessAOIIndex struct {
	M_x int
	M_y int
}

