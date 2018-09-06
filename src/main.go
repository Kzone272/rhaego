package main

import (
	"os"
	"syscall"
	"image/color"
	"image"
	"image/png"
	"math/rand"
	"reflect"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	f, err := os.OpenFile("output/image.png", syscall.O_WRONLY | syscall.O_CREAT, 0777)
	check(err)

	width := 20
	height := 20
	numPixels := width * height

	chans := make([]chan color.NRGBA, numPixels)
	cases := make([]reflect.SelectCase, numPixels)

	image := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			ch := pixel(x, y)
			i := x * height + y
			chans[i] = ch
			cases[i] = reflect.SelectCase{
				Dir: reflect.SelectRecv,
				Chan: reflect.ValueOf(ch),
			}
		}
	}

	for i := 0; i < numPixels; i++ {
		// reflect.Select is really bad here
		// It's super slow, and it breaks with more than 65536 channels
		// Currently we have 1 channel per pixel (which is probably a bad idea)
		// And we definitely want to have more than 256 * 256 pixel images
		chosen, value, _ := reflect.Select(cases)
		x := chosen / height
		y := chosen % height
		colour := value.Interface().(color.NRGBA)
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
