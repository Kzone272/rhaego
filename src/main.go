package main

import (
	"os"
	"syscall"
	"image/color"
	"image"
	"image/png"
	"math/rand"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.OpenFile("output/image.png", syscall.O_WRONLY | syscall.O_CREAT, 0777)
	check(err)

	width := 500
	height := 500
	numPixels := width * height

	ch := make(chan PixelMessage)

	image := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel(x, y, ch)
		}
	}

	for i := 0; i < numPixels; i++ {
		pixel := <- ch
		image.Set(pixel.x, pixel.y, pixel.colour)
	}

	err = png.Encode(f, image)
	check(err)

	err = f.Close()
	check(err)
}

// PixelMessage is the position and colour of a pixel
type PixelMessage struct {
	x int
	y int
	colour color.NRGBA
}

func pixel(x int, y int, ch chan PixelMessage) {
	go func () {
		colour := color.NRGBA{0, 128, 255, uint8(rand.Int() % 256)} // azure
		ch <- PixelMessage{ x, y, colour }
	}()
}
