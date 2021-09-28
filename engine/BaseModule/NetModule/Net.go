package NetModule

import (
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YNode"
	"time"
)

type NetModule struct {
	YModule.BaseInter
	m_session_pool map[uint64]*YNet.Session
	m_net_msg_pool map[string][]YMsg.Agent
}

func NewInfo(node_ *YNode.Info) *NetModule {
	_info := &NetModule{
		m_session_pool: make(map[uint64]*YNet.Session),
		m_net_msg_pool: make(map[string][]YMsg.Agent),
	}
	_info.Info = YModule.NewInfo(node_)
	return _info
}
func (m *NetModule) Init() {
	m.Info.Init(m)
}



func (m *NetModule) Loop_100(time time.Time) {
	for more := true; more; {
		select {
		case _msg := <-YNet.G_net_msg_chan:
			switch _msg.M_msg_type {
			case YNet.NET_SESSION_STATE_CONNECT:
				m.m_session_pool[_msg.M_session.GetUID()] = _msg.M_session
			case YNet.NET_SESSION_STATE_CONNECT_OTHER_SUCCESS:
				_ip_port := _msg.M_net_msg.M_msg_name
				m.Info.RPCCall("YNode", uint64(m.Info.M_node_id), "RegisterOtherNode", _ip_port, _msg.M_session.GetUID())
			case YNet.NET_SESSION_STATE_MSG:
				for _, _agent_it := range m.m_net_msg_pool[_msg.M_net_msg.M_msg_name] {
					_net_msg := &YMsg.C2S_net_msg{
						_agent_it,
						_msg.M_session.GetUID(),
						_msg.M_net_msg,
					}
					m.NetToOther(_net_msg)
				}
			case YNet.NET_SESSION_STATE_CLOSE:
				delete(m.m_session_pool, _msg.M_session.GetUID())
			}
		default:
			more = false
		}
	}
}

func (m *NetModule) Close() {
	YNet.Stop()
}
