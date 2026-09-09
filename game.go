package main

import (
	"fmt"
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
	power float32
	start rl.Vector2
	end   rl.Vector2
}

type bgLine struct {
	length float32
	start  float32
	angle  float32
}

var levelRect []rl.Rectangle = []rl.Rectangle{{X: 90, Y: 90, Width: 190, Height: 190}, {X: 200, Y: 0, Width: 190, Height: 300}}
var levelEvilRect []rl.Rectangle = []rl.Rectangle{{X: 390, Y: 390, Width: 190, Height: 190}}
var frozen = false
var bg [500]bgLine
var onFloor = false
var jumping = false
var player playDat = playDat{pos: rl.Vector2{X: 0, Y: 0}, vel: rl.Vector2{X: 0, Y: 0}}
var t float32 = 2
var level []line = []line{{40, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}, {40, rl.Vector2{X: 100, Y: -800}, rl.Vector2{X: -1300, Y: 600}}}

var jumpTimeLeft = 100
var onLine = line{0, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}

const maxJumpTime = 150

func drawPlatsAndPlayer() {

	for i := range 15 {
		r := float32((10*i + int(rl.GetTime()*10)) % 75)
		rl.DrawCircle(int32((-cam.Offset.X+800)*.7), int32((-cam.Offset.Y+500)*.7), 12*r, rl.Color{R: 255, G: 255, B: 0, A: uint8(255 * ((-r + 75) / 75))})
	}

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
		rl.DrawLineEx(level[i].start, level[i].end, 10, rl.Black)
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
		if rl.CheckCollisionRecs(rl.Rectangle{player.pos.X, player.pos.Y, 8, 8}, levelEvilRect[i]) {
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
		if rl.CheckCollisionLines(rl.Vector2{X: player.pos.X + 16, Y: player.pos.Y}, rl.Vector2{X: player.prevPos.X + 16, Y: player.prevPos.Y}, rl.Vector2{X: levelRect[i].X, Y: levelRect[i].Y - 16}, rl.Vector2{X: levelRect[i].X, Y: levelRect[i].Y + levelRect[i].Height}, &p1) && player.pos.Y+14 > levelRect[i].Y {
			fmt.Println("l")
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
			player.pos.Y = p1.Y
			player.vel.Y = max(2, player.vel.Y)
			jumpTimeLeft = 0
			player.pos.Y += player.vel.Y

		}

		if rl.CheckCollisionLines(rl.Vector2{X: player.pos.X, Y: player.pos.Y + 16}, rl.Vector2{X: player.prevPos.X, Y: player.prevPos.Y + 16}, rl.Vector2{X: levelRect[i].X - 15, Y: levelRect[i].Y}, rl.Vector2{X: levelRect[i].X + levelRect[i].Width, Y: levelRect[i].Y}, &p1) {
			player.pos.Y = p1.Y - 16
			jumpTimeLeft = maxJumpTime
			player.vel.Y = min(0, player.vel.Y)
			player.pos.Y += player.vel.Y
			onFloor = true
		}
	}
}

func checkColl(r float32) {
	for i := range level {
		next := rl.Vector2{X: player.pos.X + player.vel.X, Y: player.pos.Y + player.vel.Y}
		if rl.CheckCollisionLines(player.pos, player.prevPos, level[i].start, level[i].end, &next) || rl.CheckCollisionLines(player.pos, next, level[i].start, level[i].end, &next) {
			fmt.Print("yooooou\n")
			player.prevPos = player.pos
			player.pos = next
			onLine = level[i]
			t = rl.Vector2{X: player.pos.X - level[i].start.X, Y: player.pos.Y - level[i].start.Y}.Length() / rl.Vector2{X: level[i].end.X - level[i].start.X, Y: level[i].end.Y - level[i].start.Y}.Length()
			onFloor = true
		}
	}
}

func move() {
	player.prevPos = player.pos
	t += rl.GetFrameTime() * onLine.power / 25
	if t < 1 {
		onFloor = true

		player.pos.X = rl.Vector2Lerp(onLine.start, onLine.end, t).X
		if onLine.start.Y < onLine.end.Y {
			player.pos.Y = rl.Vector2Lerp(onLine.start, onLine.end, t).Y - 10
		} else {
			player.pos.Y = rl.Vector2Lerp(onLine.start, onLine.end, t).Y + 10
		}
		k := rl.Vector2{X: onLine.end.X - onLine.start.X, Y: onLine.end.Y - onLine.start.Y}
		player.vel = rl.Vector2{X: k.X / 25, Y: k.Y / 25}
		if rl.IsKeyPressed(rl.KeySpace) {
			t = 100
			player.vel.Y = -30
			//frozen = true
			checkColl(0)
		}
		return

	}
	if t < 1.1 && t > 1 {
		checkCollRect()
		checkColl(0)
		frozen = true
	}
	checkCollRect()
	if player.pos.Y > 777 {
		checkCollRect()
		onFloor = true
		player.pos.Y = 777
		player.vel.Y = 0
		jumpTimeLeft = maxJumpTime
	}

	checkCollRect()
	jumpTimeLeft -= 10
	if rl.IsKeyDown(rl.KeySpace) && jumpTimeLeft > 0 {
		player.vel.Y = -30
		fmt.Println(player.vel.Y)
	}

	checkCollRect()
	if rl.IsKeyDown(rl.KeyA) {
		player.vel.X -= 2
	}
	checkCollRect()
	if rl.IsKeyDown(rl.KeyD) {
		player.vel.X += 2
	}
	checkCollRect()
	if !onFloor {
		checkCollRect()
		player.vel.Y += 1
	} else {
		jumpTimeLeft = maxJumpTime
	}

	player.pos.X += player.vel.X

	checkColl(16)

	player.vel = player.vel.Multiply(rl.Vector2{X: 0.9, Y: 1.01})

	player.pos.Y += player.vel.Y
	player.pos.X += player.vel.X
	checkCollRect()

	checkCollRect()
	onFloor = false
	checkCollRect()

	checkCollRect()
}

func main() {
	for i := range 500 {
		bg[i].angle = rand.Float32()*1000 - 500
		bg[i].length = rand.Float32() * 1000
		bg[i].start = rand.Float32()*3000 - 1500
	}

	// rl.ToggleFullscreen()
	defer rl.CloseWindow()
	rl.InitWindow(1600, 900, "game")

	bloom := rl.LoadShader(
		"",
		"bloom.fs",
	)
	spire := rl.LoadShader(
		"",
		"spire.fs",
	)

	sTime := rl.GetShaderLocation(spire, "time")

	fmt.Println("shader valid:", rl.IsShaderValid(bloom))

	rl.SetTargetFPS(60)

	cam = rl.Camera2D{Offset: rl.Vector2{X: 0, Y: 0}, Target: rl.Vector2{X: 0, Y: 0}, Rotation: 0, Zoom: 1}

	tex := rl.LoadRenderTexture(1600, 900)
	tex2 := rl.LoadRenderTexture(1600, 900)
	defer rl.UnloadRenderTexture(tex)

	for !rl.WindowShouldClose() {

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
		rl.EndShaderMode()
		rl.BeginShaderMode(spire)
		rl.DrawTextureRec(tex2.Texture, rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  1600,
			Height: -900,
		}, rl.Vector2{}, rl.White)
		rl.EndShaderMode()
		rl.EndDrawing()

	}
}
