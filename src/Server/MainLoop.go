package main

import (
	"YNet"
	"YServer/Logic/Log"
	module "YServer/Logic/Module"
	"github.com/yxinyi/YEventBus"
	"time"
)

type Message struct {
	Id  int    `json:id`
	Str string `json:str`
}

func MainLoop() {
	ylog.Info("start module init ")
	err := module.Init()
	if err != nil {
		panic(" module init err")
	}

	err = module.Start()
	if err != nil {
		panic(" module Start err")
	}

	err = YNet.ListenTcp4("127.0.0.1:20000")
	if err != nil {
		panic(" ListenTcp4 err")
	}
	ylog.Info("start main loop ")
	_time_tick := time.Tick(50 * time.Millisecond)
	for {
		select {
		case _time := <-_time_tick:
			module.Update(_time)
		case <-YNet.G_close:
			err = module.Stop()
			if err != nil {
				panic(" module Stop err")
			}
			return
		case msg := <-YNet.G_net_msg_chan:
			switch msg.M_msg_type {
			case YNet.NET_SESSION_STATE_CONNECT:
				YEventBus.Send("UserLogin", msg.M_session)
			case YNet.NET_SESSION_STATE_MSG:
				_err := YNet.Dispatch(msg.M_session, msg.M_net_msg)
				if _err !=nil{
					ylog.Erro("SessionID [%v] msg id [%v] [%v]", msg.M_session.GetUID(), msg.M_uid,_err.Error())
				}
				//ylog.Info("SessionID [%v] msg id [%v]", msg.M_session.GetUID(), msg.M_uid)
			case YNet.NET_SESSION_STATE_CLOSE:
				YEventBus.Send("UserOffline", msg.M_session)
			}
		}
	}

}
