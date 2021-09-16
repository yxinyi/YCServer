package YPathFinding

import (
	ylog "github.com/yxinyi/YCServer/engine/YLog"
	"github.com/yxinyi/YCServer/engine/YTool"
)

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
	m_queue *YTool.SyncQueue

	m_cache_path map[uint64][]int

	m_call_back_idx uint32
	m_call_back     map[uint32]aStarCallback
}

func NewAStarManager() *AStarManager {
	return &AStarManager{
		m_queue:      YTool.NewSyncQueue(),
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

func (mgr *AStarManager) IsBlock(index_ int) bool {
	_row := index_ / len(mgr.m_maze[0])
	_col := index_ % len(mgr.m_maze[0])
	return mgr.m_maze[_row][_col] != 0
}

func (mgr *AStarManager) Search(st_, ed_ int, cb_ aStarCallback) {
	ylog.Info("[%v:%v]", st_, ed_)
/*	_cache_path, exists := mgr.m_cache_path[uint64(st_)<<32|uint64(ed_)]
	if exists {
		cb_(_cache_path)
		return append(_final_path, before_path_[_loop_idx])
	}*/

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
/*		ylog.Info("###### [%v:%v]", _msg.m_st, _msg.m_ed)
		mgr.m_cache_path[uint64(_msg.m_st)<<32|uint64(_msg.m_ed)] = _msg.m_search_path
		for _idx, _path_node_it := range _msg.m_search_path {
			mgr.m_cache_path[uint64(_path_node_it)<<32|uint64(_msg.m_ed)] = _msg.m_search_path[_idx:]
		}
*/
		delete(mgr.m_call_back, _msg.m_uid)
	}
}
