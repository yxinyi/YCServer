package YTimer

import (
	"fmt"
	"sync"
	"time"
)

type TimerCallBack func()

var g_timer_manager timerManagerInter
var g_add_timer_channel = make(chan *Timer, 100)
var g_cancel_timer = make(chan uint32, 100)

var G_call = make(chan *ChanTimer, 100)

var g_close = make(chan struct{})

const (
	WheelSlotCount               = 3

	TickerTime     time.Duration = time.Millisecond*100
)
type YSecond = float64

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
	getCurTime() *time.Time
	cancelTimer(uint32)
	close()
	getNextCursor() uint32
	getSlotSize() uint32
	getTimerSize() uint32
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
	_wg := &sync.WaitGroup{}
	_wg.Add(1)
	go func() {
		g_timer_manager.setTime(time.Now())
		ticker := time.Tick(TickerTime)
		_wg.Done()
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
	_wg.Wait()
}

func (t *wheelTimer) getSlot(timestamp int64) uint32 {
	_diff_tm := timestamp - convertToTickUnit(float64(t.m_cur_time.Unix()))
	_future_slot := (int64(_diff_tm) + int64(t.m_cursor)) % int64(t.getSlotSize())
	//fmt.Printf("timestamp[%v] t.m_cur_time.Unix() [%v] diff [%v] slot [%v] \n",timestamp,t.m_cur_time.Unix(),_diff_tm,_future_slot)
	fmt.Printf("_future_slot [%d] _diff_tm[%v]\n", _future_slot,_diff_tm)
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
	//fmt.Printf("setTime [%v]\n", t.m_cur_time.Unix())

}

func (t *wheelTimer) insertSlot(slot_ uint32, t_ *Timer) {
	t.m_map[t_.m_uid] = t_
	//fmt.Printf("cur [%v] insertSlot [%v]",t.m_cursor,slot_)
	_root := t.m_slots[slot_]
	if _root.m_next != nil {
		_root.m_next.m_perv = t_
	}
	t_.m_perv = _root
	t_.m_next = _root.m_next
	_root.m_next = t_
	t_.M_slot = slot_
}

func (t *wheelTimer) timeCall(t_ *Timer) {
	if t_.M_call_time <= convertToTickUnit(float64(t.m_cur_time.Unix())) {
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
	//fmt.Printf("t.m_cursor [%v] \n", t.m_cursor)
	return t.m_cursor
}

func (t *wheelTimer) getTimerSize() uint32 {
	return uint32(len(t.m_map))
}

func (t *wheelTimer) close() {
	g_close <- struct{}{}
}

func (t *wheelTimer) getCurTime() *time.Time {
	return &t.m_cur_time
}
func (t *wheelTimer) getAllCall() *ChanTimer {
	_timer_list := make([]*Timer, 0)
	_next_cursor := t.getNextCursor()
	_root := t.m_slots[_next_cursor].m_next
	_now_time := convertToTickUnit(float64(t.m_cur_time.Unix()))
	for _root != nil {
		if _root.M_call_time > _now_time {
			_root = _root.m_next
			continue
		}
		_append_timer := _root

		_timer_list = append(_timer_list, _append_timer)
		delete(t.m_map, _append_timer.m_uid)
		if _root.m_perv != nil {
			_root.m_perv.m_next = _root.m_next
		}
		if _root.m_next != nil {
			_root.m_next.m_perv = _root.m_perv
		}
		_append_timer.m_next = nil
		_append_timer.m_perv = nil
		_root = _root.m_next
	}
	//t.m_slots[_next_cursor].m_next = _root
	return &ChanTimer{
		M_timer_list: _timer_list,
		M_tick_time:  t.m_cur_time,
	}
}
func (t *wheelTimer) loop(_timer_list *ChanTimer) {
	for _, _it := range _timer_list.M_timer_list {
		_it.M_callback()
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
