package NetModule

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
)

func (m *NetModule) RPC_Close(s_ uint64) {
	_session := m.m_session_pool[s_]
	if _session == nil {
		return
	}
	_session.Close()
	delete(m.m_session_pool, s_)
}

func (m *NetModule) RPC_Listen(ip_port_ string) {
	err := YNet.ListenTcp4(ip_port_)
	if err != nil {
		panic(" ListenTcp4 err")
	}
	ylog.Info("[NetModule] Start Listen [%v]", ip_port_)
}

func (m *NetModule) RPC_Connect(ip_port_ string) {
	go func() {
		_new_connect := YNet.NewConnect()
		_new_connect.Connect(ip_port_)
		_conn_sesstion := _new_connect.GetSession()
		_conn_sesstion.StartLoop()
		
		_msg_pack := YNet.NewNetMsgPack()
		_msg_pack.M_msg_name = ip_port_
		_conn_msg := YNet.NewMessage(YNet.NET_SESSION_STATE_CONNECT_OTHER_SUCCESS, _conn_sesstion, _msg_pack)
		YNet.G_net_msg_chan <- _conn_msg
	}()
}

func (m *NetModule) RPC_SendNetMsgJson(s_ uint64, msg_ *YNet.NetMsgPack) {
	_session := m.m_session_pool[s_]
	if _session == nil {
		return
	}
	_session.Send(msg_)
}

func (m *NetModule) RPC_NetMsgRegister(msg_list_ []string, agent_ YMsg.Agent) {
	for _, _msg_it := range msg_list_ {
		_, exists := m.m_net_msg_pool[_msg_it]
		if !exists {
			m.m_net_msg_pool[_msg_it] = make([]YMsg.Agent, 0)
		}
		m.m_net_msg_pool[_msg_it] = append(m.m_net_msg_pool[_msg_it], agent_)
	}
}
