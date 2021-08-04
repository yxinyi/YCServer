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

// *秒后执行
func AfterSecondsCall(after_time_ uint32, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = int64(after_time_) + time.Now().Unix()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// *秒后执行,每隔*秒执行*次
func AfterSecondsWithIntervalAndLoopTimesCall(after_time_ uint32, inter_val_ int64, loop_times int, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = loop_times
	_t.M_interval = inter_val_ * time.Second.Nanoseconds()
	_t.M_call_time = int64(after_time_)*time.Second.Nanoseconds() + time.Now().Unix()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// *秒后执行,每隔*秒执行
func AfterSecondsWithIntervalCall(after_time_ uint32, inter_val_ int64, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_interval = inter_val_ * time.Second.Nanoseconds()
	_t.M_call_time = int64(after_time_)*time.Second.Nanoseconds() + time.Now().Unix()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行
func WhenTimeCall(timestamp_ int64, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_call_time = timestamp_ * time.Second.Nanoseconds()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行,每隔*秒执行
func WhenTimeWithIntervalCall(timestamp_ int64, inter_val_ int64, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = -1
	_t.M_interval = inter_val_ * time.Second.Nanoseconds()
	_t.M_call_time = timestamp_ * time.Second.Nanoseconds()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 在 * 时执行,每隔*秒执行 *次
func WhenTimeWithIntervalAndLoopTimesCall(timestamp_ int64, inter_val_ int64, loop_times int, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_times = loop_times
	_t.M_interval = inter_val_ * time.Second.Nanoseconds()
	_t.M_call_time = timestamp_ * time.Second.Nanoseconds()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 每天 * 点执行
func EveryDayCall(timestamp_ int64, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = timestamp_ * time.Second.Nanoseconds()
	g_add_timer_channel <- _t
	return _t.m_uid
}

// 每天 * 点到 *点 间隔*秒执行
func TimeToTimeWithIntervalCall(start_timestamp_ int64, end_timestamp_ int64, inter_val_ int64, cb_ TimerCallBack) uint32 {
	_t := newTimer()
	_t.M_callback = cb_
	_t.M_call_time = start_timestamp_ * time.Second.Nanoseconds()
	_t.M_times = int(int64(end_timestamp_-start_timestamp_) / inter_val_)
	_t.M_interval = inter_val_
	g_add_timer_channel <- _t
	return _t.m_uid
}
func Loop(_timer_list *ChanTimer) {
	g_timer_manager.loop(_timer_list)
}

func GetTimerSize() uint32 {
	return g_timer_manager.getTimerSize()
}