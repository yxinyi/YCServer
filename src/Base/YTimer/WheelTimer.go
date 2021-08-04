package YTimer

import "time"

type TimerCallBack func(time_ time.Time)

var g_timer_manager timerManagerInter
var g_add_timer_channel = make(chan *Timer, 100)
var g_cancel_timer = make(chan uint32, 100)

var G_call = make(chan *ChanTimer, 100)

var g_close = make(chan struct{})

const (
	WheelSlotCount               = 10000
	TickerTime     time.Duration = time.Second
)

type wheelTimer struct {
	m_slots    []*Timer
	m_map      map[uint32]*Timer
	m_cursor   uint32
	m_cur_time time.Time
}
type ChanTimer struct {
	M_timer_list []*Timer
	M_tick_time  time.Time
}

type timerManagerInter interface {
	timeCall(t_ *Timer)
	getAllCall() *ChanTimer
	setTime(t_ time.Time)
	cancelTimer(uint32)
	close()
	getNextCursor() uint32
	getSlotSize() uint32
	loop(*ChanTimer)
}

func NewWheelTimer(slot_count_ int) {
	_wheel_timer := &wheelTimer{
		m_slots: make([]*Timer, slot_count_),
		m_map:   make(map[uint32]*Timer),
	}
	for _idx := range _wheel_timer.m_slots {
		_wheel_timer.m_slots[_idx] = newTimer()
	}
	g_timer_manager = _wheel_timer
	go func() {
		ticker := time.Tick(TickerTime)
		for {
			select {
			case _time := <-ticker:
				g_timer_manager.setTime(_time)
				G_call <- g_timer_manager.getAllCall()
			case _timer := <-g_add_timer_channel:
				g_timer_manager.timeCall(_timer)
			case _cancel_uid := <-g_cancel_timer:
				g_timer_manager.cancelTimer(_cancel_uid)
			case <-g_close:
				return
			default:
			}
		}
	}()
}

func (t *wheelTimer) getSlot(nano_timestamp int64) uint32 {
	_diff_tm := nano_timestamp - t.m_cur_time.UnixNano()
	_future_slot := ((_diff_tm / int64(TickerTime)) + int64(t.m_cursor)) % int64(t.getSlotSize())

	return uint32(_future_slot)
}
func (t *wheelTimer) cancelTimer(t_ uint32) {
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

func (t *wheelTimer) setTime(t_ time.Time) {
	t.m_cur_time = t_
}

func (t *wheelTimer) insertSlot(slot_ uint32, t_ *Timer) {
	t.m_map[t_.m_uid] = t_
	_root := t.m_slots[slot_]
	if _root.m_next != nil {
		_root.m_next.m_perv = t_
	}
	t_.m_next = _root.m_next
	_root.m_next = t_
}

func (t *wheelTimer) timeCall(t_ *Timer) {
	if t_.M_call_time < t.m_cur_time.UnixNano() {
		t.insertSlot(t.m_cursor+1, t_)
		return
	}
	t.insertSlot(t.getSlot(t_.M_call_time), t_)
}
func (t *wheelTimer) getSlotSize() uint32 {
	return uint32(len(t.m_slots))
}
func (t *wheelTimer) getNextCursor() uint32 {
	t.m_cursor++
	t.m_cursor %= t.getSlotSize()
	return t.m_cursor
}

func (t *wheelTimer) close() {
	g_close <- struct{}{}
}

func (t *wheelTimer) getAllCall() *ChanTimer {
	_timer_list := make([]*Timer, 0)

	_first_timer := t.m_slots[t.getNextCursor()].m_next
	_now_nano_time := t.m_cur_time.UnixNano()
	for _first_timer != nil {
		if _first_timer.M_call_time > _now_nano_time {
			_first_timer = _first_timer.m_next
			continue
		}
		_append_timer := _first_timer

		_timer_list = append(_timer_list, _append_timer)
		delete(t.m_map, _append_timer.m_uid)
		if _first_timer.m_perv != nil {
			_first_timer.m_perv.m_next = _first_timer.m_next
		}
		if _first_timer.m_next != nil {
			_first_timer.m_next.m_perv = _first_timer.m_perv
		}
		_append_timer.m_next = nil
		_append_timer.m_perv = nil
		_first_timer = _first_timer.m_next
	}
	t.m_slots[t.getNextCursor()] = _first_timer
	return &ChanTimer{
		M_timer_list: _timer_list,
		M_tick_time:  t.m_cur_time,
	}
}
func (t *wheelTimer) loop(_timer_list *ChanTimer) {
	for _, _it := range _timer_list.M_timer_list {
		_it.M_callback(_timer_list.M_tick_time)
		if _it.M_times == -1 {
			_it.M_call_time += _it.M_interval
			TimerCall(_it)
			continue
		} else if _it.M_times > 0 {
			_it.M_times--
			_it.M_call_time += _it.M_interval
			TimerCall(_it)
			continue
		}
	}
}
