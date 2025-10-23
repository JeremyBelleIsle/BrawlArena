package hockey

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
	goalTimer        time.Time
	showGoal         bool
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
	g.showGoal = false
	g.Win = 0
	g.SpeedBalle = 0
	g.balleX = 451
	g.balleY = 310
	g.dir = 0
	g.ResetPositions()
}
func (p *Player) Update(g *Game) {
	if !p.Dead && !g.showGoal {
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
				g.SpeedBalle = 5
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
		if !g.p1a.balleTake && !g.p1b.balleTake && !g.p2a.balleTake && !g.p2b.balleTake && CircleRectCollide(float64(g.balleX), float64(g.balleY), BalleRadius, p.X, p.Y, playerW, playerH) {
			p.balleTake = true
		}
		//si joueur a la balle alors la balle suit le joueur
		if p.balleTake {
			g.balleX = float32(p.X)
			g.balleY = float32(p.Y)
		}
		//diminution de la vitesse de la balle
		if g.SpeedBalle > 0 && !p.balleTake {
			g.SpeedBalle -= 0.03
		}

		if !g.showGoal { // Empêche de redétecter un but pendant l’affichage
			if g.balleX >= 825 && g.balleY > 250 && g.balleY < 450 {
				g.Points1++
				g.goalTimer = time.Now()
				g.showGoal = true
				g.SpeedBalle = 0
				g.balleX = 451
				g.balleY = 310
				p.balleTake = false
				g.ResetPositions()
			}
			if g.balleX <= 75 && g.balleY > 250 && g.balleY < 450 {
				g.Points2++
				g.goalTimer = time.Now()
				g.showGoal = true
				g.balleX = 451
				g.balleY = 310
				g.SpeedBalle = 0
				p.balleTake = false
				g.ResetPositions()
			}
		}
		if g.Points1 == 3 {
			g.Win = 1
		}
		if g.Points2 == 3 {
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
	}
	if g.showGoal && time.Since(g.goalTimer) > 2*time.Second {
		g.showGoal = false
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
				g.p2a.Dead = true
				g.p2a.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p2a.deadCooldown = 5.0
				g.p2a.DeadCNT++
				g.p1a.KillCNT++
			}
			if g.p2a.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.Dead = true
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p1a.balleTake = false
				g.p1a.deadCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1a.X, g.p1a.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1a.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.Dead = true
				g.p2b.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p2b.deadCooldown = 5.0
			}
			if g.p2b.playerSpeed == 16 && g.p1a.cooldown <= 0 {
				g.p1a.Dead = true
				g.p1a.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p1a.deadCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2a.X, g.p2a.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2a.cooldown <= 0 {
				g.p2a.Dead = true
				g.p2a.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p2a.deadCooldown = 5.0
			}
			if g.p2a.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.Dead = true
				g.p1b.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p1b.deadCooldown = 5.0
			}
		}
		if RectRectCollide(g.p1b.X, g.p1b.Y, playerW, playerH, g.p2b.X, g.p2b.Y, playerW, playerH) {
			if g.p1b.playerSpeed == 16 && g.p2b.cooldown <= 0 {
				g.p2b.Dead = true
				g.p2b.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p2b.deadCooldown = 5.0
			}
			if g.p2b.playerSpeed == 16 && g.p1b.cooldown <= 0 {
				g.p1b.Dead = true
				g.p1b.balleTake = false
				g.SpeedBalle = 3.5
				g.dir = 2
				g.p1b.deadCooldown = 5.0
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
		//dessiner le terrain de hockey(IA)
		screen.Fill(color.RGBA{180, 220, 255, 255})
		margin := float32(60)
		rinkW := float32(ScreenWidth) - 2*margin
		rinkH := float32(ScreenHeight) - 2*margin
		rinkX := float32(margin)
		rinkY := float32(margin)
		cornerR := float32(40)
		vector.DrawFilledRect(screen, rinkX+cornerR, rinkY, rinkW-2*cornerR, rinkH, color.White, true)
		vector.DrawFilledRect(screen, rinkX, rinkY+cornerR, cornerR, rinkH-2*cornerR, color.White, true)
		vector.DrawFilledRect(screen, rinkX+rinkW-cornerR, rinkY+cornerR, cornerR, rinkH-2*cornerR, color.White, true)
		vector.DrawFilledCircle(screen, rinkX+cornerR, rinkY+cornerR, cornerR, color.White, true)
		vector.DrawFilledCircle(screen, rinkX+rinkW-cornerR, rinkY+cornerR, cornerR, color.White, true)
		vector.DrawFilledCircle(screen, rinkX+cornerR, rinkY+rinkH-cornerR, cornerR, color.White, true)
		vector.DrawFilledCircle(screen, rinkX+rinkW-cornerR, rinkY+rinkH-cornerR, cornerR, color.White, true)
		vector.StrokeRect(screen, rinkX, rinkY, rinkW, rinkH, 3, color.RGBA{0, 0, 0, 255}, true)
		centerLineW := float32(6)
		vector.DrawFilledRect(screen, rinkX+rinkW/2-centerLineW/2, rinkY+10, centerLineW, rinkH-20, color.RGBA{220, 20, 60, 255}, true)
		blueLineW := float32(8)
		blueOffset := float32(180)
		vector.DrawFilledRect(screen, rinkX+blueOffset, rinkY+10, blueLineW, rinkH-20, color.RGBA{10, 50, 150, 255}, true)
		vector.DrawFilledRect(screen, rinkX+rinkW-blueOffset-blueLineW, rinkY+10, blueLineW, rinkH-20, color.RGBA{10, 50, 150, 255}, true)
		centerX := rinkX + rinkW/2
		centerY := rinkY + rinkH/2
		vector.StrokeCircle(screen, centerX, centerY, 70, 4, color.RGBA{200, 30, 30, 255}, true)
		vector.DrawFilledCircle(screen, centerX, centerY, 6, color.RGBA{200, 30, 30, 255}, true)
		offsetX := float32(200)
		offsetY := float32(0)
		vector.StrokeCircle(screen, centerX-offsetX, centerY-offsetY, 40, 3, color.RGBA{10, 50, 150, 255}, true)
		vector.StrokeCircle(screen, centerX+offsetX, centerY-offsetY, 40, 3, color.RGBA{10, 50, 150, 255}, true)
		vector.StrokeCircle(screen, centerX-offsetX, centerY+offsetY, 40, 3, color.RGBA{10, 50, 150, 255}, true)
		vector.StrokeCircle(screen, centerX+offsetX, centerY+offsetY, 40, 3, color.RGBA{10, 50, 150, 255}, true)
		goalOffset := float32(30)
		creaseR := float32(60)
		vector.StrokeCircle(screen, rinkX+goalOffset+creaseR, rinkY+rinkH/2, creaseR, 2, color.RGBA{10, 50, 150, 255}, true)
		vector.StrokeCircle(screen, rinkX+rinkW-goalOffset-creaseR, rinkY+rinkH/2, creaseR, 2, color.RGBA{10, 50, 150, 255}, true)
		vector.DrawFilledCircle(screen, rinkX+goalOffset+creaseR+20, rinkY+rinkH/2-80, 5, color.RGBA{200, 30, 30, 255}, true)
		vector.DrawFilledCircle(screen, rinkX+rinkW-goalOffset-creaseR-20, rinkY+rinkH/2+80, 5, color.RGBA{200, 30, 30, 255}, true)
		goalColor := color.RGBA{255, 0, 0, 180}
		vector.DrawFilledRect(screen, 0, 250, 75, 200, goalColor, true)
		vector.StrokeRect(screen, 0, 250, 75, 200, 3, color.Black, true)
		vector.DrawFilledRect(screen, 825, 250, 75, 200, goalColor, true)
		vector.StrokeRect(screen, 825, 250, 75, 200, 3, color.Black, true)
		//balle
		vector.DrawFilledCircle(screen, g.balleX, g.balleY, BalleRadius, color.RGBA{0, 0, 0, 255}, true)
		// dessiner les joueurs
		g.p1a.Draw(screen)
		g.p1b.Draw(screen)
		g.p2a.Draw(screen)
		g.p2b.Draw(screen)
		if g.balleX >= 808 && g.balleY > 250 && g.balleY < 450 || g.balleX <= 58 && g.balleY > 250 && g.balleY < 450 {
			op := &text.DrawOptions{}
			op.GeoM.Translate(float64(900/4), float64(600/2))
			op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 0})
			text.Draw(screen, "GOAL!", &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   175,
			}, op)
		}
		// dessiner les points
		// Affichage du score du joueur 1
		switch g.Points1 {
		case 0, 1, 2:
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
		if g.showGoal {
			op := &text.DrawOptions{}
			op.GeoM.Translate((900/4)-200, (600/2)-75)
			op.ColorScale.ScaleWithColor(color.RGBA{222, 49, 99, 255})
			text.Draw(screen, "GOAL!", &text.GoTextFace{
				Source: mplusFaceSource,
				Size:   175,
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

func Hockey() *Game {
	game := NewGame()
	return game
}

// joueur ne peuvent pas sortirent du terrain - done
// smach
// vitesse progressive - done
// dash si touche personne personne: éliminé pendant 5 secondes dash cooldown:3
