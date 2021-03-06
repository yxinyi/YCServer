package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/engine/YTool"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Msg"
	"github.com/yxinyi/YCServer/examples/SeamlessExample/Server/Util"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
	"math"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	userGridSize = 5
)

var g_center_pos = YTool.PositionXY{float64(ScreenWidth / 2), float64(ScreenHeight / 2)}

var uiFont font.Face

type Map struct {
	m_user_list map[uint64]Msg.UserData
}

func NewMap() *Map {
	return &Map{
		m_user_list: make(map[uint64]Msg.UserData),
	}
}

var g_map = NewMap()
var g_main_uid uint64
var g_main_path_node []YTool.PositionXY
var g_main_check_node []YTool.PositionXY

type MapMazeInfo struct {
	M_msg *Msg.S2C_AllSyncMapInfo
	*YTool.Rectangle
	M_grid_size float64
}

func NewMapMazeInfo(msg_ *Msg.S2C_AllSyncMapInfo) *MapMazeInfo {
	_info := &MapMazeInfo{}
	_info.M_msg = msg_
	_info.Rectangle = YTool.NewRectangle()
	_up_down_offset, _left_right_offset := Util.MapOffDiff(0x7FFFFFFF<<32|0x7FFFFFFF, msg_.M_map_uid)
	
	_left_up := &YTool.PositionXY{
		M_x: float64(_left_right_offset) * msg_.M_width,
		M_y: float64(_up_down_offset) * msg_.M_height,
	}
	_right_down := &YTool.PositionXY{
		M_x: _left_up.M_x + msg_.M_width,
		M_y: _left_up.M_y + msg_.M_height,
	}
	_info.Rectangle.InitForLefUPRightDown(_left_up, _right_down)
	_info.M_grid_size = msg_.M_height / float64(len(msg_.M_maze))
	return _info
}

var g_map_maze_info = make(map[uint64]*MapMazeInfo)

func (m *Map) Init() {
	tt, err := opentype.Parse(goregular.TTF)
	if err != nil {
		log.Fatal(err.Error())
	}
	uiFont, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    12,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	ebiten.SetMaxTPS(150)
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_FirstEnterMap) {
		m.UpdateUser(msg_.M_data)
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_AllSyncMapInfo) {
		g_map_maze_info[msg_.M_map_uid] = NewMapMazeInfo(&msg_)
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_Login) {
		g_main_uid = msg_.M_main_uid
		//m.AddNewUser(msg_.M_defalut_value)
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_MapAStarNodeUpdate) {
		g_main_path_node = msg_.M_path
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapAddUser) {
		for _, _it := range msg_.M_user {
			m.AddNewUser(_it)
		}
	})
	
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapUpdateUser) {
		//atomic.AddInt32(&_msg_count, 1)
		for _, _it := range msg_.M_user {
			if _it.M_current_map_id != m.MainMapID() {
				fmt.Printf("b[%v]a[%v]\n", m.MainMapID(), _it.M_current_map_id)
			}
			fmt.Printf("pos[%v]\n", _it.M_pos.DebugString())
			
			m.UpdateUser(_it)
		}
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapDeleteUser, ) {
		for _, _it := range msg_.M_user {
			m.DeleteUser(_it.M_uid)
		}
	})
	
}
func (m *Map) DeleteUser(uid_ uint64) {
	delete(m.m_user_list, uid_)
}

func (m *Map) AddNewUser(user_data_ Msg.UserData) {
	m.m_user_list[user_data_.M_uid] = user_data_
}

var g_slope string

func (m *Map) UpdateUser(user_data_ Msg.UserData) {
	
	if g_main_uid == user_data_.M_uid {
		g_slope = fmt.Sprintf("%.2f", (user_data_.M_pos.M_y-m.m_user_list[user_data_.M_uid].M_pos.M_y)/(user_data_.M_pos.M_x-m.m_user_list[user_data_.M_uid].M_pos.M_x))
	}
	m.m_user_list[user_data_.M_uid] = user_data_
}

func (m *Map) UserMove(msg_ Msg.S2C_MOVE, _ YNet.Session) {
	m.m_user_list[msg_.M_uid] = msg_.M_data
}

