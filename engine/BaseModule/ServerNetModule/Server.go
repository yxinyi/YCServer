package ServerNetModule

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YNode"
)

const (
	ip   string = "0.0.0.0"
	port string = "10000"
)

func NewInfo(node_ *YNode.Info) *ServerNetModule {
	_info := &ServerNetModule{}
	_info.Info = YModule.NewInfo(node_)

	return _info
}

type ServerNetModule struct {
	YModule.BaseInter
	m_session_pool map[uint64]*YNet.Session
	m_net_msg_pool map[uint32][]YMsg.Agent
}

func (m *ServerNetModule) Init() {
	err := YNet.ListenTcp4("127.0.0.1:20000")
	if err != nil {
		panic(" ListenTcp4 err")
	}
}
func (m *ServerNetModule) RPC_net_msg_register(msg_list_ []uint32, agent_ YMsg.Agent) {
	for _, _msg_it := range msg_list_ {
		_, exists := m.m_net_msg_pool[_msg_it]
		if !exists {
			m.m_net_msg_pool[_msg_it] = make([]YMsg.Agent, 0)
		}
		m.m_net_msg_pool[_msg_it] = append(m.m_net_msg_pool[_msg_it], agent_)
	}
}

func (m *ServerNetModule) Loop() {
	for {
		select {
		case _msg := <-YNet.G_net_msg_chan:
			switch _msg.M_msg_type {
			case YNet.NET_SESSION_STATE_CONNECT:
				m.m_session_pool[_msg.M_session.GetUID()] = _msg.M_session
			case YNet.NET_SESSION_STATE_MSG:
				for _, _agent_it := range m.m_net_msg_pool[_msg.M_net_msg.M_msg_id] {
					_net_msg := &YMsg.C2S_net_msg{
						_agent_it,
						_msg.M_session.GetUID(),
						_msg.M_net_msg,
					}
					m.Info.PushNetMsg(_net_msg)
				}
			case YNet.NET_SESSION_STATE_CLOSE:
				delete(m.m_session_pool, _msg.M_session.GetUID())
			}
		}
	}
}

func (m *ServerNetModule) Close() {
	YNet.Stop()
}
