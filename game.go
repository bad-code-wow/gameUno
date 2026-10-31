/*
you can tell its human cus the comments suck and i need 2 4k monitors to real a full line
*/
package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"math"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var cam rl.Camera2D

type playDat struct {
	pos     rl.Vector2
	vel     rl.Vector2
	prevPos rl.Vector2
}

type line struct {
	Power float32
	Start rl.Vector2
	End   rl.Vector2
}

type bgLine struct {
	length float32
	start  float32
	angle  float32
}

type levelFull struct {
	LevelRect     []rl.Rectangle
	LevelEvilRect []rl.Rectangle
	Level         []line
}

var bomb = rl.Vector2{X: 10000, Y: 10000}

// test

var levelRect = []rl.Rectangle{{X: 90, Y: 290, Width: 190, Height: 190}, {X: 200, Y: 0, Width: 190, Height: 300}}
var levelEvilRect = []rl.Rectangle{{X: 390, Y: 390, Width: 190, Height: 190}}
var frozen = false
var bg [500]bgLine
var onFloor = false
var jumping = false
var player = playDat{pos: rl.Vector2{X: 0, Y: 0}, vel: rl.Vector2{X: 0, Y: 0}}
var t float32 = 2
var level []line
var playing bool = false
var bv float32 = 0
var mode = 0
var points [16]rl.Vector2
var pointsVel [16]rl.Vector2

const drawSunBool = false

var sunCol = rl.Yellow
var bgCol = rl.Pink
var playAndPlatCol = rl.Black

var jumpTimeLeft = 100
var onLine = line{0, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}

var colors []rl.Color = []rl.Color{rl.Black, rl.White, rl.Brown, rl.Red, rl.Orange, rl.Yellow, rl.Green, rl.Blue, rl.Pink, rl.Purple}
var numKeys []int = []int{rl.KeyOne, rl.KeyTwo, rl.KeyThree, rl.KeyFour, rl.KeyFive, rl.KeySix, rl.KeySeven, rl.KeyEight, rl.KeyNine, rl.KeyZero}
var numKeysMap map[int]int = map[int]int{
	numKeys[0]: 1,
	numKeys[1]: 2,
	numKeys[2]: 3,
	numKeys[3]: 4,
	numKeys[4]: 5,
	numKeys[5]: 6,
	numKeys[6]: 7,
	numKeys[7]: 8,
	numKeys[8]: 9,
	numKeys[9]: 0}

var colorBgBoard board = board{buttonOps: func(k int) {
	if k != -1 {
		bgCol = colors[numKeysMap[k]]
	}
}, activation: rl.KeyB,
	draw: func() {
		rl.DrawText("changing back color", 0, 0, 32, rl.Black)
	},
}

var boardList []board = []board{colorBgBoard}

const maxJumpTime = 150

type board struct {
	buttonOps  func(int)
	activation int
	draw       func()
}

func drawSun() {
	if drawSunBool {
		for i := range 15 {
			r := float32((10*i + int(rl.GetTime()*10)) % 75)
			rl.DrawCircle(int32((-cam.Offset.X+800)*.95), int32((-cam.Offset.Y+500)*.95), 12*r, rl.Color{R: sunCol.R, G: sunCol.G, B: sunCol.B, A: uint8(255 * ((-r + 75) / 75))})
		}
	}
}

func drawPlayer() {
	rl.DrawRectangleRec(rl.Rectangle{
		X:      player.pos.X,
		Y:      player.pos.Y,
		Width:  16,
		Height: 16,
	}, playAndPlatCol)
	rl.DrawCircleV(bomb, 32, playAndPlatCol)
}
func drawBg() {
	for i := range 500 {
		rl.DrawLineEx(rl.Vector2{X: bg[i].start + (-cam.Offset.X+500)*.3, Y: float32(1000)}, rl.Vector2{X: bg[i].start + bg[i].angle + (-cam.Offset.X+500)*.3, Y: -bg[i].length + 1000}, 3.5, playAndPlatCol)
	}

}
func drawPlats() {

	for i := range level {
		rl.DrawLineEx(level[i].Start, level[i].End, 10, playAndPlatCol)
	}
	for i := range levelRect {
		rl.DrawRectangleRec(levelRect[i], playAndPlatCol)
	}
}

