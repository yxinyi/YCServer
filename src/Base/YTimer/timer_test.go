package YTimer

import (
	"testing"
	time "time"
)

func init() {
	NewWheelTimer(WheelSlotCount)
}

func TestTimerIndexForeach(t_ *testing.T) {
	for _idx := 0; _idx < 20; _idx++ {
		if g_timer_manager.getNextCursor() != (uint32(_idx+1) % g_timer_manager.getSlotSize()) {
			t_.Errorf("[%v]", g_timer_manager.getNextCursor())
		}

	}
}
func TimerApiHelp(t_ *testing.T, second_ float64) {
	_before_time := time.Now()
	_close := make(chan struct{})
	AfterSecondsCall(second_, func() {
		_diff_time := time.Now().Sub(_before_time)
		if _diff_time.Seconds()-second_ > 0.01 {
			t_.Errorf("err diff [%v] right diff [%v] _before_time[%v] after_time [%v]", int(_diff_time.Seconds()), second_, _before_time.Unix(), _diff_time.Seconds())
		}
		t_.Logf("TimerSize [%v]", GetTimerSize())
		if GetTimerSize() != 0 {
			t_.Errorf("err size [%v]", GetTimerSize())
		}

		t_.Logf("diff [%v] right diff [%v] _before_time[%v] after_time [%v]", int(_diff_time.Seconds()), second_, _before_time.Unix(), _diff_time.Seconds())
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
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
	TimerApiHelp(t_, 1)
}
