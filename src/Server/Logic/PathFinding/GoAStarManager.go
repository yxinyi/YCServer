package PathFinding

import "queue"

type AStarCallbackMsg struct {
	m_uid         uint32
	m_st          int
	m_ed          int
	m_search_path []int
}
type AStarCallback func([]int)

type AStarManager struct {
	m_maze  [][]float64
	m_queue *queue.SyncQueue
	
	m_cache_path map[uint64][]int
	
	m_call_back_idx uint32
	m_call_back     map[uint32]AStarCallback
}

func NewAStarManager(maze_ [][]float64) *AStarManager {
	return &AStarManager{
		m_maze:       maze_,
		m_queue:      queue.NewSyncQueue(),
		m_call_back:  make(map[uint32]AStarCallback),
		m_cache_path: make(map[uint64][]int),
	}
}

func (mgr *AStarManager) Search(st_, ed_ int, cb_ AStarCallback) {
	
	_cache_path, exists := mgr.m_cache_path[uint64(st_)<<32|uint64(ed_)]
	if exists {
		cb_(_cache_path)
		return
	}
	
	_tmp_idx := mgr.m_call_back_idx
	mgr.m_call_back_idx++
	mgr.m_call_back[_tmp_idx] = cb_
	go func() {
		_msg := AStarCallbackMsg{}
		_msg.m_uid = _tmp_idx
		_msg.m_st = st_
		_msg.m_ed = ed_
		_a := NewAStar()
		_a.Init(mgr.m_maze)
		_ret := _a.SearchWithIndex(st_, ed_)
		_msg.m_search_path = _ret
		mgr.m_queue.Add(_msg)
	}()
}

func (mgr *AStarManager) Update() {
	for {
		if mgr.m_queue.Len() == 0 {
			break
		}
		_msg := mgr.m_queue.Pop().(AStarCallbackMsg)
		mgr.m_call_back[_msg.m_uid](_msg.m_search_path)
		
		mgr.m_cache_path[uint64(_msg.m_st)<<32|uint64(_msg.m_ed)] = _msg.m_search_path
		delete(mgr.m_call_back, _msg.m_uid)
	}
}