func die() {
	player.pos = rl.Vector2{X: 0, Y: 0}
	player.vel = rl.Vector2{X: 0, Y: 0}
	jumpTimeLeft = 0
}

func drawEvil() {
	for i := range levelEvilRect {
		rl.DrawRectangleRec(levelEvilRect[i], rl.White)
		if rl.CheckCollisionRecs(rl.Rectangle{X: player.pos.X, Y: player.pos.Y, Width: 8, Height: 8}, levelEvilRect[i]) {
			die()
		}
	}
}

type side struct {
	pos rl.Vector2
	hit bool
}

func check4To1Line(l1, l2, p11, p12, p21, p22, p31, p32, p41, p42 rl.Vector2) side {
	// println("not")
	var pos rl.Vector2
	if rl.CheckCollisionLines(l1, l2, p11, p12, &pos) {
		return side{pos, true}
	}

	if rl.CheckCollisionLines(l1, l2, p21, p22, &pos) {
		return side{pos, true}
	}

	if rl.CheckCollisionLines(l1, l2, p31, p32, &pos) {
		return side{pos, true}
	}

	if rl.CheckCollisionLines(l1, l2, p41, p42, &pos) {
		return side{pos, true}
	}
	return side{hit: false}
}

type side2 struct {
	pos rl.Vector2
	sde int
}

func getColPoints() []side2 {
	var out []side2

	tl := rl.Vector2{X: player.pos.X, Y: player.pos.Y}
	tr := rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y}
	bl := rl.Vector2{X: player.pos.X, Y: player.pos.Y + 16}
	br := rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y + 16}

	player.pos = player.pos.Add(player.vel)

	ntl := rl.Vector2{X: player.pos.X, Y: player.pos.Y}
	ntr := rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y}
	nbl := rl.Vector2{X: player.pos.X, Y: player.pos.Y + 16}
	nbr := rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y + 16}

	player.pos = player.pos.Subtract(player.vel)

	checkSide := func(l1, l2 rl.Vector2) side {
		return check4To1Line(l1, l2, tl, ntl, tr, ntr, bl, nbl, br, nbr)
	}

	for i := range levelRect {

		//here be point findin
		l := levelRect[i]
		t1 := rl.Vector2{X: l.X, Y: l.Y}
		t2 := rl.Vector2{X: l.X + l.Width, Y: l.Y}
		if checkSide(t1, t2).hit {
			out = append(out, side2{pos: checkSide(t1, t2).pos, sde: 1})
		}

		l1 := rl.Vector2{X: l.X, Y: l.Y}
		l2 := rl.Vector2{X: l.X, Y: l.Y + l.Height}

		if checkSide(l1, l2).hit {
			out = append(out, side2{pos: checkSide(l1, l2).pos, sde: 2})
		}

		r1 := rl.Vector2{X: l.X + l.Width, Y: l.Y}
		r2 := rl.Vector2{X: l.X + l.Width, Y: l.Y + l.Height}

		if checkSide(r1, r2).hit {
			out = append(out, side2{pos: checkSide(r1, r2).pos, sde: 3})
		}

		b1 := rl.Vector2{X: l.X, Y: l.Y + l.Height}
		b2 := rl.Vector2{X: l.X + l.Width, Y: l.Y + l.Height}
		if checkSide(b1, b2).hit {
			out = append(out, side2{pos: checkSide(b1, b2).pos, sde: 4})
		}

	}
	return out
}

