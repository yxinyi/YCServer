package aoi

import (
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
)

type GoNGAoiMoveCallBack func(notify_ uint64, action_ map[uint64]struct{})
type GoNGAoiEnterCallBack func(notify_ uint64, action_ map[uint64]struct{})
type GoNGAoiQuitCallBack func(notify_ uint64, action_ map[uint64]struct{})

type GoNGAoiAction struct {
	m_action     uint32
	m_notify_obj uint64
	m_action_obj map[uint64]struct{}
}

type GoNGInAoiAction struct {
	m_action uint32
	m_obj    GoAoiObj
}
type GoNineGirdAoiManager struct {
	M_height         float64
	M_width          float64
	m_aoi_list       map[uint32]*GoNineGirdAoiCell
	M_current_index  map[uint64]uint32
	m_block_height   float64
	m_block_width    float64
	m_block_size     float64
	m_obj_copy       map[uint64]*GoAoiObj
	M_quit_callback  GoNGAoiQuitCallBack
	M_move_callback  GoNGAoiMoveCallBack
	M_enter_callback GoNGAoiEnterCallBack
	M_action_out_    *YTool.SyncQueue
	m_action_in      chan GoNGInAoiAction
}

func NewGoNineGirdAoiCellManager(width_, height_, block_size_ float64) *GoNineGirdAoiManager {
	_mgr := &GoNineGirdAoiManager{
		m_aoi_list:      make(map[uint32]*GoNineGirdAoiCell),
		M_current_index: make(map[uint64]uint32),
		M_action_out_:   YTool.NewSyncQueue(),
		m_action_in:     make(chan GoNGInAoiAction, 100),
		m_obj_copy:      make(map[uint64]*GoAoiObj),
	}
	_mgr.M_height = height_
	_mgr.M_width = width_
	_mgr.m_block_height = height_ / block_size_
	_mgr.m_block_width = width_ / block_size_
	_mgr.m_block_size = block_size_
	return _mgr
}

func (mgr *GoNineGirdAoiManager) Init(move_call_ GoNGAoiMoveCallBack, enter_call_ GoNGAoiEnterCallBack, quit_call_ GoNGAoiQuitCallBack) {
	for _row_idx := uint32(0); _row_idx < uint32(mgr.m_block_size); _row_idx++ {
		for _col_idx := uint32(0); _col_idx < uint32(mgr.m_block_size); _col_idx++ {
			_cell := NewGoNineGirdAoiCell()
			mgr.m_aoi_list[mgr.buildIndex(_row_idx, _col_idx)] = _cell
		}
	}
	mgr.M_move_callback = move_call_
	mgr.M_enter_callback = enter_call_
	mgr.M_quit_callback = quit_call_
	go func() {
		for {
			select {
			case _action_obj := <-mgr.m_action_in:
				switch _action_obj.m_action {
				case GO_AOI_ACTION_ENTER:
					mgr.enter(_action_obj.m_obj)
				case GO_AOI_ACTION_UPDATE:
					mgr.move(_action_obj.m_obj)
				case GO_AOI_ACTION_QUIT:
					mgr.quit(_action_obj.m_obj)
				}
			}
		}
	}()
}

func (mgr *GoNineGirdAoiManager) getDiff(lhs_ map[uint32]struct{}, rhs_ map[uint32]struct{}) map[uint32]struct{} {
	_ret := make(map[uint32]struct{})
	for _it := range lhs_ {
		_ret[_it] = struct{}{}
	}
	for _it := range rhs_ {
		delete(_ret, _it)
	}
	
	return _ret
}

func (mgr *GoNineGirdAoiManager) FindObj(uid_ uint64) *GoAoiObj {
	return mgr.m_obj_copy[uid_]
}

func (obj *GoAoiObj) InViewRange(rhs_ *GoAoiObj) bool {
	return obj.PositionXY.Distance(rhs_.PositionXY) < obj.M_view_range
}

