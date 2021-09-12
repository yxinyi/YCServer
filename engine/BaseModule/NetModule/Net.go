package NetModule

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YModule"
	"github.com/yxinyi/YCServer/engine/YMsg"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YNode"
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
	/*err := YNet.ListenTcp4("127.0.0.1:20000")
	if err != nil {
		panic(" ListenTcp4 err")
	}*/

}

func (m *NetModule) RPC_Listen(ip_port_ string) {
	err := YNet.ListenTcp4(ip_port_)
	if err != nil {
		panic(" ListenTcp4 err")
	}
	ylog.Info("[NetModule] Start Listen [%v]", ip_port_)
}

func (m *NetModule) RPC_Connect(ip_ string, port_ string) {
	//后续改成异步的,连接会阻塞当前模块
	_new_connect := YNet.NewConnect()
	_new_connect.Connect(ip_, port_)
	_new_connect.Start()
}

func (m *NetModule) RPC_SendNetMsgJson(s_ uint64, msg_ *YNet.NetMsgPack) {
	_session := m.m_session_pool[s_]
	if _session == nil {
		return
	}
	ylog.Info("[NetModule:SendNetMsgJson] [%v]",msg_.M_msg_name)
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

func (m *NetModule) Loop() {
	for more := true; more; {
		select {
		case _msg := <-YNet.G_net_msg_chan:
			switch _msg.M_msg_type {
			case YNet.NET_SESSION_STATE_CONNECT:
				m.m_session_pool[_msg.M_session.GetUID()] = _msg.M_session
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
	m.Info.Loop()
}

func (m *NetModule) Close() {
	YNet.Stop()
}
