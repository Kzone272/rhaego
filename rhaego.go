package main

import (
	"math"
	_ "fmt"
	"os"
	"syscall"
	"image/color"
	"image"
	"image/png"
	_ "math/rand"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/barnex/fmath"
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

	ratio := float32(width) / float32(height)
	fov := 45.0
	planeHeight := float32(math.Tan(fov / 180 * math.Pi))
	planeWidth := planeHeight * ratio

	// @TODO make a proper world
	ball := Sphere{
		mgl32.Vec4{0, 0, -7, 0},
		1,
	}

	// Is this buffered channel unnecessairly large?
	ch := make(chan PixelMessage, numPixels)

	image := image.NewNRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			ray := mgl32.Vec4{
				planeWidth * (float32(x) / float32(width) * 2 - 1),
				-planeHeight * (float32(y) / float32(height) * 2 - 1),
				-1,
				0,
			}
			ray = ray.Normalize()
			cast(x, y, ray, ball, ch)
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

func cast(x int, y int, ray mgl32.Vec4, obj Object, ch chan PixelMessage) {
	go func () {
		hit, dist := obj.intersect(ray)
		var colour color.NRGBA
		if hit {
			distScaler := fmath.Min(1, (dist - 5.7) * 2);
			colour = color.NRGBA{0, uint8(distScaler * 128), uint8(distScaler * 255), 255} // azure
		} else {
			colour = color.NRGBA{0, 0, 0, 255} // black
		}
		ch <- PixelMessage{ x, y, colour }
	}()
}

// Object - Intersectable object in the scene
type Object interface {
	// returns t, distance from ray origin
	intersect(ray mgl32.Vec4) (bool, float32)
}

// Sphere - renderable sphere
type Sphere struct {
	pos mgl32.Vec4
	radius float32
}

// Equation from: https://en.wikipedia.org/wiki/Line%E2%80%93sphere_intersection
// ray must be a unit vector
func (s Sphere) intersect(ray mgl32.Vec4) (bool, float32) {
	dot := mgl32.Vec4.Dot(ray, s.pos)
	underRoot := dot - s.pos.Len() + fmath.Pow(s.radius, 2)

	if (underRoot < 0) {
		return false, 0
	}
	root := fmath.Sqrt(underRoot)

	t := fmath.Min(dot + root, dot - root)

	// This either means the object is behind us or we're inside it
	// Either way we don't want to render it
	if (t < 0) {
		return false, 0
	}

	return true, t
}