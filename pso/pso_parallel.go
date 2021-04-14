package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

func psoParallel(pixels []uint8, n_var int, wi, wf, cpi, cpf, cgi, cgf float64, particle_num, iter_num int, tsallis_order int) []uint8 {
	var all_positions = make([][]uint8, iter_num*particle_num+1)
	var thresh_num = n_var + 1
	var best_position []float64
	var best_value = math.Inf(-1)

	var deltaW = math.Abs(wi-wf) / float64(iter_num)
	var deltaCp = math.Abs(cpi-cpf) / float64(iter_num)
	var deltaCg = math.Abs(cgi-cgf) / float64(iter_num)

	var w = wi
	var cp = cpi
	var cg = cgi

	var population []Particle

	for i := 0; i < particle_num; i++ {
		var pos = get_position(thresh_num)
		var fpoz = tsallis(pos, pixels, tsallis_order)

		if fpoz > best_value {
			best_position = pos
			best_value = fpoz
		}

		population = append(population, Particle{speed: make([]float64, thresh_num), position: pos, best_position: pos, value: fpoz, best_value: fpoz})
	}

	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	for i := 0; i < iter_num; i++ {
		wg.Add(4)
		for j := 0; j < 4; j++ {
			go psoInnerLoop(i, j*len(population)/4, (j+1)*len(population)/4, particle_num, tsallis_order, w, cp, cg, population, pixels, all_positions, &best_position, &best_value, mutex, &wg)
		}
		wg.Wait()
		w += deltaW
		cp += deltaCp
		cg -= deltaCg
	}

	//upisujemo samo ako imamo 2 ili 3 praga (samo te slucajeve iscrtavamo)
	if thresh_num < 3 {
		all_positions[iter_num*particle_num-1] = convertToUint8(best_position)
		writePositionLog(all_positions)
	}

	return convertToUint8(best_position)
}

func psoInnerLoop(iter_no, startIdx, endIdx, particle_num, tsallis_order int, w, cp, cg float64, population []Particle, pixels []uint8, all_positions [][]uint8, best_position *[]float64, best_value *float64, mutex *sync.Mutex, wg *sync.WaitGroup) {
	for j := startIdx; j < endIdx; j++ {
		var r1 = rand.Float64()
		var r2 = rand.Float64()

		var sub1 []float64 = subVector(population[j].best_position, population[j].position)
		var sub2 []float64 = subVector(*best_position, population[j].position)
		var mul1 []float64 = mulVectorConst(sub1, r1*cp)
		var mul2 []float64 = mulVectorConst(sub2, r2*cg)
		var mul3 []float64 = mulVectorConst(population[j].speed, w)
		var add1 []float64 = addVector(mul1, mul3)

		population[j].speed = addVector(add1, mul2)
		population[j].position = addVector(population[j].position, population[j].speed)
		sort.Float64s(population[j].position)

		population[j].value = tsallis(population[j].position, pixels, tsallis_order)

		all_positions[iter_no*particle_num+j] = convertToUint8(population[j].position)

		if population[j].value > population[j].best_value {
			population[j].best_position = population[j].position
			population[j].best_value = population[j].value

			mutex.Lock()
			if population[j].best_value > *best_value {
				*best_position = population[j].best_position
				*best_value = population[j].best_value
			}
			mutex.Unlock()
		}
	}
	wg.Done()
}

func applyThresholds_parallel(startIdx, endIdx int, img *image.Gray, thresholds []uint8, wg *sync.WaitGroup) {
	//korak od 3 jer Pix sadrzi za svaki piksel redom R, G, B kanale
	for i := startIdx; i < endIdx; i++ {
		for j := 0; j < len(thresholds); j++ {
			if img.Pix[i] < thresholds[j] {
				if j == 0 {
					img.Pix[i] = 0
				} else {
					img.Pix[i] = uint8((uint16(thresholds[j]) + uint16(thresholds[j-1])) / 2)
				}
				break
			}
		}
		if img.Pix[i] >= thresholds[len(thresholds)-1] {
			img.Pix[i] = 255
		}
	}
	wg.Done()
}

//n_var = 1, wi = 0.9, wf = 0.4, cpi = 0.5, cpf = 2.5, cgi = 2.5, cgf = 0.5, particle_num = 10, iter_num = 10, tsallis_order = 4
func main() {
	var img, err = imgio.Open("../input/lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	//radimo segmentaciju samo greyscale slika
	var rect = img.Bounds()
	var img_grey = image.NewGray(rect)
	draw.Draw(img_grey, rect, img, rect.Min, draw.Src)

	var start = time.Now()
	var thresholds = psoSerial(img_grey.Pix, 1, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, 100, 20, 4)
	fmt.Print("Serial: " + time.Since(start).String() + " - ")
	fmt.Println(thresholds)

	start = time.Now()
	thresholds = psoParallel(img_grey.Pix, 1, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, 100, 20, 4)
	fmt.Print("Parallel: " + time.Since(start).String() + " - ")
	fmt.Println(thresholds)

	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	for i := 0; i < 4; i++ {
		go applyThresholds_parallel(i*len(img_grey.Pix)/4, (i+1)*len(img_grey.Pix)/4, img_grey, thresholds, &waitGroup)
	}
	waitGroup.Wait()

	if err := imgio.Save("../output/output.png", img_grey, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}
