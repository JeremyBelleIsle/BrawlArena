package volley

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

var hue float64

const (
	ScreenWidth  = 900
	ScreenHeight = 600
	playerW      = 40
	playerH      = 60
	borderW      = 24
	midLineW     = 12
	playerSpeed  = 4
)

// Player représente un joueur simple
type Player struct {
	X, Y                                       float64
	W, H                                       float64
	Color                                      color.RGBA
	LeftKey, RightKey, UpKey, DownKey, DashKey ebiten.Key
	MinX, MaxX                                 float64
	cooldown                                   float64
	playerSpeed                                int
	posXD                                      int
	posYD                                      int
	Dead                                       bool
	deadCooldown                               float32
	DeadCNT                                    int
	KillCNT                                    int
}

// Game contient l'état
type Game struct {
	leftBorderColor  color.RGBA
	rightBorderColor color.RGBA
	midLineColor     color.RGBA
	p1a, p1b         *Player
	p2a, p2b         *Player
	balleX           float32
	balleY           float32
	targetX          float32
	targetY          float32
	Next_targetX     float32
	Next_targetY     float32
	SpeedBalle       float32
	Win              int
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
func (p *Player) Update() {
	if !p.Dead {
		if p.cooldown > 0 {
			p.cooldown -= 1.0 / 60.0
			if p.cooldown < 0 {
				p.cooldown = 0
			}
		}
		if p.cooldown > 0 {
			p.playerSpeed = 4
		}
		if ebiten.IsKeyPressed(p.LeftKey) {
			p.X -= float64(p.playerSpeed)
			if p.playerSpeed == 16 && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
		}
		if ebiten.IsKeyPressed(p.RightKey) {
			p.X += float64(p.playerSpeed)
			if p.playerSpeed == 16 && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
		}
		if ebiten.IsKeyPressed(p.UpKey) {
			p.Y -= float64(p.playerSpeed)
			if p.playerSpeed == 16 && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
		}
		if ebiten.IsKeyPressed(p.DownKey) {
			p.Y += float64(p.playerSpeed)
			if p.playerSpeed == 16 && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
		}
		if ebiten.IsKeyPressed(p.DashKey) {
			if p.cooldown <= 0 && p.playerSpeed == 4 {
				p.playerSpeed = 16
				p.cooldown = 0
				p.posXD = int(p.X)
				p.posYD = int(p.Y)
			}
		}
		// Limites horizontales
		if p.X < p.MinX {
			p.X = p.MinX
		}
		if p.X+40 > ScreenWidth {
			p.X = ScreenWidth - 40
		}
		// Limites verticales
		if p.Y < 0 {
			p.Y = 0
		}
		if p.Y > float64(ScreenHeight)-p.H {
			p.Y = float64(ScreenHeight) - p.H
		}
	} else {
		p.deadCooldown -= 1.0 / 60.0
		if p.deadCooldown <= 0 {
			p.Dead = false
		}
	}
}

func (p *Player) Draw(screen *ebiten.Image) {
	if !p.Dead {
		vector.DrawFilledRect(screen, float32(p.X), float32(p.Y), float32(p.W), float32(p.H), p.Color, true)
	}
}

func NewGame() *Game {
	g := &Game{}
	g.balleX = 451
	g.balleY = 310
	g.targetX = 100
	g.targetY = 300
	g.Win = 0
	g.SpeedBalle = 3
	g.leftBorderColor = color.RGBA{R: 70, G: 130, B: 180, A: 255}  // steelblue
	g.rightBorderColor = color.RGBA{R: 180, G: 70, B: 130, A: 255} // rose
	g.midLineColor = color.RGBA{R: 30, G: 30, B: 30, A: 255}

	// positions initiales
	p1a := &Player{
		X:           float64(100),
		Y:           float64(ScreenHeight - 100),
		W:           playerW,
		H:           playerH,
		Color:       color.RGBA{R: 255, G: 200, B: 0, A: 255},
		LeftKey:     ebiten.KeyA,
		RightKey:    ebiten.KeyD,
		UpKey:       ebiten.KeyW,
		DownKey:     ebiten.KeyS,
		DashKey:     ebiten.KeyE,
		MinX:        borderW,
		MaxX:        ScreenWidth/2 - midLineW/2,
		playerSpeed: playerSpeed,
	}
	p1b := &Player{
		X:           float64(100),
		Y:           float64(ScreenHeight - 200),
		W:           playerW,
		H:           playerH,
		Color:       color.RGBA{R: 255, G: 220, B: 80, A: 255},
		LeftKey:     ebiten.KeyF,
		RightKey:    ebiten.KeyH,
		UpKey:       ebiten.KeyT,
		DownKey:     ebiten.KeyG,
		DashKey:     ebiten.KeyY,
		MinX:        borderW,
		MaxX:        ScreenWidth/2 - midLineW/2,
		playerSpeed: playerSpeed,
	}
	p2a := &Player{
		X:           float64(ScreenWidth - 140),
		Y:           float64(ScreenHeight - 100),
		W:           playerW,
		H:           playerH,
		Color:       color.RGBA{R: 0, G: 200, B: 255, A: 255},
		LeftKey:     ebiten.KeyArrowLeft,
		RightKey:    ebiten.KeyArrowRight,
		UpKey:       ebiten.KeyArrowUp,
		DownKey:     ebiten.KeyArrowDown,
		DashKey:     ebiten.KeyShiftRight,
		MinX:        ScreenWidth/2 + midLineW/2,
		MaxX:        ScreenWidth - borderW,
		playerSpeed: playerSpeed,
	}
	p2b := &Player{
		X:           float64(ScreenWidth - 140),
		Y:           float64(ScreenHeight - 200),
		W:           playerW,
		H:           playerH,
		Color:       color.RGBA{R: 80, G: 220, B: 255, A: 255},
		LeftKey:     ebiten.KeyJ,
		RightKey:    ebiten.KeyL,
		UpKey:       ebiten.KeyI,
		DownKey:     ebiten.KeyK,
		DashKey:     ebiten.KeyO,
		MinX:        ScreenWidth/2 + midLineW/2,
		MaxX:        ScreenWidth - borderW,
		playerSpeed: playerSpeed,
	}
	g.p1a = p1a
	g.p1b = p1b
	g.p2a = p2a
	g.p2b = p2b
	return g
}
func CircleRectCollide(cx, cy, r, rx, ry, rw, rh float64) bool {
	closestX := math.Max(rx, math.Min(cx, rx+rw))
	closestY := math.Max(ry, math.Min(cy, ry+rh))
	dx := cx - closestX
	dy := cy - closestY
	return (dx*dx + dy*dy) <= (r * r)
}
func RectRectCollide(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 && y1 < y2+h2 && y1+h1 > y2
}
func (g *Game) Update() error {
	if g.Win == 0 {
		hue += 0.02
		if hue > 2*math.Pi {
			hue = 0
		}

		g.SpeedBalle += 0.0032
		g.p1a.Update()
		g.p1b.Update()
		g.p2a.Update()
		g.p2b.Update()
		if RectRectCollide(g.p1a.X, g.p1a.Y, playerW, playerH, g.p2a.X, g.p2a.Y, playerW, playerH) {
			if g.p1a.playerSpeed == 16 && g.p2a.cooldown <= 0 {
				g.p2a.Dead = true
				g.p2a.deadCooldown = 5.0
				g.p2a.DeadCNT++
				g.p1a.KillCNT++
			}
			if g.p2a.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.Dead = true
				g.p1a.deadCooldown = 5.0
				g.p1a.DeadCNT++
				g.p2a.KillCNT++
			}
		}
		if RectRectCollide(g.p1a.X, g.p1a.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1a.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.Dead = true
				g.p2b.deadCooldown = 5.0
				g.p2b.DeadCNT++
				g.p1a.KillCNT++
			}
			if g.p2b.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.Dead = true
				g.p1a.deadCooldown = 5.0
				g.p1a.DeadCNT++
				g.p2b.KillCNT++
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2a.X, g.p2a.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2a.cooldown <= 0 {
				g.p2a.Dead = true
				g.p2a.deadCooldown = 5.0
				g.p2a.DeadCNT++
			}
			if g.p2a.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.Dead = true
				g.p1b.deadCooldown = 5.0
				g.p1b.DeadCNT++
				g.p2a.KillCNT++
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.Dead = true
				g.p2b.deadCooldown = 5.0
				g.p2b.DeadCNT++
				g.p1b.KillCNT++
			}
			if g.p2b.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.Dead = true
				g.p1b.deadCooldown = 5.0
				g.p1b.DeadCNT++
				g.p2b.KillCNT++
			}
		}
		// vecteur direction
		dx := g.targetX - g.balleX
		dy := g.targetY - g.balleY
		dist := float32(math.Hypot(float64(dx), float64(dy)))

		if dist > 30 {
			// normalisation + déplacement
			g.balleX += (dx / dist) * g.SpeedBalle
			g.balleY += (dy / dist) * g.SpeedBalle
		} else {
			if math.Abs(float64(g.targetX-g.Next_targetX)) > 500 {
				g.targetX = g.Next_targetX
				g.targetY = g.Next_targetY
			} else if CircleRectCollide(float64(g.targetX), float64(g.targetY), 30, g.p1a.X, g.p1a.Y, playerW, playerH) || CircleRectCollide(float64(g.targetX), float64(g.targetY), 30, g.p1b.X, g.p1b.Y, playerW, playerH) || CircleRectCollide(float64(g.targetX), float64(g.targetY), 30, g.p2a.X, g.p2a.Y, playerW, playerH) || CircleRectCollide(float64(g.targetX), float64(g.targetY), 30, g.p2b.X, g.p2b.Y, playerW, playerH) {
				for math.Abs(float64(g.targetX-g.Next_targetX)) < 500 {
					g.Next_targetX = float32(rand.Intn(ScreenWidth-200) + 100)
					g.Next_targetY = float32(rand.Intn(ScreenHeight-200) + 100)
					if g.Next_targetX > 300 || g.Next_targetX < 600 {
						for g.Next_targetX > 300 && g.Next_targetX < 600 {
							g.Next_targetX = float32(rand.Intn(ScreenWidth-200) + 100)
							g.Next_targetY = float32(rand.Intn(ScreenHeight-200) + 100)
						}
					}
				}
			} else {
				if g.balleX < 450 {
					g.Win = 2
				}
				if g.balleX > 450 {
					g.Win = 1
				}
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Win == 0 {
		cr := uint8((math.Sin(hue) + 1) * 127)
		cg := uint8((math.Sin(hue+2*math.Pi/3) + 1) * 127)
		cb := uint8((math.Sin(hue+4*math.Pi/3) + 1) * 127)
		dynamicColor := color.RGBA{cr, cg, cb, 255}
		// fond uni
		screen.Fill(dynamicColor)

		// bords de terrain (gauche et droite)
		vector.DrawFilledRect(screen, 0, 0, borderW, float32(ScreenHeight), g.leftBorderColor, true)
		vector.DrawFilledRect(screen, float32(ScreenWidth-borderW), 0, borderW, float32(ScreenHeight), g.rightBorderColor, true)

		// grosse ligne au milieu
		vector.DrawFilledRect(screen, float32(ScreenWidth/2-midLineW/2), 0, midLineW, float32(ScreenHeight), g.midLineColor, true)

		// dessiner les joueurs
		g.p1a.Draw(screen)
		g.p1b.Draw(screen)
		g.p2a.Draw(screen)
		g.p2b.Draw(screen)

		//dessiner la balle
		centerX := float32(ScreenWidth) / 2
		centerY := float32(ScreenHeight) / 2
		distToCenter := float32(math.Hypot(float64(g.balleX-centerX), float64(g.balleY-centerY)))
		maxDist := float32(math.Hypot(float64(centerX), float64(centerY)))
		minRadius := float32(35)
		maxRadius := float32(77)
		balleRadius := minRadius + (maxRadius-minRadius)*(1.0-distToCenter/maxDist)
		vector.DrawFilledCircle(screen, g.balleX, g.balleY, balleRadius, color.White, true)
		//prochain point
		vector.DrawFilledCircle(screen, g.targetX, g.targetY, 30, color.RGBA{0, 255, 0, 255}, true)
	} else {
		if g.Win > 0 {
			// Affichage du titre du gagnant
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(900/4), float64(200))
			op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
			text.Draw(screen, fmt.Sprintf("Team %d WIN", g.Win), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   53,
			}, op)

			// Stats équipe 1
			vector.DrawFilledRect(screen, 150, 300, 100, 100, g.p1a.Color, true)
			vector.DrawFilledRect(screen, 300, 300, 100, 100, g.p1b.Color, true)
			op = &text.DrawOptions{}
			op.GeoM.Translate(150, 420)
			text.Draw(screen, fmt.Sprintf("K:%d D:%d", g.p1a.KillCNT, g.p1a.DeadCNT), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, op)
			op.GeoM.Translate(150, 0)
			text.Draw(screen, fmt.Sprintf("K:%d D:%d", g.p1b.KillCNT, g.p1b.DeadCNT), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, op)

			// Stats équipe 2
			vector.DrawFilledRect(screen, 500, 300, 100, 100, g.p2a.Color, true)
			vector.DrawFilledRect(screen, 650, 300, 100, 100, g.p2b.Color, true)
			op = &text.DrawOptions{}
			op.GeoM.Translate(500, 420)
			text.Draw(screen, fmt.Sprintf("K:%d D:%d", g.p2a.KillCNT, g.p2a.DeadCNT), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, op)
			op.GeoM.Translate(150, 0)
			text.Draw(screen, fmt.Sprintf("K:%d D:%d", g.p2b.KillCNT, g.p2b.DeadCNT), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   20,
			}, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func Volley() *Game {
	game := NewGame()
	return game
}

// joueur ne peuvent pas sortirent du terrain - done
// smach
// vitesse progressive - done
// dash si touche personne personne: éliminé pendant 5 secondes dash cooldown:3
