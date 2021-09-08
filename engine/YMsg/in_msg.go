package YMsg

type ModuleInfo struct {
	M_name    string
	M_uid     uint64
	M_node_id uint64
}

type S2S_rpc_msg struct {
	M_uid            uint64
	M_source         ModuleInfo
	M_tar            ModuleInfo
	M_marshal_type   uint32
	M_func_name      string
	M_func_parameter [][]byte
}

