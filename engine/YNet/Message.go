package YNet

import "sync/atomic"

const (
	NET_SESSION_STATE_CONNECT uint8 = iota
	NET_SESSION_STATE_MSG
	NET_SESSION_STATE_CLOSE
)

type Message struct {
	M_msg_type uint8
	M_uid      uint64
	M_session  *Session
	M_net_msg  *NetMsgPack
}

var g_msg_unique_id = uint64(0)

func NewMessage(type_ uint8, s_ *Session, pack_ *NetMsgPack) *Message {
	atomic.AddUint64(&g_msg_unique_id, 1)
	return &Message{
		M_msg_type: type_,
		M_uid:      g_msg_unique_id,
		M_session:  s_,
		M_net_msg:  pack_,
	}
}
