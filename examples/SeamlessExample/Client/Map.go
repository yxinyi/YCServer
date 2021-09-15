package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/yxinyi/YCServer/engine/YNet"
	"github.com/yxinyi/YCServer/examples/AoiAstarExample/Msg"
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
	gridSize     = 10
	userGridSize = 5
)

var g_center_pos = Msg.PositionXY{float64(ScreenWidth / 2), float64(ScreenHeight / 2)}

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
var g_main_path_node []Msg.PositionXY
var g_main_check_node []Msg.PositionXY
var g_map_maze_info Msg.S2C_FirstEnterMap

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
	ebiten.SetMaxTPS(30)
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_FirstEnterMap) {
		g_map_maze_info = msg_
		m.UpdateUser(msg_.M_data)
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_Login) {
		g_main_uid = msg_.M_data.M_uid
		m.AddNewUser(msg_.M_data)
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2C_MapAStarNodeUpdate) {
		g_main_path_node = msg_.M_path
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapAddUser) {
		for _, _it := range msg_.M_user {
			m.AddNewUser(_it)
		}
		
	})
	/*	_msg_count := int32(0)
		go func() {
			_ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-_ticker.C:
					fmt.Printf("msg count [%v]\n", atomic.LoadInt32(&_msg_count))
					atomic.StoreInt32(&_msg_count, 0)
				}
			}
	
		}()*/
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapUpdateUser) {
		//atomic.AddInt32(&_msg_count, 1)
		for _, _it := range msg_.M_user {
			m.UpdateUser(_it)
		}
	})
	YNet.Register(func(_ YNet.Session, msg_ Msg.S2CMapDeleteUser, ) {
		for _, _it := range msg_.M_user {
			m.DeleteUser(_it.M_uid)
		}
	})
	
	/*	YNet.Register(Msg.MsgID_S2CUserMove, m.UserMove)
		YNet.Register(Msg.MSG_S2C_MAP_FULL_SYNC, func(msg_ Msg.S2CMapFullSync, _ YNet.Session) {
			for _, _it := range msg_.M_user {
				m.AddNewUser(_it)
			}
		})
		YNet.Register(Msg.MSG_S2C_MAP_ADD_USER, func(msg_ Msg.S2CMapAddUser, _ YNet.Session) {
			for _, _it := range msg_.M_user {
				m.AddNewUser(_it)
			}
	
		})
		YNet.Register(Msg.MSG_S2C_MAP_UPDATE_USER, func(msg_ Msg.S2CMapUpdateUser, _ YNet.Session) {
			for _, _it := range msg_.M_user {
				m.UpdateUser(_it)
			}
		})
		YNet.Register(Msg.MSG_S2C_MAP_DELETE_USER, func(msg_ Msg.S2CMapDeleteUser, _ YNet.Session) {
			for _, _it := range msg_.M_user {
				m.DeleteUser(_it.M_uid)
			}
		})
		YNet.Register(Msg.MsgID_S2CUserSuccessLogin, func(msg_ Msg.S2CUserSuccessLogin, _ YNet.Session) {
			g_main_uid = msg_.M_data.M_uid
			m.AddNewUser(msg_.M_data)
		})
		YNet.Register(Msg.MSG_S2C_MAP_ASTAR_NODE_UPDATE, func(msg_ Msg.S2C_MapAStarNodeUpdate, _ YNet.Session) {
			g_main_path_node = msg_.M_path
		})
		YNet.Register(Msg.MSG_S2C_MAP_FLUSH_MAP_MAZE, func(msg_ Msg.S2CFlushMapMaze, _ YNet.Session) {
			g_map_maze_info = msg_
		})*/
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

func (m *Map) MainPos() Msg.PositionXY {
	return m.m_user_list[g_main_uid].M_pos
}

func (m *Map) Update() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		_tar_x, _tar_y := ebiten.CursorPosition()
		_x_diff := float64(_tar_x) - g_center_pos.M_x
		_y_diff := float64(_tar_y) - g_center_pos.M_y
		
		//_touch_pos := m.PosConvert(Msg.PositionXY{ float64(_tar_x), float64(_tar_y)})
		g_client_cnn.SendJson(Msg.C2S_UserMove{
			Msg.PositionXY{
				m.MainPos().M_x + _x_diff,
				m.MainPos().M_y + _y_diff,
			},
		})
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

