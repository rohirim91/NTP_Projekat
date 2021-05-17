package main

import (
	"fmt"
	"image"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
)

type Cluster struct {
	pixels        [][]int
	h, w, l, a, b int
}

func (c *Cluster) update(h, w, l, a, b int) {
	c.h = h
	c.w = w
	c.l = l
	c.a = a
	c.b = b
}

func (c *Cluster) removePixelsValue(h, w int) {
	var remove_idx = 0

	for idx, v := range c.pixels {
		if v[0] == h && v[1] == w {
			remove_idx = idx
			break
		}
	}

	var new_pixels = make([][]int, 0)
	new_pixels = append(new_pixels, c.pixels[:remove_idx]...)
	c.pixels = append(new_pixels, c.pixels[remove_idx+1:]...)
}

func (c *Cluster) removePixelsValue_parallel(h, w int, mutex *sync.Mutex) {
	var remove_idx = 0

	mutex.Lock()
	for idx, v := range c.pixels {
		if v[0] == h && v[1] == w {
			remove_idx = idx
			break
		}
	}

	var new_pixels = make([][]int, 0)
	new_pixels = append(new_pixels, c.pixels[:remove_idx]...)
	c.pixels = append(new_pixels, c.pixels[remove_idx+1:]...)
	mutex.Unlock()
}

type Label struct {
	h, w    int
	cluster *Cluster
}

func repackPixels(img image.Image) [][][]int {
	var bounds = img.Bounds()
	var width, height = bounds.Max.X, bounds.Max.Y

	var pixels [][][]int
	for y := 0; y < height; y++ {
		var row [][]int
		for x := 0; x < width; x++ {
			row = append(row, rgbaToPixel(img.At(x, y).RGBA()))
		}
		pixels = append(pixels, row)
	}

	return pixels
}

func rgbaToPixel(r uint32, g uint32, b uint32, a uint32) []int {
	return []int{int(r / 257), int(g / 257), int(b / 257)}
}

func loadImage(path string) image.Image {
	var img, err = imgio.Open(path)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return img
}
