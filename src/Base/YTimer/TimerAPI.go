package YTimer

func TimerCall(t_ *Timer) uint32 {
	return t_.m_uid
}

func CancelTimer(uid_ uint32) {
	g_timer_manager.cancelTimer(uid_)
}

// *秒后执行
func AfterSecondsCall(after_time_ uint32, cb_ TimerCallBack) uint32 {
	_t := newTimer()

	return _t.m_uid
}

// *秒后执行,每隔*秒执行*次
func AfterSecondsWithIntervalAndLoopTimesCall(after_time_ uint32, inter_val_ uint32, loop_times int, cb_ TimerCallBack) uint32 {
	return 0
}

// *秒后执行,每隔*秒执行
func AfterSecondsWithIntervalCall(after_time__ uint32, inter_val_ uint32, cb_ TimerCallBack) uint32 {
	return 0
}

// 在 * 时执行
func WhenTimeCall(timestamp_ int64, cb_ TimerCallBack) uint32 {
	return 0
}

// *秒后执行,每隔*秒执行
func WhenTimeWithIntervalCall(timestamp_ int64, inter_val_ uint32, cb_ TimerCallBack) uint32 {
	return 0
}

// *秒后执行,每隔*秒执行
func WhenTimeWithIntervalAndLoopTimesCall(timestamp_ int64, inter_val_ uint32, loop_times int, cb_ TimerCallBack) uint32 {
	return 0
}

// 每天 * 点执行
func EveryDayCall(timestamp_ int64, cb_ TimerCallBack) uint32 {
	return 0
}

// 每天 * 点到 *点 间隔*秒执行
func TimeToTimeWithIntervalCall(start_timestamp_ int64, end_timestamp_ int64, inter_val_ uint32, cb_ TimerCallBack) uint32 {
	return 0
}
