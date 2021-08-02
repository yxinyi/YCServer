package YTimer

import "time"

type TimerCallBack func()

var g_timer_manager TimerManagerInter
var g_add_timer_channel = make(chan *Timer, 100)
var g_cancel_timer = make(chan uint32, 100)
var G_close = make(chan *Timer, 100)
var G_call = make(chan []*Timer, 100)

const (
	WheelSlotCount       = 10000
	TickerTime     int64 = int64(time.Millisecond * 100)
)

type WheelTimer struct {
	m_slots    []*Timer
	m_map      map[uint32]*Timer
	m_cursor   uint32
	m_cur_time time.Time
}

func init() {
	go func() {
		ticker := time.Tick(time.Millisecond * 100)
		for {
			select {
			case _time := <-ticker:
				g_timer_manager.setTime(_time)
				G_call <- g_timer_manager.getAllCall()
			case _timer := <-g_add_timer_channel:
				g_timer_manager.timeCall(_timer)
			case _cancel_uid := <-g_cancel_timer:
				g_timer_manager.cancelTimer(_cancel_uid)
			case <-G_close:
				return
			}
		}
	}()
}

type TimerManagerInter interface {
	timeCall(t_ *Timer)
	getAllCall() []*Timer
	setTime(t_ time.Time)
	cancelTimer(uint32)
}

func NewWheelTimer() {
	g_timer_manager = &WheelTimer{
		m_slots: make([]*Timer, WheelSlotCount),
		m_map:   make(map[uint32]*Timer),
	}
}

func (t *WheelTimer) getSlot(nano_timestamp int64) uint32 {
	_diff_tm := nano_timestamp - t.m_cur_time.UnixNano()
	_future_slot := ((_diff_tm / TickerTime) + int64(t.m_cursor)) % WheelSlotCount

	return uint32(_future_slot)
}
func (t *WheelTimer) cancelTimer(t_ uint32) {
	_timer, exists := t.m_map[t_]
	if !exists {
		return
	}
	if _timer.m_perv != nil {
		_timer.m_perv.m_next = _timer.m_next
	}
	if _timer.m_next != nil {
		_timer.m_next.m_perv = _timer.m_perv
	}
	_timer.m_perv = nil
	_timer.m_next = nil
}

func (t *WheelTimer) setTime(t_ time.Time) {
	t.m_cur_time = t_
}

func (t *WheelTimer) insertSlot(slot_ uint32, t_ *Timer) {
	t.m_map[t_.m_uid] = t_
	_frst_time := t.m_slots[slot_]
	t_.m_next = _frst_time
	if _frst_time != nil {
		_frst_time.m_perv = t_
	}
	t.m_slots[slot_] = _frst_time
}

func (t *WheelTimer) timeCall(t_ *Timer) {
	if t_.M_call_time < t.m_cur_time.UnixNano() {
		t.insertSlot(t.m_cursor+1, t_)
		return
	}
	t.insertSlot(t.getSlot(t_.M_call_time), t_)
}

func (t *WheelTimer) getNextCursor() uint32 {
	t.m_cursor++
	return t.m_cursor % WheelSlotCount
}

func (t *WheelTimer) getAllCall() []*Timer {
	_timer_list := make([]*Timer, 0)

	_first_timer := t.m_slots[t.getNextCursor()]

	_now_nano_time := t.m_cur_time.UnixNano()
	for _first_timer != nil {
		if _first_timer.M_call_time > _now_nano_time {
			continue
		}
		_timer_list = append(_timer_list, _first_timer)
		if _first_timer.m_perv != nil {
			_first_timer.m_perv.m_next = _first_timer.m_next
		}
		if _first_timer.m_next != nil {
			_first_timer.m_next.m_perv = _first_timer.m_perv
		}
		_first_timer = _first_timer.m_next
	}

	return _timer_list
}