func (mgr *GoNineGirdAoiManager) getObjNotInViewRangeMap(obj_ *GoAoiObj, cell_index_ uint32) map[uint64]map[uint64]struct{} {
	_enter_sync_list := make(map[uint64]map[uint64]struct{})
	_cell := mgr.m_aoi_list[cell_index_]
	if _cell == nil {
		return _enter_sync_list
	}
	for _it := range _cell.GetWatch() {
		obj := mgr.FindObj(_it)
		if obj != nil {
			if !obj_.InViewRange(obj) {
				_, exists := _enter_sync_list[obj_.M_uid]
				if !exists {
					_enter_sync_list[obj_.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[obj_.M_uid][obj.M_uid] = struct{}{}
			}
			if !obj.InViewRange(obj_) {
				_, exists := _enter_sync_list[obj.M_uid]
				if !exists {
					_enter_sync_list[obj.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[obj.M_uid][obj_.M_uid] = struct{}{}
				
			}
			
		}
	}
	return _enter_sync_list
}

func (mgr *GoNineGirdAoiManager) getObjMap(obj_ *GoAoiObj, cell_index_ uint32) map[uint64]map[uint64]struct{} {
	_enter_sync_list := make(map[uint64]map[uint64]struct{})
	_cell := mgr.m_aoi_list[cell_index_]
	if _cell == nil {
		return _enter_sync_list
	}
	for _it := range _cell.GetWatch() {
		_tmp_obj := mgr.FindObj(_it)
		if _tmp_obj != nil {
			{
				_, exists := _enter_sync_list[obj_.M_uid]
				if !exists {
					_enter_sync_list[obj_.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[obj_.M_uid][_tmp_obj.M_uid] = struct{}{}
			}
			{
				_, exists := _enter_sync_list[_tmp_obj.M_uid]
				if !exists {
					_enter_sync_list[_tmp_obj.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[_tmp_obj.M_uid][obj_.M_uid] = struct{}{}
				
			}
			
		}
	}
	return _enter_sync_list
}

func (mgr *GoNineGirdAoiManager) getObjInViewRangeMap(obj_ *GoAoiObj, cell_index_ uint32) map[uint64]map[uint64]struct{} {
	_enter_sync_list := make(map[uint64]map[uint64]struct{})
	_cell := mgr.m_aoi_list[cell_index_]
	if _cell == nil {
		return _enter_sync_list
	}
	for _it := range _cell.GetWatch() {

		_tmp_obj := mgr.FindObj(_it)
		if _tmp_obj != nil {
			if obj_.InViewRange(_tmp_obj) {
				_, exists := _enter_sync_list[obj_.M_uid]
				if !exists {
					_enter_sync_list[obj_.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[obj_.M_uid][_tmp_obj.M_uid] = struct{}{}
			}
			if _tmp_obj.InViewRange(obj_) {
				_, exists := _enter_sync_list[_tmp_obj.M_uid]
				if !exists {
					_enter_sync_list[_tmp_obj.M_uid] = make(map[uint64]struct{})
				}
				_enter_sync_list[_tmp_obj.M_uid][obj_.M_uid] = struct{}{}
				
			}
			
		}
	}
	return _enter_sync_list
}

func (mgr *GoNineGirdAoiManager) Update() {
	for {
		if mgr.M_action_out_.Len() == 0 {
			break
		}
		_act := mgr.M_action_out_.Pop().(GoNGAoiAction)
		switch _act.m_action {
		case GO_AOI_ACTION_ENTER:
			mgr.M_enter_callback(_act.m_notify_obj, _act.m_action_obj)
		case GO_AOI_ACTION_UPDATE:
			mgr.M_move_callback(_act.m_notify_obj, _act.m_action_obj)
		case GO_AOI_ACTION_QUIT:
			mgr.M_quit_callback(_act.m_notify_obj, _act.m_action_obj)
		}
	}
	
}
func (mgr *GoNineGirdAoiManager) sendOutUpdateAction(map_ map[uint64]map[uint64]struct{}) {
	for _key, _set_it := range map_ {
		mgr.M_action_out_.Add(GoNGAoiAction{
			GO_AOI_ACTION_UPDATE,
			_key,
			_set_it,
		})
	}
}

func (mgr *GoNineGirdAoiManager) updateCell(enter_ *GoAoiObj, cell_set_ map[uint32]struct{}) {
	{
		_enter_map_sync := make(map[uint64]map[uint64]struct{})
		for _it := range cell_set_ {
			_enter_map_sync = YTool.Uint64MapUint64SetMerge(_enter_map_sync, mgr.getObjInViewRangeMap(enter_, _it))
		}
		mgr.sendOutUpdateAction(_enter_map_sync)
	}
	{
		_quit_map_sycn := make(map[uint64]map[uint64]struct{})
		for _it := range cell_set_ {
			_quit_map_sycn = YTool.Uint64MapUint64SetMerge(_quit_map_sycn, mgr.getObjNotInViewRangeMap(enter_, _it))
		}
		mgr.sendOutQuitAction(_quit_map_sycn)
	}
	
}

func (mgr *GoNineGirdAoiManager) sendOutEnterAction(map_ map[uint64]map[uint64]struct{}) {
	for _key, _set_it := range map_ {
		mgr.M_action_out_.Add(GoNGAoiAction{
			GO_AOI_ACTION_ENTER,
			_key,
			_set_it,
		})
	}
}
func (mgr *GoNineGirdAoiManager) enterCell(enter_ *GoAoiObj, cell_set_ map[uint32]struct{}) {
	_enter_map_sync := make(map[uint64]map[uint64]struct{})
	for _it := range cell_set_ {
		_enter_map_sync = YTool.Uint64MapUint64SetMerge(_enter_map_sync, mgr.getObjInViewRangeMap(enter_, _it))
	}
	mgr.sendOutEnterAction(_enter_map_sync)
}

func (mgr *GoNineGirdAoiManager) Enter(enter_ GoAoiObj) {
	mgr.m_action_in <- GoNGInAoiAction{
		GO_AOI_ACTION_ENTER,
		enter_,
	}
}
func (mgr *GoNineGirdAoiManager) enter(enter_ GoAoiObj, ) {
	_current_index := mgr.CalcIndex(enter_.PositionXY)
	_round_arr := mgr.getRoundBlock(_current_index)
	mgr.enterCell(&enter_, _round_arr)
	mgr.m_aoi_list[_current_index].Watch(enter_.M_uid)
	mgr.m_obj_copy[enter_.M_uid] = &enter_
	mgr.M_current_index[enter_.M_uid] = _current_index
}
func (mgr *GoNineGirdAoiManager) sendOutQuitAction(map_ map[uint64]map[uint64]struct{}) {
	for _key, _set_it := range map_ {
		mgr.M_action_out_.Add(GoNGAoiAction{
			GO_AOI_ACTION_QUIT,
			_key,
			_set_it,
		})
	}
}
func (mgr *GoNineGirdAoiManager) quitCell(enter_ *GoAoiObj, cell_set_ map[uint32]struct{}) {
	_quit_map_sync := make(map[uint64]map[uint64]struct{})
	for _it := range cell_set_ {
		_quit_map_sync = YTool.Uint64MapUint64SetMerge(_quit_map_sync, mgr.getObjMap(enter_, _it))
	}
	mgr.sendOutQuitAction(_quit_map_sync)
}

func (mgr *GoNineGirdAoiManager) Quit(quit_ GoAoiObj) {
	mgr.m_action_in <- GoNGInAoiAction{
		GO_AOI_ACTION_QUIT,
		quit_,
	}
}
func (mgr *GoNineGirdAoiManager) quit(quit_ GoAoiObj) {
	_current_index := mgr.CalcIndex(quit_.PositionXY)
	
	_round_arr := mgr.getRoundBlock(_current_index)
	mgr.quitCell(&quit_, _round_arr)
	delete(mgr.M_current_index, quit_.M_uid)
	mgr.m_aoi_list[_current_index].Forget(quit_.M_uid)
	delete(mgr.m_obj_copy, quit_.M_uid)
}

func (mgr *GoNineGirdAoiManager) ActionUpdate(move_ GoAoiObj) {
	mgr.m_action_in <- GoNGInAoiAction{
		GO_AOI_ACTION_UPDATE,
		move_,
	}
}

func (mgr *GoNineGirdAoiManager) move(move_ GoAoiObj) {
	mgr.m_obj_copy[move_.M_uid] = &move_
	_old_round_arr := mgr.getOldRoundBlock(move_.M_uid)
	
	_current_index := mgr.CalcIndex(move_.PositionXY)
	_new_round_arr := mgr.getRoundBlock(_current_index)
	
	_enter_cell := getDiff(_new_round_arr, _old_round_arr)
	mgr.enterCell(&move_, _enter_cell)
	
	if _current_index != mgr.M_current_index[move_.M_uid] {
		mgr.m_aoi_list[_current_index].Watch(move_.M_uid)
	}
	
	_update_cell := getDiff(_new_round_arr, _enter_cell)
	mgr.updateCell(&move_, _update_cell)
	
	if _current_index != mgr.M_current_index[move_.M_uid] {
		mgr.m_aoi_list[mgr.M_current_index[move_.M_uid]].Forget(move_.M_uid)
	}
	
	_quit_cell := getDiff(_old_round_arr, _new_round_arr)
	mgr.quitCell(&move_, _quit_cell)
	
	mgr.M_current_index[move_.M_uid] = _current_index
	
}

func (mgr *GoNineGirdAoiManager) CalcIndex(xy_ Msg.PositionXY) uint32 {
	return mgr.buildIndex(uint32(xy_.M_x/mgr.m_block_width), uint32(xy_.M_y/mgr.m_block_height))
}

func (mgr *GoNineGirdAoiManager) buildIndex(col_, row_ uint32) uint32 {
	return col_ + row_*uint32(mgr.m_block_size)
}

func (mgr *GoNineGirdAoiManager) getOldRoundBlock(uid_ uint64) map[uint32]struct{} {
	_old_index := mgr.M_current_index[uid_]
	return mgr.getRoundBlock(_old_index)
}

func (mgr *GoNineGirdAoiManager) getRoundBlock(cent_index_ uint32) map[uint32]struct{} {
	_ret_round := make(map[uint32]struct{})
	_cent_idex := int(cent_index_)
	_block_size := int(mgr.m_block_size)
	
	_max_idx := int(mgr.m_block_size * mgr.m_block_size)
	
	_cent_row := int(cent_index_ / uint32(mgr.m_block_size))
	_ret_round[cent_index_] = struct{}{}
	{
		_left_up := _cent_idex - _block_size - 1
		if _left_up >= 0 && (_left_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_left_up)] = struct{}{}
		}
	}
	
	{
		_up := _cent_idex - _block_size
		if _up >= 0 && (_up/_block_size+1) == _cent_row {
			_ret_round[uint32(_up)] = struct{}{}
		}
	}
	{
		_up_right := _cent_idex - _block_size + 1
		if _up_right >= 0 && (_up_right/_block_size+1) == _cent_row {
			_ret_round[uint32(_up_right)] = struct{}{}
		}
	}
	
	{
		_left := _cent_idex - 1
		if _left >= 0 && (_left/_block_size) == _cent_row {
			_ret_round[uint32(_left)] = struct{}{}
		}
	}
	{
		_right := _cent_idex + 1
		if _right >= 0 && (_right/_block_size) == _cent_row {
			_ret_round[uint32(_right)] = struct{}{}
		}
	}
	
	{
		_down_left := _cent_idex + _block_size - 1
		if _down_left < _max_idx && (_down_left/_block_size-1) == _cent_row {
			_ret_round[uint32(_down_left)] = struct{}{}
		}
	}
	
	{
		_down := _cent_idex + _block_size
		if _down < _max_idx && (_down/_block_size-1) == _cent_row {
			_ret_round[uint32(_down)] = struct{}{}
		}
	}
	{
		_down_right := _cent_idex + _block_size + 1
		if _down_right < _max_idx && (_down_right/_block_size-1) == _cent_row {
			_ret_round[uint32(_down_right)] = struct{}{}
		}
	}
	
	return _ret_round
}