func (m *Map) InViewRange(pos Msg.PositionXY) bool {
	_distance := pos.DistancePosition(m.MainPos())
	if math.Abs( _distance.M_x) > ScreenWidth/2-10 || math.Abs(_distance.M_y) > ScreenHeight/2-10{
		return false
	}
	return true
}
func (m *Map) PosConvert(pos Msg.PositionXY) Msg.PositionXY {
	
	_main_user_pos := m.MainPos()
	
	_x_diff := g_center_pos.M_x - _main_user_pos.M_x
	_y_diff := g_center_pos.M_y - _main_user_pos.M_y
	pos.M_x += _x_diff
	pos.M_y += _y_diff
	return pos
}

func (m *Map) Draw(screen *ebiten.Image) {
	//text.Draw(screen, g_slope, uiFont, int(100), int(100), color.White)
	
	_grid_size := g_map_maze_info.M_height / float64(len(g_map_maze_info.M_maze))
	for _row_idx_it, _row_it := range g_map_maze_info.M_maze {
		_row_idx := _row_idx_it
		for _col_idx_it, _block_val := range _row_it {
			_col_idx := _col_idx_it
			if _block_val != 0 {
				_block_pos :=Msg.PositionXY{float64(_col_idx) * _grid_size, float64(_row_idx) * _grid_size}
				if !m.InViewRange(_block_pos){
					continue
				}
				_block_pos = m.PosConvert(_block_pos)
				ebitenutil.DrawRect(screen, _block_pos.M_x, _block_pos.M_y, _grid_size, _grid_size, color.Black)
			} /* else {
				ebitenutil.DrawRect(screen, float64(_col_idx)*_grid_size, float64(_row_idx)*_grid_size, _grid_size, _grid_size, color.White)
			}*/
			/*detailStr := fmt.Sprintf("%d", _row_idx*len(g_map_maze_info.M_maze[0])+_col_idx)
			text.Draw(screen, detailStr, uiFont, int (float64(_col_idx)*_grid_size + _grid_size/2), int (float64(_row_idx)*_grid_size+ _grid_size/2), color.White) */
		}
	}
	
	for _, it := range m.m_user_list {
		for _, path_it := range it.M_path {
			if !m.InViewRange(path_it){
				continue
			}
			_path_pos := m.PosConvert(path_it)
			ebitenutil.DrawRect(screen, _path_pos.M_x, _path_pos.M_y, gridSize, gridSize, color.RGBA{0xff, 0x00, 0x00, 0xff})
		}
	}
	
	for _, path_it := range g_main_path_node {
		if !m.InViewRange(path_it){
			continue
		}
		_path_pos := m.PosConvert(path_it)
		ebitenutil.DrawRect(screen, _path_pos.M_x, _path_pos.M_y, gridSize, gridSize, color.RGBA{0xff, 0x00, 0x00, 0xff})
	}
	
	for _uid_it, it := range m.m_user_list {
		if !m.InViewRange(it.M_pos){
			continue
		}
		if m.m_user_list[_uid_it].M_pos.Distance(it.M_pos) > 100 {
			panic("1")
		}
		
		if g_main_uid == _uid_it {
			//detailStr := fmt.Sprintf("%.2f,%.2f", it.M_pos.M_x, it.M_pos.M_y)
			//text.Draw(screen, detailStr, uiFont, int(it.M_pos.M_x), int(it.M_pos.M_y+20), color.White)
			_main_user := m.PosConvert(Msg.PositionXY{it.M_pos.M_x + (gridSize-userGridSize)/2, it.M_pos.M_y + (gridSize-userGridSize)/2})
			ebitenutil.DrawRect(screen, _main_user.M_x, _main_user.M_y, userGridSize, userGridSize, color.RGBA{0xff, 0xa0, 0x00, 0xff})
		} else {
			//detailStr := fmt.Sprintf("%.2f,%.2f", it.M_pos.M_x, it.M_pos.M_y)
			//text.Draw(screen, detailStr, uiFont, int(it.M_pos.M_x), int(it.M_pos.M_y+20), color.White)
			_main_user := m.PosConvert(Msg.PositionXY{it.M_pos.M_x + (gridSize-userGridSize)/2, it.M_pos.M_y + (gridSize-userGridSize)/2})
			ebitenutil.DrawRect(screen, _main_user.M_x, _main_user.M_y, userGridSize, userGridSize, color.RGBA{0x80, 0xa0, 0x00, 0xff})
			//ebitenutil.DrawRect(screen, it.M_pos.M_x+(gridSize-userGridSize)/2, it.M_pos.M_y+(gridSize-userGridSize)/2, userGridSize, userGridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
		}
	}
	
	ebitenutil.DebugPrint(screen, fmt.Sprintf("MAX: %d\nTPS: %0.2f\nFPS: %0.2f", ebiten.MaxTPS(), ebiten.CurrentTPS(), ebiten.CurrentFPS()))
}
