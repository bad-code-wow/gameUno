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
var levelRect = []rl.Rectangle{{X: 90, Y: 90, Width: 190, Height: 190}, {X: 200, Y: 0, Width: 190, Height: 300}}
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

var jumpTimeLeft = 100
var onLine = line{0, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}

const maxJumpTime = 150

func drawPlatsAndPlayer() {

	for i := range 15 {
		r := float32((10*i + int(rl.GetTime()*10)) % 75)
		rl.DrawCircle(int32((-cam.Offset.X+800)*.7), int32((-cam.Offset.Y+500)*.7), 12*r, rl.Color{R: 255, G: 255, B: 0, A: uint8(255 * ((-r + 75) / 75))})
	}
	rl.DrawCircleV(bomb, 32, rl.Black)

	for i := range 500 {
		rl.DrawLineEx(rl.Vector2{X: bg[i].start + (-cam.Offset.Y+500)*.27, Y: float32(1000)}, rl.Vector2{X: bg[i].start + bg[i].angle + (-cam.Offset.Y+500)*.27, Y: -bg[i].length + 1000}, 3.5, rl.Color{R: 0, G: 0, B: 0, A: 255})
	}
	rl.DrawRectangleRec(rl.Rectangle{
		X:      player.pos.X,
		Y:      player.pos.Y,
		Width:  16,
		Height: 16,
	}, rl.Black)
	for i := range level {
		rl.DrawLineEx(level[i].Start, level[i].End, 10, rl.Black)
	}
	for i := range levelRect {
		rl.DrawRectangleRec(levelRect[i], rl.Black)
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

func checkCollRect() {
	/*
	*~*~*~*~*~*~*

	*magic numbers = magic code*

	*~*~*~*~*~*~*

	 */
	for i := range levelRect {
		p1 := rl.Vector2{X: 0, Y: 0}
		if rl.CheckCollisionLines(rl.Vector2{X: player.pos.X, Y: player.pos.Y + 16}, rl.Vector2{X: player.prevPos.X, Y: player.prevPos.Y + 16}, rl.Vector2{X: levelRect[i].X - 15, Y: levelRect[i].Y}, rl.Vector2{X: levelRect[i].X + levelRect[i].Width, Y: levelRect[i].Y}, &p1) {
			player.pos.Y = p1.Y - 16
			jumpTimeLeft = maxJumpTime
			player.vel.Y = min(0, player.vel.Y)
			player.pos.Y += player.vel.Y
			onFloor = true
			jumping = true
			return
		}

		if rl.CheckCollisionLines(rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y}, rl.Vector2{X: player.prevPos.X + 16, Y: player.prevPos.Y}, rl.Vector2{X: levelRect[i].X, Y: levelRect[i].Y - 16}, rl.Vector2{X: levelRect[i].X, Y: levelRect[i].Y + levelRect[i].Height}, &p1) && player.pos.Y+14 > levelRect[i].Y {
			//fmt.Println("l")
			if rl.IsKeyDown(rl.KeyA) || player.vel.X < 0 && !rl.IsKeyDown(rl.KeyD) {
				player.prevPos = player.pos
				player.pos.X += player.vel.X
			} else {

				player.pos.X = p1.X - 16
				player.vel.X = 0
			}
		}
		if rl.CheckCollisionLines(player.pos, player.prevPos, rl.Vector2{X: levelRect[i].X + levelRect[i].Width, Y: levelRect[i].Y - 16}, rl.Vector2{X: levelRect[i].X + levelRect[i].Width, Y: levelRect[i].Y + levelRect[i].Height}, &p1) && player.pos.Y+14 > levelRect[i].Y {

			if rl.IsKeyDown(rl.KeyD) {
				player.prevPos = player.pos
				player.pos.X += player.vel.X
			} else {

				player.pos.X = p1.X
				player.vel.X = 0
			}
		}

		if rl.CheckCollisionLines(player.pos, player.prevPos, rl.Vector2{X: levelRect[i].X - 15, Y: levelRect[i].Y + levelRect[i].Height}, rl.Vector2{X: levelRect[i].X + levelRect[i].Width, Y: levelRect[i].Y + levelRect[i].Height}, &p1) {
			player.pos.Y = p1.Y + 1
			player.vel.Y = max(2, player.vel.Y)
			jumpTimeLeft = 0
			//player.pos.Y += player.vel.Y

		}

	}
}

func checkColl() {
	for i := range level {
		next := rl.Vector2{X: player.pos.X + player.vel.X, Y: player.pos.Y + player.vel.Y}
		if rl.CheckCollisionLines(player.pos, player.prevPos, level[i].Start, level[i].End, &next) || rl.CheckCollisionLines(player.pos, next, level[i].Start, level[i].End, &next) {
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
		bv = player.vel.Y * 1.1
	}

	player.prevPos = player.pos
	t += rl.GetFrameTime() * onLine.Power / 25
	if t < 1 {
		onFloor = true

		player.pos.X = rl.Vector2Lerp(onLine.Start, onLine.End, t).X
		if onLine.Start.Y < onLine.End.Y {
			player.pos.Y = rl.Vector2Lerp(onLine.Start, onLine.End, t).Y - 10
		} else {
			player.pos.Y = rl.Vector2Lerp(onLine.Start, onLine.End, t).Y + 10
		}
		k := rl.Vector2{X: onLine.End.X - onLine.Start.X, Y: onLine.End.Y - onLine.Start.Y}
		player.vel = rl.Vector2{X: k.X / 25, Y: k.Y / 25}
		if rl.IsKeyPressed(rl.KeySpace) {
			t = 100
			player.vel.Y = -30
			//frozen = true
			checkColl()
		}
		return

	}
	if t < 1.1 && t > 1 {
		//checkCollRect()
		checkColl()
		frozen = true
	}
	//checkCollRect()
	if player.pos.Y > 777 {
		//checkCollRect()
		onFloor = true
		player.pos.Y = 777
		player.vel.Y = 0
		jumping = true
		jumpTimeLeft = maxJumpTime
	}

	//checkCollRect()
	jumpTimeLeft -= 10
	if rl.IsKeyDown(rl.KeySpace) && jumpTimeLeft > 0 && jumping {
		player.vel.Y = -30
	}

	checkCollRect()
	if rl.IsKeyDown(rl.KeyA) {
		player.vel.X -= 2
	}
	checkCollRect()
	if rl.IsKeyDown(rl.KeyD) {
		player.vel.X += 2
	}
	//checkCollRect()
	if !onFloor {
		//checkCollRect()
		player.vel.Y += rl.GetFrameTime() * 120
	} else {
		jumpTimeLeft = maxJumpTime
	}

	player.pos.X += player.vel.X

	checkColl()

	player.vel = player.vel.Multiply(rl.Vector2{X: 0.9, Y: 1.01})

	player.pos.Y += player.vel.Y
	player.pos.X += player.vel.X
	checkCollRect()

	onFloor = false
	//checkCollRect()
	//println(onFloor)
	checkCollRect()

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

	bloom := rl.LoadShader(
		"",
		"bloom.fs",
	)
	spire := rl.LoadShader(
		"",
		"spire.fs",
	)
	tex := rl.LoadRenderTexture(1600, 900)
	tex2 := rl.LoadRenderTexture(1600, 900)
	defer rl.UnloadRenderTexture(tex)

	for !rl.WindowShouldClose() {

		if rl.IsKeyPressed(rl.KeySpace) && mode == 0 {
			mode = 1
		}
		if rl.IsKeyPressed(rl.KeyL) && mode == 0 {
			mode = 2
		}

		if mode == 0 {
			rl.BeginDrawing()
			rl.ClearBackground(rl.Orange)
			for i := range 15 {
				r := float32((10*i + int(rl.GetTime()*10)) % 75)
				rl.DrawCircle(int32((-cam.Offset.X+800)*.7), int32((-cam.Offset.Y+500)*.7), 12*r, rl.Color{R: 255, G: 255, B: 0, A: uint8(255 * ((-r + 75) / 75))})
			}
			rl.EndDrawing()
		}
		if mode == 1 {
			if rl.IsKeyPressed(rl.KeyEnter) && playing {
				mode = 2
				playing = false
			}
			sTime := rl.GetShaderLocation(spire, "time")

			rl.SetShaderValue(
				spire,
				sTime,
				[]float32{float32(rl.GetTime())},
				rl.ShaderUniformFloat,
			)

			move()
			if frozen {
				frozen = false
				//continue
			}
			const easing = 3
			cam.Offset.X = ((player.pos.X*-1 + 800) + cam.Offset.X*(easing-1)) / easing
			cam.Offset.Y = ((player.pos.Y*-1 + 550) + cam.Offset.Y*(easing-1)) / easing

			rl.BeginTextureMode(tex)
			rl.ClearBackground(rl.Orange)
			rl.BeginMode2D(cam)
			drawPlatsAndPlayer()
			rl.DrawRectangle(-100000, 792, 10000890, 1000, rl.Black)
			rl.EndMode2D()
			rl.EndTextureMode()

			rl.BeginTextureMode(tex2)
			rl.ClearBackground(rl.Blank)
			rl.BeginMode2D(cam)
			drawEvil()
			rl.EndMode2D()
			rl.EndTextureMode()

			rl.BeginDrawing()
			rl.BeginShaderMode(bloom)
			rl.ClearBackground(rl.Orange)
			rl.DrawTextureRec(tex.Texture, rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  1600,
				Height: -900,
			}, rl.Vector2{}, rl.White)
			rl.BeginShaderMode(spire)
			rl.DrawTextureRec(tex2.Texture, rl.Rectangle{
				X:      0,
				Y:      0,
				Width:  1600,
				Height: -900,
			}, rl.Vector2{}, rl.White)
			rl.EndShaderMode()
			rl.EndShaderMode()
			rl.EndDrawing()

		}
		if mode == 2 {
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
			rl.ClearBackground(rl.Orange)
			drawPlatsAndPlayer()
			drawEvil()
			rl.DrawRectangle(-100000, 792, 10000890, 1000, rl.Black)

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
					rl.DrawText("level saved", 0, 0, 32, rl.Black)
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
					rl.DrawText("level loaded", 0, 0, 32, rl.Black)
				}
			}
			rl.EndDrawing()
		}
	}
}
