package YTimer

import "time"

func TimerCall(t_ *Timer) uint32 {
	g_add_timer_channel <- t_
	return t_.m_uid
}

func Close() {
	g_close <- struct{}{}
}
func CancelTimer(uid_ uint32) {
	g_timer_manager.cancelTimer(uid_)
}

func convertToTickUnit(t_ YSecond) int64 {
	return int64((float64(t_)) * float64(int64(time.Second.Nanoseconds())/int64(TickerTime)))
}

// *秒后执行
func AfterSecondsCall(after_time_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = convertToTickUnit(after_time_ + YSecond(time.Now().Unix()))
	g_add_timer_channel <- _t
	return _t.m_uid
}

// *秒后执行,每隔*秒执行*次
func AfterSecondsWithIntervalAndLoopTimesCall(after_time_ YSecond, inter_val_ YSecond, loop_times int, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = loop_times
	_t.M_interval = convertToTickUnit(inter_val_)
	_t.M_call_time = convertToTickUnit(after_time_ + YSecond(time.Now().Unix()))
	g_add_timer_channel <- _t
	return _t.m_uid
}

// *秒后执行,每隔*秒执行
func AfterSecondsWithIntervalCall(after_time_ YSecond, inter_val_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_interval = convertToTickUnit(inter_val_)
	_t.M_call_time = convertToTickUnit(after_time_ + YSecond(time.Now().Unix()))
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行
func WhenTimeCall(timestamp_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_call_time = convertToTickUnit(timestamp_)
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行,每隔*秒执行
func WhenTimeWithIntervalCall(timestamp_ YSecond, inter_val_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_interval = convertToTickUnit(inter_val_)
	_t.M_call_time = convertToTickUnit(timestamp_)
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行,每隔*秒执行 *次
func WhenTimeWithIntervalAndLoopTimesCall(timestamp_ YSecond, inter_val_ YSecond, loop_times int, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = loop_times
	_t.M_interval = convertToTickUnit(inter_val_)
	_t.M_call_time = convertToTickUnit(timestamp_)
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 每天 * 点执行
func EveryDayCall(timestamp_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = convertToTickUnit(timestamp_)
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 每天 * 点到 *点 间隔*秒执行
func TimeToTimeWithIntervalCall(start_timestamp_ YSecond, end_timestamp_ YSecond, inter_val_ YSecond, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = convertToTickUnit(start_timestamp_)
	_t.M_times = int((end_timestamp_-start_timestamp_) / inter_val_)
	_t.M_interval = convertToTickUnit(inter_val_)
	g_add_timer_channel <- _t
	return _t.m_uid
}
func Loop(_timer_list *ChanTimer) {
	g_timer_manager.loop(_timer_list)
}

func GetTimerSize() uint32 {
	return g_timer_manager.getTimerSize()
}
