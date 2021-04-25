package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"

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

type Label struct {
	h, w    int
	cluster *Cluster
}

type SuperpixelsProcessor struct {
	image        [][][]int
	img_w, img_h int
	K, N, S      int
	M            float64
	clusters     []Cluster
	label        []Label
	dist         [][]float64
}

func (sp *SuperpixelsProcessor) initialize() {
	sp.N = sp.img_h * sp.img_w
	sp.S = int(math.Sqrt(float64(sp.N / sp.K)))
	sp.dist = make([][]float64, sp.img_h)

	for i := 0; i < sp.img_h; i++ {
		sp.dist[i] = make([]float64, sp.img_w)

		for j := 0; j < sp.img_w; j++ {
			sp.dist[i][j] = math.Inf(1)
		}
	}
}

func (sp *SuperpixelsProcessor) initClusters() {
	var w = int(sp.S / 2)
	var h = int(sp.S / 2)

	for {
		if h >= sp.img_h {
			break
		}
		for {
			if w >= sp.img_w {
				break
			}
			sp.clusters = append(sp.clusters, Cluster{h: h, w: w, l: sp.image[h][w][0], a: sp.image[h][w][1], b: sp.image[h][w][2]})
			w += sp.S
		}
		w = sp.S / 2
		h += sp.S
	}
}

func (sp *SuperpixelsProcessor) getGradient(w, h int) int {
	if w+1 >= sp.img_w {
		w = sp.img_w - 2
	}
	if h+1 >= sp.img_h {
		h = sp.img_h - 2
	}

	return sp.image[h+1][w+1][0] - sp.image[h][w][0] +
		sp.image[h+1][w+1][1] - sp.image[h][w][1] +
		sp.image[h+1][w+1][2] - sp.image[h][w][2]
}

func (sp *SuperpixelsProcessor) moveClusters() {
	var c_grad, _h, _w, new_grad = 0, 0, 0, 0

	for i := 0; i < len(sp.clusters); i++ {
		c_grad = sp.getGradient(sp.clusters[i].h, sp.clusters[i].w)

		for dh := -1; dh < 2; dh++ {
			for dw := -1; dw < 2; dw++ {
				_h = int(math.Max(float64(sp.clusters[i].h+dh), 0))
				_w = int(math.Max(float64(sp.clusters[i].w+dw), 0))

				new_grad = sp.getGradient(_h, _w)

				if new_grad < c_grad {
					sp.clusters[i].update(_h, _w, sp.image[_h][_w][0], sp.image[_h][_w][1], sp.image[_h][_w][2])
					c_grad = new_grad
				}
			}
		}
	}
}

func (sp *SuperpixelsProcessor) updateCluster() {
	for i := 0; i < len(sp.clusters); i++ {
		if len(sp.clusters[i].pixels) != 0 {
			var sum_h, sum_w = 0, 0

			for j := 0; j < len(sp.clusters[i].pixels); j++ {
				sum_h += sp.clusters[i].pixels[j][0]
				sum_w += sp.clusters[i].pixels[j][1]
			}

			var _h = int(sum_h / len(sp.clusters[i].pixels))
			var _w = int(sum_w / len(sp.clusters[i].pixels))
			sp.clusters[i].update(_h, _w, sp.image[_h][_w][0], sp.image[_h][_w][1], sp.image[_h][_w][2])
		}
	}
}

func (sp *SuperpixelsProcessor) checkLabel(h, w int) bool {
	for _, v := range sp.label {
		if v.h == h && v.w == w {
			return true
		}
	}
	return false
}

func (sp *SuperpixelsProcessor) updateLabel(h, w int, cluster *Cluster) {
	for i := 0; i < len(sp.label); i++ {
		if sp.label[i].h == h && sp.label[i].w == w {
			sp.label[i].cluster = cluster
			break
		}
	}
}

func (sp *SuperpixelsProcessor) removeLabelPixelsValue(h, w int) {
	for i := 0; i < len(sp.label); i++ {
		if sp.label[i].h == h && sp.label[i].w == w {
			sp.label[i].cluster.removePixelsValue(h, w)
			break
		}
	}
}

func (sp *SuperpixelsProcessor) assign() {
	for i := 0; i < len(sp.clusters); i++ {
		for h := sp.clusters[i].h - 2*sp.S; h < sp.clusters[i].h+2*sp.S; h++ {
			if h < 0 || h >= sp.img_h {
				continue
			}
			for w := sp.clusters[i].w - 2*sp.S; w < sp.clusters[i].w+2*sp.S; w++ {
				if w < 0 || w >= sp.img_w {
					continue
				}
				var l, a, b = sp.image[h][w][0], sp.image[h][w][1], sp.image[h][w][2]
				var dc = math.Sqrt(math.Pow(float64(l-sp.clusters[i].l), 2) +
					math.Pow(float64(a-sp.clusters[i].a), 2) +
					math.Pow(float64(b-sp.clusters[i].b), 2))
				var ds = math.Sqrt(math.Pow(float64(h-sp.clusters[i].h), 2) +
					math.Pow(float64(w-sp.clusters[i].w), 2))
				var d = math.Sqrt(math.Pow(dc/sp.M, 2) + math.Pow(ds/float64(sp.S), 2))

				if d < sp.dist[h][w] {
					if !sp.checkLabel(h, w) {
						sp.clusters[i].pixels = append(sp.clusters[i].pixels, []int{h, w})
						sp.label = append(sp.label, Label{h: h, w: w, cluster: &sp.clusters[i]})
					} else {
						sp.removeLabelPixelsValue(h, w)
						sp.clusters[i].pixels = append(sp.clusters[i].pixels, []int{h, w})
						sp.updateLabel(h, w, &sp.clusters[i])
					}
					sp.dist[h][w] = d
				}
			}
		}

	}
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

func (sp *SuperpixelsProcessor) saveImage(path string, img image.Image) {
	var bounds = img.Bounds()
	var width, height = bounds.Max.X, bounds.Max.Y

	var upLeft = image.Point{0, 0}
	var lowRight = image.Point{width, height}

	var new_img = image.NewRGBA(image.Rectangle{upLeft, lowRight})

	for i := 0; i < len(sp.clusters); i++ {
		for j := 0; j < len(sp.clusters[i].pixels); j++ {
			new_img.Set(sp.clusters[i].pixels[j][1], sp.clusters[i].pixels[j][0], color.RGBA{uint8(sp.clusters[i].l), uint8(sp.clusters[i].a), uint8(sp.clusters[i].b), 0xff})
		}
		//new_img.Set(sp.clusters[i].w, sp.clusters[i].h, color.RGBA{0x0, 0x0, 0x0, 0xff})
	}

	f, _ := os.Create(path)
	png.Encode(f, new_img)
}