func (m *Map) MainMapID() uint64 {
	return m.m_user_list[g_main_uid].M_current_map_id
}

func (m *Map) MainMapInfo() *MapMazeInfo {
	return g_map_maze_info[m.MainMapID()]
}

func (m *Map) MainPos() YTool.PositionXY {
	return m.m_user_list[g_main_uid].M_pos
}

func (m *Map) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_tar_x, _tar_y := ebiten.CursorPosition()
		_x_diff := float64(_tar_x) - g_center_pos.M_x
		_y_diff := float64(_tar_y) - g_center_pos.M_y
		
		_fix_tar_pos := &YTool.PositionXY{
			m.MainPos().M_x + _x_diff,
			m.MainPos().M_y + _y_diff,
		}
		
		_tar_map := uint64(0)
		
		_msg := Msg.C2S_UserMove{
			_tar_map,
			*_fix_tar_pos,
		}
		//fmt.Printf("[%v]", _msg.M_tar_map_uid)
		g_client_cnn.SendJson(_msg)
	}
	
	if inpututil.IsKeyJustPressed(ebiten.KeyT) {
		switch ebiten.MaxTPS() {
		case 30:
			ebiten.SetMaxTPS(60)
		case 60:
			ebiten.SetMaxTPS(90)
		case 90:
			ebiten.SetMaxTPS(120)
		case 120:
			ebiten.SetMaxTPS(150)
		case 150:
			ebiten.SetMaxTPS(30)
		}
		
	}
}

func (m *Map) InViewRange(pos YTool.PositionXY) bool {
	_distance := pos.GetOffset(m.MainPos())
	if math.Abs(_distance.M_x) > ScreenWidth/2-10 || math.Abs(_distance.M_y) > ScreenHeight/2-10 {
		return false
	}
	return true
}

func (m *Map) PosConvert(pos YTool.PositionXY) YTool.PositionXY {
	
	_main_user_pos := m.MainPos()
	
	_x_diff := g_center_pos.M_x - _main_user_pos.M_x
	_y_diff := g_center_pos.M_y - _main_user_pos.M_y
	pos.M_x += _x_diff
	pos.M_y += _y_diff
	return pos
}

