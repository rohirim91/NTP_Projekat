package main

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"sync"
)

type SuperpixelsProcessorParallel struct {
	image        [][][]int
	img_w, img_h int
	K, N, S      int
	M            float64
	clusters     []Cluster
	label        []Label
	dist         [][]float64
}

func (sp *SuperpixelsProcessorParallel) initialize() {
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

func (sp *SuperpixelsProcessorParallel) initClusters() {
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

func (sp *SuperpixelsProcessorParallel) getGradient(w, h int) int {
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

func (sp *SuperpixelsProcessorParallel) moveClusters() {
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

func (sp *SuperpixelsProcessorParallel) updateCluster() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	for i := 0; i < 4; i++ {
		go sp._updateCluster_parallel(i*len(sp.clusters)/4, (i+1)*len(sp.clusters)/4, &waitGroup)
	}
	waitGroup.Wait()
}

func (sp *SuperpixelsProcessorParallel) _updateCluster_parallel(startIdx, endIdx int, wg *sync.WaitGroup) {
	for i := startIdx; i < endIdx; i++ {
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
	wg.Done()
}

func (sp *SuperpixelsProcessorParallel) checkLabel(h, w int) bool {
	for _, v := range sp.label {
		if v.h == h && v.w == w {
			return true
		}
	}
	return false
}

func (sp *SuperpixelsProcessorParallel) updateLabel(h, w int, cluster *Cluster, mutex *sync.Mutex) {
	for i := 0; i < len(sp.label); i++ {
		if sp.label[i].h == h && sp.label[i].w == w {
			mutex.Lock()
			sp.label[i].cluster = cluster
			mutex.Unlock()
			break
		}
	}
}

func (sp *SuperpixelsProcessorParallel) removeLabelPixelsValue(h, w int, mutex *sync.Mutex) {
	for i := 0; i < len(sp.label); i++ {
		if sp.label[i].h == h && sp.label[i].w == w {
			mutex.Lock()
			sp.label[i].cluster.removePixelsValue(h, w)
			mutex.Unlock()
			break
		}
	}
}

func (sp *SuperpixelsProcessorParallel) assign() {
	var mutex = &sync.Mutex{}
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	for i := 0; i < 4; i++ {
		go sp._assign_parallel(i*len(sp.clusters)/4, (i+1)*len(sp.clusters)/4, mutex, &waitGroup)
	}
	waitGroup.Wait()
}

func (sp *SuperpixelsProcessorParallel) _assign_parallel(startIdx, endIdx int, mutex *sync.Mutex, wg *sync.WaitGroup) {
	for i := startIdx; i < endIdx; i++ {
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
						mutex.Lock()
						sp.clusters[i].pixels = append(sp.clusters[i].pixels, []int{h, w})
						sp.label = append(sp.label, Label{h: h, w: w, cluster: &sp.clusters[i]})
						mutex.Unlock()
					} else {
						sp.removeLabelPixelsValue(h, w, mutex)
						mutex.Lock()
						sp.clusters[i].pixels = append(sp.clusters[i].pixels, []int{h, w})
						mutex.Unlock()
						sp.updateLabel(h, w, &sp.clusters[i], mutex)
					}
					sp.dist[h][w] = d
				}
			}
		}
	}
	wg.Done()
}

func (sp *SuperpixelsProcessorParallel) saveImage(path string, img image.Image) {
	var bounds = img.Bounds()
	var width, height = bounds.Max.X, bounds.Max.Y

	var upLeft = image.Point{0, 0}
	var lowRight = image.Point{width, height}

	var new_img = image.NewRGBA(image.Rectangle{upLeft, lowRight})

	var waitGroup sync.WaitGroup
	for i := 0; i < len(sp.clusters); i++ {
		waitGroup.Add(4)
		for j := 0; j < 4; j++ {
			go sp._saveImage_parallel(i, j*len(sp.clusters[i].pixels)/4, (j+1)*len(sp.clusters[i].pixels)/4, new_img, &waitGroup)
		}
		waitGroup.Wait()
	}

	f, _ := os.Create(path)
	png.Encode(f, new_img)
	f.Close()
}

func (sp *SuperpixelsProcessorParallel) _saveImage_parallel(cluster_idx, startIdx, endIdx int, new_img *image.RGBA, wg *sync.WaitGroup) {
	for i := startIdx; i < endIdx; i++ {
		new_img.Set(sp.clusters[cluster_idx].pixels[i][1], sp.clusters[cluster_idx].pixels[i][0], color.RGBA{uint8(sp.clusters[cluster_idx].l), uint8(sp.clusters[cluster_idx].a), uint8(sp.clusters[cluster_idx].b), 0xff})
	}
	wg.Done()
}
