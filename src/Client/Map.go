package main

import (
	"YMsg"
	"YNet"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
	m_user_list map[uint32]YMsg.PositionXY
}

func NewMap() *Map {
	return &Map{
		m_user_list: make(map[uint32]YMsg.PositionXY),
	}
}

var g_map = NewMap()

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
		m.AddNewUser(msg_.M_user)
	})
	YNet.Register(YMsg.MSG_S2C_MAP_UPDATE_USER, func(msg_ YMsg.S2CMapUpdateUser, _ YNet.Session) {
		m.UpdateUser(msg_.M_user)
	})
	YNet.Register(YMsg.MSG_S2C_MAP_DELETE_USER, func(msg_ YMsg.S2CMapDeleteUser, _ YNet.Session) {
		m.DeleteUser(msg_.M_user.M_uid)
	})

}
func (m *Map) DeleteUser(uid_ uint32) {
	delete(m.m_user_list, uid_)
}

func (m *Map) AddNewUser(user_data_ YMsg.UserData) {
	m.m_user_list[user_data_.M_uid] = user_data_.M_pos
}
func (m *Map) UpdateUser(user_data_ YMsg.UserData) {
	m.m_user_list[user_data_.M_uid] = user_data_.M_pos
}

func (m *Map) UserMove(msg_ YMsg.S2C_MOVE, _ YNet.Session) {
	m.m_user_list[msg_.M_uid] = msg_.M_pos
}

func (m *Map) Update() {

}

func (m *Map) Draw(screen *ebiten.Image) {
	for _, it := range m.m_user_list {
		ebitenutil.DrawRect(screen, it.M_x, it.M_y, gridSize, gridSize, color.RGBA{0x80, 0xa0, 0xc0, 0xff})
	}
	/*
		detailStr := fmt.Sprintf("%d", 10)
		text.Draw(screen, detailStr, uiFont, 100, 100, color.White)*/

}
