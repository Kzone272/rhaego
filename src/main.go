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

	chans := make([]chan color.NRGBA, numPixels)

	image := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			ch := pixel(x, y)
			i := x * height + y
			chans[i] = ch
		}
	}

	for i := 0; i < numPixels; i++ {
		x := i / height
		y := i % height
		colour := <- chans[i]
		image.Set(x, y, colour)
	}

	err = png.Encode(f, image)
	check(err)

	err = f.Close()
	check(err)
}

func pixel(x int, y int) chan color.NRGBA {
	ch := make(chan color.NRGBA)

	go func () {
		colour := color.NRGBA{0, 128, 255, uint8(rand.Int() % 256)} // azure
		ch <- colour
	}()

	return ch
}
