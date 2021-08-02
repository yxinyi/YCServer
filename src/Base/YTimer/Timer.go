package YTimer

var time_uid_count uint32 = 0

type Timer struct {
	m_perv      *Timer
	m_next      *Timer
	m_uid       uint32
	M_interval  int64
	M_call_time int64
	M_times     int32
	M_callback  TimerCallBack
}

func newTimer() *Timer {
	time_uid_count++
	return &Timer{
		m_uid: time_uid_count,
	}
}
