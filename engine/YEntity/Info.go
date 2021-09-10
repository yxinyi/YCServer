package YEntity


type Inter interface {
	GetInfo()*Info
}

type Info struct {
	M_uid  uint64
	M_type uint32
}

type BaseInfo struct {
	Info
}

func (u *BaseInfo) GetInfo() *Info {
	return &u.Info
}