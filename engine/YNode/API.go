package YNode

func (n *Info) RPC_NewModule(create_func_name string, uid_ uint64) {
	
	_new_module := n.m_moduele_factory[create_func_name](obj, uid_)
	obj.register(_new_module)
	go obj.startModule(_new_module)
}

func (n *Info) RPC_RegisterOtherNode(ip_port_ string, session_ uint64) {
	_node_id, exists := n.M_node_ip_to_id[ip_port_]
	if !exists {
		n.RPCCall("NetModule", uint64(n.M_node_id), "Close", session_)
		return
	}
	n.M_node_id_to_session[_node_id] = session_
}
