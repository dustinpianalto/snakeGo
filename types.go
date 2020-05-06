package main

import "github.com/veandco/go-sdl2/sdl"

type snakeBody struct {
	sdl.Rect
	dx		int32
	dy		int32
	tex 	*sdl.Texture
	next	*snakeBody
}

type fruit struct {
	sdl.Rect
	tex 	*sdl.Texture
}

func (s *snakeBody) detectCollision(rect sdl.Rect) bool {
	return s.X == rect.X && s.Y == rect.Y
}

func (s *snakeBody) foundFruit(fruit fruit) bool {
	return s.detectCollision(fruit.Rect)
}

func (s *snakeBody) detectBodyCollision() bool {
	nextChunk := s.next
	for nextChunk != nil {
		if s.detectCollision(nextChunk.Rect) {
			return true
		}
		nextChunk = nextChunk.next
	}
	return false
}

func (s *snakeBody) draw(renderer *sdl.Renderer) {
	chunk := s
	for chunk != nil {
		renderer.Copy(chunk.tex, nil, &chunk.Rect)
		chunk = chunk.next
	}
}

func (s *snakeBody) update() {
	prevRect := s.Rect
	s.X += s.dx
	s.Y += s.dy
	next := s.next
	for next != nil {
		nextRect := next.Rect
		next.Rect = prevRect
		prevRect = nextRect
		next = next.next
	}
}

func (s *snakeBody) grow() {
	last := s
	for {
		if last.next == nil {
			break
		}
		last = last.next
	}
	var newSnake *snakeBody
	for i := 0; i < 3; i++ {
		newSnake = &snakeBody{
			Rect: last.Rect,
			dx:   0,
			dy:   0,
			tex:  last.tex,
			next: newSnake,
		}
	}
	last.next = newSnake
}

func (s *snakeBody) canMoveTo(dx, dy int32) bool {
	if s.next == nil {
		return true
	}
	nextRect := s.next.Rect
	if s.X + dx == nextRect.X && s.Y + dy == nextRect.Y {
		return false
	}
	return true
}

func (f *fruit) draw(renderer *sdl.Renderer) {
	renderer.Copy(f.tex, nil, &f.Rect)
}

func (f *fruit) newLoc(s *snakeBody, w, h int) {
	rect := getRandomRect(w, h)
	for !isEmpty(rect, s) {
		rect = getRandomRect(w, h)
	}
	f.X = rect.X
	f.Y = rect.Y
}