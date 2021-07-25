module YClient

go 1.16

require (
	YMsg v0.0.0-00010101000000-000000000000
	YNet v0.0.0-00010101000000-000000000000
	github.com/hajimehoshi/ebiten/v2 v2.1.3
	golang.org/x/image v0.0.0-20210220032944-ac19c3e999fb
)

replace YNet => ../Base/YNet

replace YMsg => ../Base/YMsg
