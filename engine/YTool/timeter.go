package YTool

import (
	"time"
)

type TimerMeterInfo struct {
	Name  string
	Timer int64
}

type TimerMeter struct {
	timerInfoPool map[string]int
	timerPool     map[int]*TimerMeterInfo
}

func NewTimerMeter() *TimerMeter {
	return &TimerMeter{
		timerInfoPool: make(map[string]int),
		timerPool:     make(map[int]*TimerMeterInfo),
	}
}
func (t *TimerMeter) Start(name string) {
	timer := &TimerMeterInfo{
		Name:  name,
		Timer: time.Now().UnixNano(),
	}
	_, exists := t.timerInfoPool[name]
	if exists {
		t.timerPool[t.timerInfoPool[name]] = timer
	} else {
		idx := len(t.timerPool)
		t.timerInfoPool[name] = idx
		t.timerPool[idx] = timer
	}
}

func (t *TimerMeter) End(name string) {
	idx, exists := t.timerInfoPool[name]
	if !exists {
		return
	}
	delete(t.timerInfoPool, name)
	t.timerPool[idx].Timer = time.Now().UnixNano() - t.timerPool[idx].Timer
}

/*func (t *TimerMeter) Print(name string) {
	idx, exists := t.timerInfoPool[name]
	if !exists {
		return
	}
	
	l4g.Info("[TimerMeter] name [%v] time [%04f]", name, t.timerPool[idx].Timer/int64(time.Millisecond))
}

func (t *TimerMeter) PrintAll() {
	meterList := make([]*TimerMeterInfo, 0)
	for idx := 0; idx < len(t.timerPool); idx++ {
		meterList = append(meterList, t.timerPool[idx])
		//l4g.Info("[TimerMeter] name [%v] time [%04f]", t.timerPool[idx].Name, float64(t.timerPool[idx].Timer/int64(time.Millisecond))/float64(1000))
	}
	sort.Slice(meterList, func(lhs int, rhs int) bool {
		return meterList[lhs].Timer > meterList[rhs].Timer
	})
	
	for _, it := range meterList {
		l4g.Info("[TimerMeter] name [%v] time [%04f]", it.Name, float64(it.Timer/int64(time.Millisecond))/float64(1000))
	}
}*/