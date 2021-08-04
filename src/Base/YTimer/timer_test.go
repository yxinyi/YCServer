package YTimer

import (
	"testing"
	time "time"
)

func init(){
	NewWheelTimer(10)
}

func TestTimerIndexForeach(t_ *testing.T) {
	for _idx := 0; _idx < 20; _idx++ {
		t_.Logf("[%v]", g_timer_manager.getNextCursor())
	}
}
func TimerApiHelp(t_ *testing.T, second_ uint32){
	_before_time := time.Now()
	_close := make(chan struct{})
	AfterSecondsCall(second_, func(tick_time_ time.Time) {
		_diff_time := tick_time_.Sub(_before_time)
		if uint32(_diff_time.Seconds()) != second_ {
			t_.Errorf("err diff [%v]", int(_diff_time.Seconds()))
		}
		t_.Logf("info diff [%v]", int(_diff_time.Seconds()))
		close(_close)
	})
	for {
		select {
		case _t := <-G_call:
			Loop(_t)
		case <-_close:
			return
		}
	}
}
func TestTimerApi(t_ *testing.T) {
	TimerApiHelp(t_,1)
	TimerApiHelp(t_,2)
	TimerApiHelp(t_,3)
}
