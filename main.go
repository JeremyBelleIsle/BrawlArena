package main

import (
	"BrawlArena/pkg/dodge"
	"BrawlArena/pkg/hockey"
	"BrawlArena/pkg/volley"
	"bytes"
	"image/color"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Game struct {
	screenW       int
	screenH       int
	currentscreen string
	volley        *volley.Game
	hockey        *hockey.Game
	dodge         *dodge.Game
}

var (
	mplusFaceSource *text.GoTextFaceSource
)

func init() {
	rand.Seed(time.Now().UnixNano())
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.PressStart2P_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}
func within(x, y, rx, ry, rw, rh int) bool {
	return x >= rx && x <= rx+rw && y >= ry && y <= ry+rh
}
func (g *Game) Update() error {
	if g.currentscreen == "Menu" && ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		if within(mx, my, 90, 360, 120, 120) {
			g.currentscreen = "HockeyBrawl"
			ebiten.SetWindowSize(hockey.ScreenWidth, hockey.ScreenHeight)

			// Centre la fenêtre sur l'écran
			w, h := ebiten.Monitor().Size()
			x := (w - hockey.ScreenWidth) / 2
			y := (h - hockey.ScreenHeight) / 2
			ebiten.SetWindowPosition(x, y)

			g.hockey = hockey.Hockey()
		}
		if within(mx, my, 270, 360, 120, 120) {
			g.currentscreen = "VolleyBrawl"
			ebiten.SetWindowSize(volley.ScreenWidth, volley.ScreenHeight)

			// Centre la fenêtre sur l'écran
			w, h := ebiten.Monitor().Size()
			x := (w - volley.ScreenWidth) / 2
			y := (h - volley.ScreenHeight) / 2
			ebiten.SetWindowPosition(x, y)

			g.volley = volley.Volley()
		}
		if within(mx, my, 450, 360, 120, 120) {
			g.currentscreen = "DodgeBrawl"
			ebiten.SetWindowSize(dodge.ScreenWidth, dodge.ScreenHeight)

			// Centre la fenêtre sur l'écran
			w, h := ebiten.Monitor().Size()
			x := (w - dodge.ScreenWidth) / 2
			y := (h - dodge.ScreenHeight) / 2
			ebiten.SetWindowPosition(x, y)

			g.dodge = dodge.Dodge()
		}
	}

	switch g.currentscreen {
	case "HockeyBrawl":
		return g.hockey.Update()
	case "VolleyBrawl":
		return g.volley.Update()
	case "DodgeBrawl":
		return g.dodge.Update()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.currentscreen == "Menu" {
		ebitenutil.DrawRect(screen, 90, 360, 120, 120, color.RGBA{0, 255, 0, 255})
		ebitenutil.DrawRect(screen, 270, 360, 120, 120, color.RGBA{0, 255, 0, 255})
		ebitenutil.DrawRect(screen, 450, 360, 120, 120, color.RGBA{0, 255, 0, 255})
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(30), float64(100))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "Brawl Arena", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   55,
		}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(90), float64(330))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "HockeyBrawl", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, op)
		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(270), float64(330))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "VolleyBrawl", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, op)

		op = &text.DrawOptions{}
		op.GeoM.Translate(float64(450), float64(330))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "DodgeBrawl", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   14,
		}, op)
	}
	switch g.currentscreen {
	case "DodgeBrawl":
		g.dodge.Draw(screen)
	case "VolleyBrawl":
		g.volley.Draw(screen)
	case "HockeyBrawl":
		g.hockey.Draw(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	switch g.currentscreen {
	case "DodgeBrawl":
		return g.dodge.Layout(outsideWidth, outsideHeight)
	case "VolleyBrawl":
		return g.volley.Layout(outsideWidth, outsideHeight)
	case "HockeyBrawl":
		return g.hockey.Layout(outsideWidth, outsideHeight)
	}

	return g.screenW, g.screenH
}

func main() {
	game := &Game{
		screenW:       640,
		screenH:       480,
		currentscreen: "Menu",
	}

	ebiten.SetWindowSize(game.screenW, game.screenH)
	ebiten.SetWindowTitle("Brawl Arena")

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
