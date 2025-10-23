package dodge

import (
	"bytes"
	_ "embed"
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

const (
	ScreenWidth  = 900
	ScreenHeight = 600
	playerW      = 40
	playerH      = 60
	borderW      = 24
	midLineW     = 12
	playerSpeed  = 4
	BalleRadius  = 45
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
	frezze                                     bool
	frezzeCooldown                             float32
	balleTake                                  bool
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
	SpeedBalle       float32
	Win              int
	dir              int
	Points1          int
	Points2          int
	background       *ebiten.Image
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
func (g *Game) RestartGame() {
	g.Points1 = 0
	g.Points2 = 0
	g.Win = 0
	g.SpeedBalle = 0
	g.balleX = 451
	g.balleY = 310
	g.dir = 0
	g.ResetPositions()
}
func (p *Player) Update(g *Game) {
	if !p.Dead && !p.frezze {
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
			if p.playerSpeed == 16 && !p.balleTake && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 && !p.balleTake {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
			//set la direction de la balle
			if p.balleTake && ebiten.IsKeyPressed(p.DashKey) {
				g.dir = 1
			}
		}
		if ebiten.IsKeyPressed(p.RightKey) {
			p.X += float64(p.playerSpeed)
			if p.playerSpeed == 16 && !p.balleTake && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 && !p.balleTake {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
			//set la direction de la balle
			if p.balleTake && ebiten.IsKeyPressed(p.DashKey) {
				g.dir = 2
			}
		}
		if ebiten.IsKeyPressed(p.UpKey) {
			p.Y -= float64(p.playerSpeed)
			if p.playerSpeed == 16 && !p.balleTake && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 && !p.balleTake {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
			//set la direction de la balle
			if p.balleTake && ebiten.IsKeyPressed(p.DashKey) {
				g.dir = 3
			}
		}
		if ebiten.IsKeyPressed(p.DownKey) {
			p.Y += float64(p.playerSpeed)
			if p.playerSpeed == 16 && !p.balleTake && math.Abs(float64(p.posXD)-p.X) >= 100 || math.Abs(float64(p.posYD)-p.Y) >= 100 && p.playerSpeed == 16 && !p.balleTake {
				p.cooldown = 3.0
				p.posXD = 0
				p.posYD = 0
			}
			//set la direction de la balle
			if p.balleTake && ebiten.IsKeyPressed(p.DashKey) {
				g.dir = 4
			}
		}
		//dash/lancée
		if ebiten.IsKeyPressed(p.DashKey) {
			if p.cooldown <= 0 && p.playerSpeed == 4 {
				if !p.balleTake {
					p.playerSpeed = 16
					p.cooldown = 0
					p.posXD = int(p.X)
					p.posYD = int(p.Y)
				}
			}
			//petit déplacement de la balle qu'elle est lancée
			if p.balleTake && p.cooldown <= 0 {
				switch g.dir {
				case 1:
					g.balleX -= 100
				case 2:
					g.balleX += 100
				case 3:
					g.balleY -= 100
				case 4:
					g.balleY += 150
				}
				p.balleTake = false
				p.cooldown = 3
				g.SpeedBalle = 2.3
			}
		}
		//déplacement de la balle
		if g.dir == 1 && g.SpeedBalle > 0 && !p.balleTake {
			g.balleX -= g.SpeedBalle
		}
		if g.dir == 2 && g.SpeedBalle > 0 && !p.balleTake {
			g.balleX += g.SpeedBalle
		}
		if g.dir == 3 && g.SpeedBalle > 0 && !p.balleTake {
			g.balleY -= g.SpeedBalle
		}
		if g.dir == 4 && g.SpeedBalle > 0 && !p.balleTake {
			g.balleY += g.SpeedBalle
		}
		//collision balle et joueur
		if g.SpeedBalle <= 0 && !g.p1a.balleTake && !g.p1b.balleTake && !g.p2a.balleTake && !g.p2b.balleTake && CircleRectCollide(float64(g.balleX), float64(g.balleY), BalleRadius, p.X, p.Y, playerW, playerH) {
			p.balleTake = true
		}
		if g.SpeedBalle > 0 && !g.p1a.balleTake && !g.p1b.balleTake && !g.p2a.balleTake && !g.p2b.balleTake && CircleRectCollide(float64(g.balleX), float64(g.balleY), BalleRadius, p.X, p.Y, playerW, playerH) {
			p.Dead = true
			p.deadCooldown = 5.0
			if g.p2a.Dead || g.p2b.Dead {
				g.Points1++
			}
			if g.p1a.Dead || g.p1b.Dead {
				g.Points2++
			}
		}
		//si joueur a la balle alors la balle suit le joueur
		if p.balleTake {
			g.balleX = float32(p.X)
			g.balleY = float32(p.Y)
		}
		//diminution de la vitesse de la balle
		if g.SpeedBalle > 0 && !p.balleTake {
			g.SpeedBalle -= 0.01
		}
		if g.Points1 == 10 {
			g.Win = 1
		}
		if g.Points2 == 10 {
			g.Win = 2
		}
		//rebonds verticaux
		if g.balleY <= 75 {
			g.dir = 4
		}
		if g.balleY >= 625 {
			g.dir = 3
		}
		//rebonds horizontaux
		if g.balleX <= 75 {
			g.dir = 2
		}
		if g.balleX >= 825 {
			g.dir = 1
		}
		// Limites horizontales
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
		p.frezzeCooldown -= 1.0 / 60.0
		if p.frezzeCooldown <= 0 {
			p.frezze = false
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
	g.Win = 0
	g.SpeedBalle = 0
	g.leftBorderColor = color.RGBA{R: 70, G: 130, B: 180, A: 255}  // steelblue
	g.rightBorderColor = color.RGBA{R: 180, G: 70, B: 130, A: 255} // rose
	g.midLineColor = color.RGBA{R: 30, G: 30, B: 30, A: 255}

	// valeurs initiales
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
	bg := ebiten.NewImage(ScreenWidth, ScreenHeight)
	for i := 0; i < 40; i++ {
		x := float32(rand.Intn(ScreenWidth))
		h := float32(rand.Intn(100) + 120)
		c := color.RGBA{R: uint8(20 + rand.Intn(30)), G: uint8(90 + rand.Intn(80)), B: uint8(20), A: 255}
		vector.StrokeLine(bg, x, ScreenHeight-10, x, ScreenHeight-10-h, 3, c, true)
		for j := 0; j < 3; j++ {
			bx := x + float32(rand.Intn(40)-20)
			by := ScreenHeight - 10 - h/3*float32(j)
			vector.StrokeLine(bg, x, by, bx, by-float32(rand.Intn(25)), 2, c, true)
		}
	}
	for i := 0; i < 25; i++ {
		x := float32(rand.Intn(ScreenWidth))
		h := float32(rand.Intn(80) + 40)
		c := color.RGBA{R: uint8(40 + rand.Intn(30)), G: uint8(130 + rand.Intn(100)), B: uint8(40), A: 255}
		vector.StrokeLine(bg, x, ScreenHeight-5, x, ScreenHeight-5-h, 2, c, true)
	}
	vector.StrokeLine(bg, 0, ScreenHeight-5, ScreenWidth, ScreenHeight-5, 10, color.RGBA{0, 90, 0, 255}, true)
	for y := 0; y < ScreenHeight; y++ {
		a := uint8(60 - y/8)
		if a > 0 {
			vector.DrawFilledRect(bg, 0, float32(y), float32(ScreenWidth), 1, color.RGBA{0, 60, 0, a}, true)
		}
	}
	g.background = bg
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
func (g *Game) ResetPositions() {
	// joueurs équipe 1
	g.p1a.X, g.p1a.Y = 100, float64(ScreenHeight-100)
	g.p1b.X, g.p1b.Y = 100, float64(ScreenHeight-200)

	// joueurs équipe 2
	g.p2a.X, g.p2a.Y = float64(ScreenWidth-140), float64(ScreenHeight-100)
	g.p2b.X, g.p2b.Y = float64(ScreenWidth-140), float64(ScreenHeight-200)

	// aucun joueur ne possède la balle
	g.p1a.balleTake = false
	g.p1b.balleTake = false
	g.p2a.balleTake = false
	g.p2b.balleTake = false
}
func (g *Game) Update() error {
	if g.Win == 0 {
		//player dash code
		g.p1a.Update(g)
		g.p1b.Update(g)
		g.p2a.Update(g)
		g.p2b.Update(g)
		if RectRectCollide(g.p1a.X, g.p1a.Y, playerW, playerH, g.p2a.X, g.p2a.Y, playerW, playerH) {
			if g.p1a.playerSpeed == 16 && g.p2a.cooldown <= 0 {
				g.p2a.frezze = true
				g.p2a.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p2a.frezzeCooldown = 5.0
			}
			if g.p2a.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.frezze = true
				g.SpeedBalle = 2
				g.dir = 2
				g.p1a.balleTake = false
				g.p1a.frezzeCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1a.X, g.p1a.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1a.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.frezze = true
				g.p2b.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p2b.frezzeCooldown = 5.0
			}
			if g.p2b.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.frezze = true
				g.p1a.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p1a.frezzeCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2a.X, g.p2a.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2a.cooldown <= 0 {
				g.p2a.frezze = true
				g.p2a.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p2a.frezzeCooldown = 5.0
			}
			if g.p2a.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.frezze = true
				g.p1b.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p1b.frezzeCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.frezze = true
				g.p2b.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p2b.frezzeCooldown = 5.0
			}
			if g.p2b.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.frezze = true
				g.p1b.balleTake = false
				g.SpeedBalle = 2
				g.dir = 2
				g.p1b.frezzeCooldown = 5.0
			}
		}

	}
	if g.Win != 0 {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			bx, by, bw, bh := 350, 350, 200, 60
			if mx >= bx && mx <= bx+bw && my >= by && my <= by+bh {
				g.RestartGame()
			}
		}
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.Win == 0 {
		screen.DrawImage(g.background, nil)
		vector.DrawFilledCircle(screen, g.balleX, g.balleY, BalleRadius, color.RGBA{0, 0, 0, 255}, true)
		//balle
		vector.DrawFilledCircle(screen, g.balleX, g.balleY, BalleRadius, color.RGBA{0, 0, 0, 255}, true)
		// dessiner les joueurs
		g.p1a.Draw(screen)
		g.p1b.Draw(screen)
		g.p2a.Draw(screen)
		g.p2b.Draw(screen)
		// dessiner les points
		// Affichage du score du joueur 1
		switch g.Points1 {
		case 0, 1, 2, 3, 4, 5, 6, 7, 8, 9:
			op := &text.DrawOptions{}
			op.GeoM.Translate(10, 100)
			op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
			text.Draw(screen, fmt.Sprintf("%d", g.Points1), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   50,
			}, op)
		}
		// Affichage du score du joueur 2
		switch g.Points2 {
		case 0, 1, 2:
			op := &text.DrawOptions{}
			op.GeoM.Translate(850, 100)
			op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
			text.Draw(screen, fmt.Sprintf("%d", g.Points2), &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   50,
			}, op)
		}
	}
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
	if g.Win != 0 {
		x, y, w, h := 350, 350, 200, 60
		vector.DrawFilledRect(screen, float32(x), float32(y), float32(w), float32(h), color.RGBA{100, 100, 255, 255}, false)
		op := &text.DrawOptions{}
		op.GeoM.Translate(float64(x+10), float64(y+20))
		op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
		text.Draw(screen, "RESTART", &text.GoTextFace{
			Source: mplusFaceSource,
			Size:   25,
		}, op)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}

func Dodge() *Game {
	game := NewGame()
	return game
}

// joueur ne peuvent pas sortirent du terrain - done
// smach
// vitesse progressive - done
// dash si touche personne personne: éliminé pendant 5 secondes dash cooldown:3
