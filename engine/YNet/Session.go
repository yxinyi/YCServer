package YNet

import (
	"fmt"
	"net"
	"sync/atomic"
)

type Session struct {
	m_uid           uint64
	m_conn          net.Conn
	m_stop          chan struct{}
	m_send_msg_chan chan *NetMsgPack
	M_is_rotbot     bool
}

var g_uni_id = uint64(0)

func NewSession(conn_ net.Conn) *Session {
	atomic.AddUint64(&g_uni_id, 1)
	return &Session{
		m_uid:           g_uni_id,
		m_conn:          conn_,
		m_send_msg_chan: make(chan *NetMsgPack, MAX_PACKAGE_LEN),
		m_stop:          make(chan struct{}),
	}
}

func (s *Session) GetUID() uint64 {
	return s.m_uid
}

func (s *Session) Close() {
	close(s.m_send_msg_chan)
	close(s.m_stop)
}

func (s *Session) SendJson(json_ interface{}) error {
	if s.M_is_rotbot {
		return nil
	}
	_msg := NewNetMsgPackWithJson(json_)
	if _msg == nil {
		return fmt.Errorf("[Session:SendJson] pack error")
	}
	s.Send(_msg)
	return nil
}

func (s *Session) Send(msg_ *NetMsgPack) bool {
	if s.M_is_rotbot {
		return true
	}
	if s.m_send_msg_chan == nil {
		return false
	}
	s.m_send_msg_chan <- msg_
	return true
}

func (s *Session) StartLoop() {
	go func() {
		for {
			select {
			case <-s.m_stop:
				return
			case pack, ok := <-s.m_send_msg_chan:
				if !ok {
					s.m_conn.Close()
					return
				}
				_msg_byte := pack.ToByteStream()
				_, err := s.m_conn.Write(_msg_byte)
				if err != nil {
					break
				}
			}
		}
	}()
	go func() {
		for {
			select {
			case <-s.m_stop:
				return
			default:
				_msg_pack := NewNetMsgPack()
				if !_msg_pack.InitFromIO(s.m_conn) {
					connMsg := NewMessage(NET_SESSION_STATE_CLOSE, s, nil)
					G_net_msg_chan <- connMsg
					return
				}

				connMsg := NewMessage(NET_SESSION_STATE_MSG, s, _msg_pack)
				G_net_msg_chan <- connMsg
			}
		}
	}()
}
