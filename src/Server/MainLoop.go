package main

import (
	"YNet"
	"YServer/Logic/Log"
	module "YServer/Logic/Module"
	"YTimer"
	"fmt"
	"github.com/yxinyi/YEventBus"
	"math/rand"
	"syscall"
	"time"
	"unsafe"
)


func setTitle(title string) {
	kernel32, _ := syscall.LoadLibrary(`kernel32.dll`)
	sct, _ := syscall.GetProcAddress(kernel32, `SetConsoleTitleW`)
	syscall.Syscall(sct, 1, uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(title))), 0, 0)
	syscall.FreeLibrary(kernel32)
}

func MainLoop() {
	rand.Seed(time.Now().Unix())
	SingleLogicRegister()
	ylog.Info("start module init ")
	YTimer.NewWheelTimer(YTimer.WheelSlotCount)
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
	_time_tick := time.Tick(10 * time.Millisecond)
	_last_time := time.Now()
	_tick_cout := 0

	for {
		select {
		case _time := <-_time_tick:
			module.Update(_time)
			_tick_cout++
			if _time.Sub(_last_time).Seconds() >= 1 {
				//ylog.Info("[%v] tick count [%v]", _time.String(), _tick_cout/10)
				title := fmt.Sprintf("fps [%v]",_tick_cout)
				setTitle(title)
				_tick_cout = 0
				_last_time = _time
			}
		case <-YNet.G_close:
			err = module.Stop()
			if err != nil {
				panic(" module Stop err")
			}
			return
		case _timer_list := <-YTimer.G_call:
			YTimer.Loop(_timer_list)
		case msg := <-YNet.G_net_msg_chan:
			switch msg.M_msg_type {
			case YNet.NET_SESSION_STATE_CONNECT:
				YEventBus.Send("UserLogin", msg.M_session)
			case YNet.NET_SESSION_STATE_MSG:
				_err := YNet.Dispatch(msg.M_session, msg.M_net_msg)
				if _err != nil {
					ylog.Erro("SessionID [%v] msg id [%v] [%v]", msg.M_session.GetUID(), msg.M_uid, _err.Error())
				}
			case YNet.NET_SESSION_STATE_CLOSE:
				YEventBus.Send("UserOffline", msg.M_session)
			}
		}
	}

}
