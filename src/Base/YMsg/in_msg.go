package YMsg

type ModuleInfo struct {
	M_name    string
	M_uid     uint64
	M_node_id uint64
}

type S2S_rpc_msg struct {
	M_uid            uint64
	M_from           ModuleInfo
	M_tar            ModuleInfo
	M_marshal_type   uint32
	M_func_name      string
	M_func_parameter [][]byte
}

type TestParam struct {
	M_uint64 uint64
	M_str    string
	M_slice  []int
}
