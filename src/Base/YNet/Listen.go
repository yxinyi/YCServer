package YNet

import (
	"errors"
	"fmt"
	"net"
)

var g_listener net.Listener
var G_close = make(chan struct{})
var g_connect_list = make(map[uint32]*Session)

var G_net_msg_chan = make(chan *Message,100)


func ListenTcp4(ip_port_ string) error {
	listen, err := net.Listen("tcp4", ip_port_)
	if err != nil {
		return errors.New("listen err [%v]" + err.Error())
	}
	go func() {
		for {
			select {
			case <-G_close:
				return
			default:
				_conn, err := listen.Accept()
				if err != nil {
					fmt.Printf("accept err " + err.Error())
					continue
				}
				_session := NewSession(_conn)
				g_connect_list[_session.m_uid] = _session
				_session.StartLoop()

				connMsg := NewMessage(NET_SESSION_STATE_CONNECT,_session,nil)
				G_net_msg_chan <- connMsg
			}

		}
	}()

	return nil
}

func Stop() {
	G_close <- struct{}{}
}
