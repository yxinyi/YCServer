package YMsg

const (
	MESSAGE_TEST uint32 = iota
	C2S_MESSAGE_MOVE
	S2C_MESSAGE_MOVE

	MSG_S2C_MAP_FULL_SYNC
	MSG_S2C_MAP_ADD_USER
	MSG_S2C_MAP_UPDATE_USER
	MSG_S2C_MAP_DELETE_USER
)

type Message struct {
	Id     int
	Number int
}

type UserData struct {
	M_uid uint32
	M_pos   PositionXY
}

type PositionXY struct {
	M_x float64
	M_y float64
}

type C2S_MOVE struct {
	M_uid uint32
	M_pos   PositionXY
}


type S2C_MOVE struct {
	M_uid uint32
	M_pos   PositionXY
}

type S2CMapFullSync struct {
	M_user []UserData
}

type S2CMapAddUser struct {
	M_user UserData
}
type S2CMapUpdateUser struct {
	M_user UserData
}
type S2CMapDeleteUser struct {
	M_user UserData
}

