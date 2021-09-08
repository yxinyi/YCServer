package main

import (
	maze_map "YServer/Logic/Map"
	module "YServer/Logic/Module"
	user "YServer/Logic/User"
)

func SingleLogicRegister() {
	module.Register("MazeMapManager", maze_map.NewMazeMapManager())
	module.Register("UserManager", user.NewModuleUserLogin())
}