func checkCollRect() {

	var add rl.Vector2
	var sideHits [4]bool
	colPoints := getColPoints()
	var minDist float32 = 1e25
	minIndex := -1
	for i := range colPoints {
		sideHits[colPoints[i].sde-1] = true
		if colPoints[i].pos.Distance(player.pos) < minDist {
			minIndex = i
		}
	}
	if minIndex != -1 {

		player.pos = colPoints[minIndex].pos.Add(add)
	}

	if sideHits[0] {
		player.pos.Y -= 16
		player.vel.Y = 0
		onFloor = true
	}
	if sideHits[1] {
		player.pos.X -= 16
		player.vel.X = min(player.vel.X, 0)
	}
	if sideHits[2] {
		player.vel.X = max(player.vel.X, 0)
	}
	if sideHits[3] {
		player.vel.Y = 0
		jumpTimeLeft = 0
	}

}

func checkColl() {
	for i := range level {
		next := rl.Vector2{X: player.pos.X + player.vel.X, Y: player.pos.Y + player.vel.Y}
		if rl.CheckCollisionLines(next, player.pos, level[i].Start, level[i].End, &next) || rl.CheckCollisionLines(player.pos, next, level[i].Start, level[i].End, &next) {
			player.prevPos = player.pos
			player.pos = next
			onLine = level[i]
			t = rl.Vector2{X: player.pos.X - level[i].Start.X, Y: player.pos.Y - level[i].Start.Y}.Length() / rl.Vector2{X: level[i].End.X - level[i].Start.X, Y: level[i].End.Y - level[i].Start.Y}.Length()
			onFloor = true
		}
	}
}

func move() {
	player.vel.Y = max(player.vel.Y, -59)
	bv += rl.GetFrameTime() * 120
	bv *= 1.01
	bomb.Y += bv
	//if bv > 0 {
	//	bv -= 0.05
	//}
	if rl.IsKeyPressed(rl.KeyK) {
		leng := rl.Vector2{X: bomb.X - player.pos.X, Y: bomb.Y - player.pos.Y}
		p := float32(math.Sqrt(math.Pow(float64(leng.X), 2) + math.Pow(float64(leng.Y), 2)))
		if leng.X > 0 {
			player.vel.X -= rl.Clamp(1000/p, -60, 60)
		} else {
			player.vel.X += rl.Clamp(1000/p, -60, 60)
		}
		if leng.Y > 0 {
			player.vel.Y -= rl.Clamp(5000/p, -60, 60)
		} else {
			player.vel.Y += rl.Clamp(400/p, -60, 60)
		}
		bomb = rl.Vector2{X: 10000, Y: 10000}
		bv = 0
	}

	if rl.IsKeyPressed(rl.KeyJ) {
		bomb = player.pos
		bv = player.vel.Y + 12
	}

	player.prevPos = player.pos

	if t < 1.1 && t > 1 {
		checkColl()
		frozen = true
	}
	if player.pos.Y > 777 {
		onFloor = true
		player.pos.Y = 777
		player.vel.Y = 0
		jumping = true
		jumpTimeLeft = maxJumpTime
	}

	jumpTimeLeft -= 10
	if rl.IsKeyDown(rl.KeySpace) && jumpTimeLeft > 0 && jumping {
		player.vel.Y = -30
	}

	if rl.IsKeyDown(rl.KeyA) {
		player.vel.X -= 2
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.vel.X += 2
	}
	checkCollRect()
	if !onFloor {
		player.vel.Y += rl.GetFrameTime() * 120
		player.pos.Y += player.vel.Y
	} else {
		player.vel.Y = -1
		if rl.IsKeyDown(rl.KeySpace) {
			player.pos.Y -= 10
			player.vel.Y = -30
		}
		jumpTimeLeft = maxJumpTime
	}

	player.pos.X += player.vel.X * 2

	checkColl()

	player.vel = player.vel.Multiply(rl.Vector2{X: 0.9, Y: 1.01})

	onFloor = false
	//println(onFloor)

}