func (m *Map) Draw(screen *ebiten.Image) {
	//fmt.Printf("[%v] [%v]\n", m.MainMapID(), m.MainPos().DebugString())
	_is_show := 0
	_show_pos := YTool.PositionXY{}
	
	_round_map := make(map[uint64]struct{})
	_round_map[m.MainMapID()] = struct{}{}
	_round_map[m.MainMapID()-1] = struct{}{}
	_round_map[m.MainMapID()+1] = struct{}{}
	_round_map[m.MainMapID()+(1<<32)] = struct{}{}
	_round_map[m.MainMapID()-(1<<32)] = struct{}{}
	_round_map[m.MainMapID()+(1<<32)+1] = struct{}{}
	_round_map[m.MainMapID()+(1<<32)-1] = struct{}{}
	_round_map[m.MainMapID()-(1<<32)+1] = struct{}{}
	_round_map[m.MainMapID()-(1<<32)-1] = struct{}{}
	
	for _map_uid_it := range _round_map {
		_map_it := g_map_maze_info[_map_uid_it]
		if _map_it == nil{
			continue
		}
		for _row_idx_it, _row_it := range _map_it.M_msg.M_maze {
			_row_idx := _row_idx_it
			for _col_idx_it, _block_val := range _row_it {
				_col_idx := _col_idx_it
				if _block_val != 0 {
					_block_pos := YTool.PositionXY{float64(_col_idx)*_map_it.M_grid_size + _map_it.LeftUp.M_x, float64(_row_idx)*_map_it.M_grid_size + _map_it.LeftUp.M_y}
					if !m.InViewRange(_block_pos) {
						continue
					}
					_block_pos = m.PosConvert(_block_pos)
					if _map_it.M_msg.M_map_uid == 9223372034707292159 && _is_show == 0 {
						_is_show = 1
						_show_pos = _block_pos
					}
					_rgb := color.RGBA{
						uint8(((_map_it.M_msg.M_map_uid>>32)+100-0x7fffffff)*77) & 0xff,
						uint8(((_map_it.M_msg.M_map_uid)+133-0x7fffffff)*155) & 0xff,
						uint8(((_map_it.M_msg.M_map_uid>>32)+211-0x7fffffff)*211) & 0xff,
						0xff,
					}
					/*					if _row_idx == 0 && _col_idx == 0 {
										fmt.Printf("first block [%v]\n", _block_pos.DebugString())
									}*/
					ebitenutil.DrawRect(screen, _block_pos.M_x, _block_pos.M_y, _map_it.M_grid_size, _map_it.M_grid_size, _rgb)
				}
			}
		}
	}
	gridSize := float64(10)
	for _, it := range m.m_user_list {
		for _, path_it := range it.M_path {
			if !m.InViewRange(path_it) {
				continue
			}
			_path_pos := m.PosConvert(path_it)
			ebitenutil.DrawRect(screen, _path_pos.M_x, _path_pos.M_y, gridSize, gridSize, color.RGBA{0xff, 0x00, 0x00, 0xff})
		}
	}
	
	for _, path_it := range g_main_path_node {
		if !m.InViewRange(path_it) {
			continue
		}
		_path_pos := m.PosConvert(path_it)
		ebitenutil.DrawRect(screen, _path_pos.M_x, _path_pos.M_y, gridSize, gridSize, color.RGBA{0xff, 0x00, 0x00, 0xff})
	}
	
	for _uid_it, it := range m.m_user_list {
		if !m.InViewRange(it.M_pos) {
			continue
		}
		/*		if m.m_user_list[_uid_it].M_pos.Distance(it.M_pos) > 100 {
				panic("1")
			}*/
		
		if g_main_uid == _uid_it {
			//detailStr := fmt.Sprintf("%.2f,%.2f", it.M_pos.M_x, it.M_pos.M_y)
			//text.Draw(screen, detailStr, uiFont, int(it.M_pos.M_x), int(it.M_pos.M_y+20), color.White)
			_main_user := m.PosConvert(YTool.PositionXY{it.M_pos.M_x + (gridSize-userGridSize)/2, it.M_pos.M_y + (gridSize-userGridSize)/2})
			ebitenutil.DrawRect(screen, _main_user.M_x, _main_user.M_y, userGridSize, userGridSize, color.RGBA{0xff, 0xa0, 0x00, 0xff})
		} else {
			//detailStr := fmt.Sprintf("%.2f,%.2f", it.M_pos.M_x, it.M_pos.M_y)
			//text.Draw(screen, detailStr, uiFont, int(it.M_pos.M_x), int(it.M_pos.M_y+20), color.White)
			_main_user := m.PosConvert(YTool.PositionXY{it.M_pos.M_x + (gridSize-userGridSize)/2, it.M_pos.M_y + (gridSize-userGridSize)/2})
			ebitenutil.DrawRect(screen, _main_user.M_x, _main_user.M_y, userGridSize, userGridSize, color.RGBA{0x80, 0xa0, 0x00, 0xff})
			//ebitenutil.DrawRect(screen, it.M_pos.M_x+(gridSize-userGridSize)/2, it.M_pos.M_y+(gridSize-userGridSize)/2, userGridSize, userGridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
		}
	}
	
	ebitenutil.DebugPrint(screen, fmt.Sprintf("MAX: %d\nTPS: %0.2f\nFPS: %0.2f \nCM[%v]\nPOS[%v] \nFP[%v]", ebiten.MaxTPS(), ebiten.CurrentTPS(), ebiten.CurrentFPS(), m.MainMapID(), m.MainPos().DebugString(), _show_pos.DebugString()))
	
	{
		detailStr := fmt.Sprintf("%d,%d", 100, 100)
		text.Draw(screen, detailStr, uiFont, int(100), int(100), color.White)
	}
	{
		detailStr := fmt.Sprintf("%d,%d", 400, 100)
		text.Draw(screen, detailStr, uiFont, int(400), int(100), color.White)
	}
	{
		detailStr := fmt.Sprintf("%d,%d", 100, 400)
		text.Draw(screen, detailStr, uiFont, int(100), int(400), color.White)
	}
	
}
