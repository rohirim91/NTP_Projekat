package main

import (
	"image"
	"math"
	"sort"
	"sync"
)

func psoParallel(pixels []uint8, thresh_num int, wi, wf, cpi, cpf, cgi, cgf float64, particle_num, iter_num int, tsallis_order int) ([]uint8, [][]uint8, []float64) {
	var all_positions = make([][]uint8, iter_num*particle_num+1)
	var all_values = make([]float64, iter_num*particle_num+1)
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
	//fmt.Println(population)
	var mutex = &sync.Mutex{}
	var wg sync.WaitGroup
	var n_cpu = 4
	wg.Add(n_cpu)
	for j := 0; j < n_cpu-1; j++ {
		go psoInnerLoop(iter_num, deltaW, deltaCp, deltaCg, j*len(population)/n_cpu, (j+1)*len(population)/n_cpu, particle_num, tsallis_order, w, cp, cg, population, pixels, all_positions, all_values, &best_position, &best_value, mutex, &wg)
	}
	psoInnerLoop(iter_num, deltaW, deltaCp, deltaCg, 3*len(population)/n_cpu, (3+1)*len(population)/n_cpu, particle_num, tsallis_order, w, cp, cg, population, pixels, all_positions, all_values, &best_position, &best_value, mutex, &wg)
	wg.Wait()

	//upisujemo samo ako imamo 2 ili 3 praga (samo te slucajeve iscrtavamo)
	if thresh_num < 3 {
		all_positions[iter_num*particle_num] = convertToUint8(best_position)
		all_values[iter_num*particle_num] = best_value
	}

	return convertToUint8(best_position), all_positions, all_values
}

func psoInnerLoop(iter_num int, deltaW, deltaCp, deltaCg float64, startIdx, endIdx, particle_num, tsallis_order int, w, cp, cg float64, population []Particle, pixels []uint8, all_positions [][]uint8, all_values []float64, best_position *[]float64, best_value *float64, mutex *sync.Mutex, wg *sync.WaitGroup) {
	for iter_no := 0; iter_no < iter_num; iter_no++ {
		for j := startIdx; j < endIdx; j++ {
			var sub1 []float64 = subVector(population[j].best_position, population[j].position)

			var sub2 []float64 = subVector(*best_position, population[j].position)
			var mul1 []float64 = mulVectorConst(sub1, cp/2)
			var mul2 []float64 = mulVectorConst(sub2, cg/2)
			var mul3 []float64 = mulVectorConst(population[j].speed, w)
			var add1 []float64 = addVector(mul1, mul3)

			population[j].speed = addVector(add1, mul2)
			population[j].position = addVector(population[j].position, population[j].speed)

			population[j].position = checkPositionInvalid(&population[j].position)

			sort.Float64s(population[j].position)

			population[j].value = tsallis(population[j].position, pixels, tsallis_order)

			all_positions[iter_no*particle_num+j] = convertToUint8(population[j].position)
			all_values[iter_no*particle_num+j] = population[j].value

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
		w += deltaW
		cp += deltaCp
		cg -= deltaCg
	}
	wg.Done()
}

func _applyThresholdsParallel(startIdx, endIdx int, img *image.Gray, thresholds []uint8, wg *sync.WaitGroup) {
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

func applyThresholdsParallel(img_grey *image.Gray, thresholds []uint8) {
	var waitGroup sync.WaitGroup
	waitGroup.Add(4)
	for i := 0; i < 4; i++ {
		go _applyThresholdsParallel(i*len(img_grey.Pix)/4, (i+1)*len(img_grey.Pix)/4, img_grey, thresholds, &waitGroup)
	}
	waitGroup.Wait()
}
