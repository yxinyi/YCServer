package PathFinding

import "queue"

type aStarCallbackMsg struct {
	m_uid         uint32
	m_st          int
	m_ed          int
	m_search_path []int
}

func NewAStarCallbackMsg() *aStarCallbackMsg {
	return &aStarCallbackMsg{}
}

func (a *aStarCallbackMsg) Init(uid_ uint32, st_, ed_ int, search_path_ []int) {
	a.m_uid = uid_
	a.m_st = st_
	a.m_ed = ed_
	a.m_search_path = search_path_
}

type aStarCallback func([]int)

type AStarManager struct {
	m_maze  [][]float64
	m_queue *queue.SyncQueue
	
	m_cache_path map[uint64][]int
	
	m_call_back_idx uint32
	m_call_back     map[uint32]aStarCallback
}

func NewAStarManager() *AStarManager {
	return &AStarManager{
		m_queue:      queue.NewSyncQueue(),
		m_call_back:  make(map[uint32]aStarCallback),
		m_cache_path: make(map[uint64][]int),
	}
}

func (mgr *AStarManager) GetMaze() [][]float64 {
	return mgr.m_maze
}

func (mgr *AStarManager) Init(maze_ [][]float64) {
	mgr.m_maze = maze_
}
func (mgr *AStarManager) Search(st_, ed_ int, cb_ aStarCallback) {
	
	_cache_path, exists := mgr.m_cache_path[uint64(st_)<<32|uint64(ed_)]
	if exists {
		cb_(_cache_path)
		return
	}
	
	_tmp_idx := mgr.m_call_back_idx
	mgr.m_call_back_idx++
	mgr.m_call_back[_tmp_idx] = cb_
	go func() {
		_a := NewAStar()
		_a.Init(mgr.m_maze)
		_ret := _a.SearchBetterWithIndex(st_, ed_)
		_msg := NewAStarCallbackMsg()
		_msg.Init(_tmp_idx, st_, ed_, _ret)
		mgr.m_queue.Add(_msg)
	}()
}

func (mgr *AStarManager) Update() {
	for {
		if mgr.m_queue.Len() == 0 {
			break
		}
		_msg := mgr.m_queue.Pop().(*aStarCallbackMsg)
		mgr.m_call_back[_msg.m_uid](_msg.m_search_path)
		
		mgr.m_cache_path[uint64(_msg.m_st)<<32|uint64(_msg.m_ed)] = _msg.m_search_path
		delete(mgr.m_call_back, _msg.m_uid)
	}
}