func getDownNumberKeys() int32 {
	var output int32
	for i := range 10 {
		if rl.IsKeyDown(int32(numKeys[i])) {
			output = int32(numKeys[i])
			return output
		}
	}
	return -1
}

func main() {

	//decode := gob.NewDecoder(&net)

	cam = rl.Camera2D{Offset: rl.Vector2{X: 0, Y: 0}, Target: rl.Vector2{X: 0, Y: 0}, Rotation: 0, Zoom: 1}
	clickStart := rl.Vector2{X: 999999999, Y: 99999999}
	isHolding := false
	for i := range 500 {
		bg[i].angle = rand.Float32()*3000 - 1500
		bg[i].length = rand.Float32() * 2500
		bg[i].start = rand.Float32()*30000 - 15000
	}

	// rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.InitWindow(1600, 900, "game")
	rl.SetTargetFPS(60)

	// bloom := rl.LoadShader(
	// 	"",
	// 	"bloom.fs",
	// )
	// spire := rl.LoadShader(
	// 	"",
	// 	"spire.fs",
	// )
	tex := rl.LoadRenderTexture(1600, 900)
	tex2 := rl.LoadRenderTexture(1600, 900)
	defer rl.UnloadRenderTexture(tex)

	for i := 0; i < 16; i++ {
		points[i] = rl.Vector2{X: rand.Float32(), Y: rand.Float32()}
	}

	for !rl.WindowShouldClose() {
		for i := 0; i < 16; i++ {
			pointsVel[i].X += rand.Float32() - 0.5
			pointsVel[i].X *= 0.95
			pointsVel[i].Y *= 0.95
			pointsVel[i].Y += rand.Float32() - 0.5
			points[i].X = float32(math.Mod(float64(points[i].X+pointsVel[i].X/1000), 1))
			points[i].Y = float32(math.Mod(float64(points[i].Y+pointsVel[i].Y/1000), 1))

		}

		if rl.IsKeyPressed(rl.KeySpace) && mode == 0 {
			mode = 1
		}
		if rl.IsKeyPressed(rl.KeyL) && mode == 0 {
			mode = 2
		}

		if mode == 0 {

			rl.BeginDrawing()

			rl.ClearBackground(bgCol)

			drawSun()
			rl.EndDrawing()
		}
		if mode == 1 {
			if rl.IsKeyPressed(rl.KeyEnter) && playing {
				mode = 2
				playing = false
			}
			// loc := rl.GetShaderLocation(spire, "points")
			data := make([]float32, len(points)*2)
			for i, v := range points {
				data[i*2] = v.X
				data[i*2+1] = v.Y
			}

			// rl.SetShaderValueV(spire, loc, data, rl.ShaderUniformVec2, int32(len(points)))

			move()
			if frozen {
				frozen = false
				//continue
			}
			const easing = 3
			cam.Offset.X = ((player.pos.X*-1 + 800) + cam.Offset.X*(easing-1)) / easing
			cam.Offset.Y = ((player.pos.Y*-1 + 550) + cam.Offset.Y*(easing-1)) / easing

			rl.BeginTextureMode(tex)
			rl.ClearBackground(bgCol)
			rl.BeginMode2D(cam)
			drawSun()
			drawBg()
			drawPlats()
			drawPlayer()

			rl.DrawRectangle(-100000, 792, 10000890, 1000, playAndPlatCol)
			rl.EndMode2D()
			rl.EndTextureMode()

			rl.BeginTextureMode(tex2)
			rl.ClearBackground(rl.Blank)
			rl.BeginMode2D(cam)
			drawEvil()
			rl.EndMode2D()
			rl.DrawFPS(0, 0)
			rl.EndTextureMode()

			rl.BeginDrawing()
			// rl.BeginShaderMode(bloom)
			rl.ClearBackground(bgCol)
			rl.DrawTextureRec(tex.Texture, rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  1600,
				Height: -900,
			}, rl.Vector2{}, rl.White)
			// rl.EndShaderMode()
			// rl.BeginShaderMode(spire)
			rl.DrawTextureRec(tex2.Texture, rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  1600,
				Height: -900,
			}, rl.Vector2{}, rl.White)
			// rl.EndShaderMode()
			rl.EndDrawing()

		}
		if mode == 2 {

			for i := range boardList {
				if rl.IsKeyDown(int32(boardList[i].activation)) {
					boardList[i].draw()
					boardList[i].buttonOps(int(getDownNumberKeys()))
				}
			}

			if rl.IsKeyPressed(rl.KeyEnter) {
				mode = 1
				playing = true
			}

			if !playing {
				if rl.IsKeyDown(rl.KeyW) {
					cam.Offset.Y += 20
				}
				if rl.IsKeyDown(rl.KeyA) {
					cam.Offset.X += 20
				}
				if rl.IsKeyDown(rl.KeyS) {
					cam.Offset.Y -= 20
				}
				if rl.IsKeyDown(rl.KeyD) {
					cam.Offset.X -= 20
				}
			}

			rl.BeginDrawing()
			rl.BeginMode2D(cam)

			mouse := rl.GetScreenToWorld2D(rl.GetMousePosition(), cam)
			rl.ClearBackground(bgCol)
			drawPlats()
			drawBg()
			drawEvil()
			drawPlayer()
			rl.DrawRectangle(-100000, 792, 10000890, 1000, playAndPlatCol)

			if isHolding {
				sx := min(mouse.X, clickStart.X)
				bx := max(mouse.X, clickStart.X)
				sy := min(mouse.Y, clickStart.Y)
				by := max(mouse.Y, clickStart.Y)
				placedRect := rl.Rectangle{X: sx, Y: sy, Width: bx - sx, Height: by - sy}
				if rl.IsKeyDown(rl.KeyLeftShift) {
					rl.DrawRectangleRec(placedRect, rl.White)
				} else {
					rl.DrawRectangleRec(placedRect, rl.Blue)

				}
				if !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
					if !rl.IsKeyDown(rl.KeyLeftShift) {
						levelRect = append(levelRect, placedRect)
					} else {
						levelEvilRect = append(levelEvilRect, placedRect)
					}
					isHolding = false
				}

			}
			if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !isHolding {
				clickStart = mouse
				isHolding = true
			}

			if !rl.IsMouseButtonDown(rl.MouseButtonLeft) {
				isHolding = false
			}

			rl.DrawFPS(0, 0)

			rl.EndMode2D()
			if rl.IsKeyDown(rl.KeyH) {
				var net bytes.Buffer
				net.Reset()
				encode := gob.NewEncoder(&net)
				err := encode.Encode(levelFull{
					LevelRect:     levelRect,
					LevelEvilRect: levelEvilRect,
					Level:         level,
				})
				rl.SetClipboardText(base64.StdEncoding.EncodeToString(net.Bytes()))

				if err != nil {
					println(err.Error())
				} else {
					rl.DrawText("level saved", 0, 0, 32, playAndPlatCol)
				}

			}
			if rl.IsKeyDown(rl.KeyG) {

				var net bytes.Buffer
				var decodedLev levelFull
				net.Reset()
				WHY, _ := base64.StdEncoding.DecodeString(rl.GetClipboardText())
				net.WriteString(string(WHY))
				decode := gob.NewDecoder(&net)

				err := decode.Decode(&decodedLev)
				level = decodedLev.Level
				levelRect = decodedLev.LevelRect
				levelEvilRect = decodedLev.LevelEvilRect
				if err != nil {
					println(err.Error())
				} else {
					rl.DrawText("level loaded", 0, 0, 32, playAndPlatCol)
				}
			}
			rl.EndDrawing()
		}
	}
}
