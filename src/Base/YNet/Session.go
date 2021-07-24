package YNet

import (
	"net"
	"sync/atomic"
)

type Session struct {
	m_uid           uint32
	m_conn          net.Conn
	m_stop          chan struct{}
	m_send_msg_chan chan *NetMsgPack
}

var g_uni_id = uint32(0)

func NewSession(conn_ net.Conn) *Session {
	atomic.AddUint32(&g_uni_id, 1)
	return &Session{
		m_uid:           g_uni_id,
		m_conn:          conn_,
		m_send_msg_chan: make(chan *NetMsgPack, MAX_PACKAGE_LEN),
		m_stop:          make(chan struct{}),
	}
}

func (s *Session) GetUID() uint32 {
	return s.m_uid
}

func (s *Session) Close() {
	close(s.m_send_msg_chan)
	close(s.m_stop)
}

func (s *Session) Send(msg_ *NetMsgPack) bool {
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
			case pack,ok := <-s.m_send_msg_chan:
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
				_msg_pack := NewMessagePack()
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
