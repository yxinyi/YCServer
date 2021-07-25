package YNet

import (
	"encoding/json"
	"net"
	"sync"
	"time"
)

const (
	MAX_PACKAGE_LEN uint32 = 1000
)

type Connect struct {
	m_conn    net.Conn
	m_session *Session
	m_wg      sync.WaitGroup
}

func NewConnect() *Connect {
	return &Connect{}
}

func (c *Connect) Connect(ip_, port_ string) bool {
	var _conn net.Conn
	for {
		_tmp_conn, _err := net.Dial("tcp4", ip_+":"+port_)
		if _err != nil {
			time.Sleep(1*time.Second)
			continue
		}
		_conn = _tmp_conn
		break
	}
	c.m_conn = _conn
	c.m_conn.(*net.TCPConn).SetNoDelay(true)
	c.m_session = NewSession(_conn)
	return true
}

func (c *Connect) SendMsg(msg_id_ uint32, msg_ interface{}) {
	_net_msg := NewNetMsgPack()
	_net_msg.M_msg_id = msg_id_
	json_data, err := json.Marshal(msg_)
	if err == nil {
		_net_msg.M_msg_data = json_data
		_net_msg.M_msg_length = uint32(len(json_data))
	}
	c.m_session.Send(_net_msg)
}

func (c *Connect) Start() bool {
	c.m_session.StartLoop()
	return true
}

func (c *Connect) End() bool {
	c.m_session.Close()
	return true
}
