package main

import (
	"YNet"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	g_map.Init()
}

// Game represents a game state.
type Game struct {
}

// NewGame generates a new Game object.
func NewMainGame() (*Game, error) {
	g := &Game{
	}
	var err error
	if err != nil {
		return nil, err
	}
	return g, nil
}

// Layout implements ebiten.Game's Layout.
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

// Update updates the current game state.
func (g *Game) Update() error {

	for _net_msg := range YNet.G_net_msg_chan{
		YNet.Dispatch(_net_msg.M_session, _net_msg.M_net_msg)
		if len(YNet.G_net_msg_chan) == 0 {
			break
		}
	}

	g_map.Update()

	return nil
}

// Draw draws the current game to the given screen.
func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.NRGBA{0x00, 0x40, 0x80, 0xff})
	g_map.Draw(screen)
}
