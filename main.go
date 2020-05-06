package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"math"
	"math/rand"
	"time"
)

var (
	winWidth = 800
	winHeight = 600
	speed = 10
	blockSize = 10
	random = rand.New(rand.NewSource(time.Now().Unix()))
)



func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow(
		"PONG",
		sdl.WINDOWPOS_CENTERED,
		sdl.WINDOWPOS_CENTERED,
		int32(winWidth),
		int32(winHeight),
		sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatal(err)
	}
	defer renderer.Destroy()

	centerX, centerY := getCenter()
	snakeRect := sdl.Rect{
		X: int32(centerX),
		Y: int32(centerY),
		W: 10,
		H: 10,
	}
	fruitRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: 10,
		H: 10,
	}

	background, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		int32(winWidth),
		int32(winHeight))
	if err != nil {
		log.Fatal(err)
	}
	defer background.Destroy()
	snakeTex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		snakeRect.W,
		snakeRect.H)
	if err != nil {
		log.Fatal(err)
	}
	defer snakeTex.Destroy()
	fruitTex, err := renderer.CreateTexture(
		sdl.PIXELFORMAT_ABGR8888,
		sdl.TEXTUREACCESS_STREAMING,
		fruitRect.W,
		fruitRect.H)
	if err != nil {
		log.Fatal(err)
	}
	defer fruitTex.Destroy()

	fruitPixels := make([]byte, fruitRect.H * fruitRect.W * 4)
	for y := 1; y < int(fruitRect.H) - 1; y++ {
		for x := 1; x < int(fruitRect.W) - 1; x++ {
			i := (y * int(fruitRect.W) + x) * 4
			fruitPixels[i]     = 255
			fruitPixels[i + 1] = 0
			fruitPixels[i + 2] = 0
			fruitPixels[i + 3] = 255
		}
	}
	fruitTex.Update(nil, fruitPixels, int(fruitRect.W) * 4)

	snakePixels := make([]byte, snakeRect.H * snakeRect.W * 4)
	for y := 1; y < int(snakeRect.H) - 1; y++ {
		for x := 1; x < int(snakeRect.W) - 1; x++ {
			i := (y * int(snakeRect.W) + x) * 4
			snakePixels[i]     = 255
			snakePixels[i + 1] = 255
			snakePixels[i + 2] = 255
			snakePixels[i + 3] = 255
		}
	}
	snakeTex.Update(nil, snakePixels, int(snakeRect.W) * 4)

	keyState := sdl.GetKeyboardState()
	pixels := make([]byte, winHeight*winWidth*4)
	background.Update(nil, pixels, winWidth * 4)
	renderer.Copy(background, nil, nil)

	snake := &snakeBody{
		Rect: snakeRect,
		dx:   0,
		dy:   -snakeRect.H,
		tex:  snakeTex,
		next: nil,
	}
	fruit := &fruit{
		Rect: fruitRect,
		tex:  fruitTex,
	}

	fruit.newLoc(snake, winWidth - int(fruit.W), winHeight - int(fruit.H))

	snake.draw(renderer)
	fruit.draw(renderer)
	renderer.Present()
	frameStart := time.Now()
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				return
			}
		}
		if keyState[sdl.SCANCODE_UP] != 0 {
			if snake.canMoveTo(0, -snakeRect.H) {
				snake.dy = -snakeRect.H
				snake.dx = 0
			}
		} else if keyState[sdl.SCANCODE_DOWN] != 0 {
			if snake.canMoveTo(0, snakeRect.H) {
				snake.dy = snakeRect.H
				snake.dx = 0
			}
		} else if keyState[sdl.SCANCODE_RIGHT] != 0 {
			if snake.canMoveTo(snakeRect.W, 0) {
				snake.dy = 0
				snake.dx = snakeRect.W
			}
		} else if keyState[sdl.SCANCODE_LEFT] != 0 {
			if snake.canMoveTo(-snakeRect.W, 0) {
				snake.dy = 0
				snake.dx = -snakeRect.W
			}
		}
		if time.Since(frameStart).Milliseconds() >= 150 {
			renderer.Clear()
			renderer.Copy(background, nil, nil)

			snake.update()
			if snake.detectBodyCollision() || snake.X < 0 || snake.X >= int32(winWidth) || snake.Y < 0 || snake.Y >= int32(winHeight) {
				log.Println("Game Over")
				break
			}

			if snake.foundFruit(*fruit) {
				log.Println("Found Fruit")
				snake.grow()
				log.Println("Grown")
				fruit.newLoc(snake, winWidth-int(fruit.W), winHeight-int(fruit.H))
				log.Println("New Fruit location")
			}

			snake.draw(renderer)
			fruit.draw(renderer)

			renderer.Present()
			frameStart = time.Now()
		}
		sdl.Delay(10)
	}
}

func isEmpty(rect sdl.Rect, s *snakeBody) bool {
	for s != nil {
		if rect.X == s.Rect.X && rect.Y == s.Rect.Y {
			return false
		}
		s = s.next
	}
	return true
}

func getRandomRect(w, h int) sdl.Rect {
	randX := int32(random.Intn(w))
	randY := int32(random.Intn(h))
	return sdl.Rect{
		X: round(randX, 10),
		Y: round(randY, 10),
	}
}

func round(x, unit int32) int32 {
	return int32(math.Round(float64(x/unit))) * unit
}

func getCenter() (int, int) {
	return winWidth / 2, winHeight / 2
}