package main

import (
	"YMsg"
	"YNet"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/goregular"
	"golang.org/x/image/font/opentype"
	"image/color"
	"log"
)

const (
	ScreenWidth  = 1280
	ScreenHeight = 720
	gridSize     = 10
)

var uiFont font.Face

type Map struct {
	m_user_list map[uint64]YMsg.PositionXY
}

func NewMap() *Map {
	return &Map{
		m_user_list: make(map[uint64]YMsg.PositionXY),
	}
}

var g_map = NewMap()
var g_main_uid uint64
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

	YNet.Register(YMsg.S2C_MESSAGE_MOVE, m.UserMove)
	YNet.Register(YMsg.MSG_S2C_MAP_FULL_SYNC, func(msg_ YMsg.S2CMapFullSync, _ YNet.Session) {
		for _, _it := range msg_.M_user {
			m.AddNewUser(_it)
		}
	})
	YNet.Register(YMsg.MSG_S2C_MAP_ADD_USER, func(msg_ YMsg.S2CMapAddUser, _ YNet.Session) {
		for _,_it := range msg_.M_user{
			m.AddNewUser(_it)
		}

	})
	YNet.Register(YMsg.MSG_S2C_MAP_UPDATE_USER, func(msg_ YMsg.S2CMapUpdateUser, _ YNet.Session) {
		for _,_it := range msg_.M_user{
			m.UpdateUser(_it)
		}
	})
	YNet.Register(YMsg.MSG_S2C_MAP_DELETE_USER, func(msg_ YMsg.S2CMapDeleteUser, _ YNet.Session) {
		for _,_it := range msg_.M_user{
			m.DeleteUser(_it.M_uid)
		}
	})
	YNet.Register(YMsg.MSG_S2C_USER_SUCCESS_LOGIN, func(msg_ YMsg.S2CUserSuccessLogin, _ YNet.Session) {
		
		g_main_uid = msg_.M_uid
	})
	
}
func (m *Map) DeleteUser(uid_ uint64) {
	delete(m.m_user_list, uid_)
}

func (m *Map) AddNewUser(user_data_ YMsg.UserData) {
	m.m_user_list[user_data_.M_uid] = user_data_.M_pos
}
var g_slope string
func (m *Map) UpdateUser(user_data_ YMsg.UserData) {

	if g_main_uid== user_data_.M_uid {
		g_slope = fmt.Sprintf("%.2f", (user_data_.M_pos.M_y -m.m_user_list[user_data_.M_uid].M_y)/(user_data_.M_pos.M_x -m.m_user_list[user_data_.M_uid].M_x))
	}
	m.m_user_list[user_data_.M_uid] = user_data_.M_pos
}

func (m *Map) UserMove(msg_ YMsg.S2C_MOVE, _ YNet.Session) {
	m.m_user_list[msg_.M_uid] = msg_.M_pos
}

func (m *Map) Update() {

}

func (m *Map) Draw(screen *ebiten.Image) {
	text.Draw(screen, g_slope, uiFont, int(100), int(100), color.White)
	for _uid_it, it := range m.m_user_list {
		if m.m_user_list[_uid_it].Distance(it) > 100 {
			panic("1")
		}
		if g_main_uid== _uid_it{
			detailStr := fmt.Sprintf("%.2f,%.2f", it.M_x,it.M_y)
			text.Draw(screen, detailStr, uiFont, int(it.M_x), int(it.M_y + 50), color.White)
			ebitenutil.DrawRect(screen, it.M_x, it.M_y, gridSize, gridSize, color.RGBA{0xff, 0xa0, 0x00, 0xff})
		}else{
			ebitenutil.DrawRect(screen, it.M_x, it.M_y, gridSize, gridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
		}
	}
	/*
		detailStr := fmt.Sprintf("%d", 10)
		text.Draw(screen, detailStr, uiFont, 100, 100, color.White)*/

}
