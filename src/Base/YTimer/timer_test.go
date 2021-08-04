package YTimer

import (
	"testing"
	time "time"
)

func TestTimerIndexForeach(t_ *testing.T) {
	NewWheelTimer(10)
	for _idx := 0; _idx < 20; _idx++ {
		t_.Logf("[%v]", g_timer_manager.getNextCursor())
	}
}

func TestTimerApi(t_ *testing.T) {
	NewWheelTimer(10)
	_before_time := time.Now()
	_close := make(chan struct{})
	AfterSecondsCall(1, func(tick_time_ time.Time) {
		_diff_time := tick_time_.Sub(_before_time)
		if int(_diff_time.Seconds()) != 1 {
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
