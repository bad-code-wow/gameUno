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

var frozen = false
var bg [500]bgLine
var onFloor = false
var player playDat = playDat{pos: rl.Vector2{X: 0, Y: 0}, vel: rl.Vector2{X: 0, Y: 0}}
var t float32 = 2
var level []line = []line{{40, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}, {40, rl.Vector2{X: 100, Y: -800}, rl.Vector2{X: -1300, Y: 600}}}

var jumpTimeLeft = 100
var onLine = line{0, rl.Vector2{X: 1300, Y: 600}, rl.Vector2{X: 100, Y: -800}}

const maxJumpTime = 100

func drawPlatsAndPlayer() {

	for i := range 15 {
		r := float32((10*i + int(rl.GetTime()*10)) % 75)
		rl.DrawCircle(int32((-cam.Offset.X+800)*.7), int32((-cam.Offset.Y+500)*.7), 12*r, rl.Color{R: 255, G: 255, B: 0, A: uint8(255 * ((-r + 75) / 75))})
	}

	for i := range 100 {
		rl.DrawLineEx(rl.Vector2{X: bg[i].start + (-cam.Offset.Y+500)*.27, Y: float32(1000)}, rl.Vector2{X: bg[i].start + bg[i].angle + (-cam.Offset.Y+500)*.27, Y: -bg[i].length + 1000}, 1.5, rl.Color{R: 0, G: 0, B: 0, A: 128})
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
			//frozen = true
			checkColl(0)
		}
		return

	}
	if t < 1.1 && t > 1 {
		checkColl(0)
		frozen = true
	}

	jumpTimeLeft -= 10
	if rl.IsKeyDown(rl.KeySpace) && jumpTimeLeft > 0 && player.vel.Y > -30 {
		player.vel.Y = -30
	} else {
	}

	if rl.IsKeyDown(rl.KeyA) {
		player.vel.X -= 2
	}
	if rl.IsKeyDown(rl.KeyD) {
		player.vel.X += 2
	}

	if !onFloor {
		player.vel.Y += 1
		player.pos.Y += player.vel.Y
	} else {
		jumpTimeLeft = maxJumpTime
	}
	player.pos.X += player.vel.X
	player.vel.Y += 1
	if player.pos.Y > 777 {
		player.pos.Y = 777
		player.vel.Y = 0
		jumpTimeLeft = 200
	}
	checkColl(16)
	player.vel = player.vel.Multiply(rl.Vector2{X: 0.9, Y: 1.01})
	player.pos.X += player.vel.X

	onFloor = false

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

	fmt.Println("shader valid:", rl.IsShaderValid(bloom))

	rl.SetTargetFPS(60)

	cam = rl.Camera2D{Offset: rl.Vector2{X: 0, Y: 0}, Target: rl.Vector2{X: 0, Y: 0}, Rotation: 0, Zoom: 1}

	tex := rl.LoadRenderTexture(1600, 900)
	defer rl.UnloadRenderTexture(tex)

	for !rl.WindowShouldClose() {

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

		rl.BeginDrawing()
		rl.ClearBackground(rl.Orange)
		//rl.BeginShaderMode(bloom)
		rl.DrawTextureRec(tex.Texture, rl.Rectangle{
			X:      0,
			Y:      0,
			Width:  1600,
			Height: -900,
		}, rl.Vector2{}, rl.White)
		//rl.EndShaderMode()
		rl.EndDrawing()

	}
}
