package main

import (
	"fmt"
	"image"
	"image/draw"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"

	"github.com/anthonynsimon/bild/imgio"
)

type Particle struct {
	speed                   []float64
	position, best_position []float64
	value, best_value       float64
}

func convertToUint8(position []float64) []uint8 {
	var retval []uint8

	for _, val := range position {
		retval = append(retval, uint8(val))
	}

	return retval
}

func addVector(a []float64, b []float64) []float64 {
	var retval = make([]float64, len(a))

	for i := 0; i < len(a); i++ {
		retval[i] = a[i] + b[i]
	}
	return retval
}

func subVector(a []float64, b []float64) []float64 {
	var retval = make([]float64, len(a))

	for i := 0; i < len(a); i++ {
		retval[i] = a[i] - b[i]
	}
	return retval
}

func mulVectorConst(a []float64, b float64) []float64 {
	var retval = make([]float64, len(a))

	for i := 0; i < len(a); i++ {
		retval[i] = a[i] * b
	}
	return retval
}

func divVectorConst(a []float64, b float64) []float64 {
	var retval = make([]float64, len(a))

	for i := 0; i < len(a); i++ {
		retval[i] = a[i] / b
	}
	return retval
}

func setupProbs(thresholds []float64, pixels []uint8) []float64 {
	var N = len(pixels)
	var level_nums []float64

	for i := 0; i < len(thresholds)+1; i++ {
		level_nums = append(level_nums, 0)
	}

	//odredjujemo broj piksela u svakom opsegu
	for i := 0; i < len(pixels); i++ {
		for idx, threshold := range thresholds {
			if float64(pixels[i]) < threshold {
				level_nums[idx] += 1
				break
			} else if float64(pixels[i]) >= threshold && idx == len(thresholds)-1 {
				level_nums[idx+1] += 1
				break
			}
		}
	}
	//racunamo kolki udeo piksela svaki opseg sadrzi
	return divVectorConst(level_nums, float64(N))
}

func tsallis(thresholds []float64, pixels []uint8, order int) float64 {
	var probs = setupProbs(thresholds, pixels)
	var sum float64 = 0

	for _, prob := range probs {
		sum += math.Pow(prob, float64(order))
	}
	return (1 / float64(order-1)) * (1 - sum)
}

func get_position(thresh_num int) []float64 {
	var position []float64

	for i := 0; i < thresh_num; i++ {
		position = append(position, rand.Float64()*255)
	}

	sort.Float64s(position)

	return position
}

func writePositionLog(all_positions [][]uint8) {
	f, err := os.Create("output/all_positions.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(all_positions); i++ {
		var position string

		for j := 0; j < len(all_positions[i]); j++ {
			position += strconv.Itoa(int(all_positions[i][j]))

			if j != len(all_positions[i])-1 {
				position += ";"
			} else {
				position += "\n"
			}
		}

		_, err := f.WriteString(position)
		if err != nil {
			fmt.Println(err)
			f.Close()
			return
		}
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}

func pso_serial(pixels []uint8, n_var int, wi, wf, cpi, cpf, cgi, cgf float64, particle_num, iter_num int, tsallis_order int) []uint8 {
	var all_positions [][]uint8
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

	for i := 0; i < iter_num; i++ {
		for j := 0; j < particle_num; j++ {
			var r1 = rand.Float64()
			var r2 = rand.Float64()

			var sub1 []float64 = subVector(population[j].best_position, population[j].position)
			var sub2 []float64 = subVector(best_position, population[j].position)
			var mul1 []float64 = mulVectorConst(sub1, r1*cp)
			var mul2 []float64 = mulVectorConst(sub2, r2*cg)
			var mul3 []float64 = mulVectorConst(population[j].speed, w)
			var add1 []float64 = addVector(mul1, mul3)

			population[j].speed = addVector(add1, mul2)
			population[j].position = addVector(population[j].position, population[j].speed)
			sort.Float64s(population[j].position)

			population[j].value = tsallis(population[j].position, pixels, tsallis_order)

			all_positions = append(all_positions, convertToUint8(population[j].position))

			if population[j].value > population[j].best_value {
				population[j].best_position = population[j].position
				population[j].best_value = population[j].value
			}
			if population[j].best_value > best_value {
				best_position = population[j].best_position
				best_value = population[j].best_value
			}
		}

		w += deltaW
		cp += deltaCp
		cg -= deltaCg
	}

	//upisujemo samo ako imamo 2 ili 3 praga (samo te slucajeve iscrtavamo)
	if thresh_num < 3 {
		all_positions = append(all_positions, convertToUint8(best_position))
		writePositionLog(all_positions)
	}

	return convertToUint8(best_position)
}

func applyThresholds(img *image.Gray, thresholds []uint8) {
	//korak od 3 jer Pix sadrzi za svaki piksel redom R, G, B kanale
	for i := 0; i < len(img.Pix); i++ {
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
}

//n_var = 1, wi = 0.9, wf = 0.4, cpi = 0.5, cpf = 2.5, cgi = 2.5, cgf = 0.5, particle_num = 10, iter_num = 10, tsallis_order = 4
func main() {
	var img, err = imgio.Open("input/lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	//radimo segmentaciju samo greyscale slika
	var rect = img.Bounds()
	var img_grey = image.NewGray(rect)
	draw.Draw(img_grey, rect, img, rect.Min, draw.Src)

	var thresholds = pso_serial(img_grey.Pix, 1, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, 100, 20, 4)
	applyThresholds(img_grey, thresholds)
	fmt.Println(thresholds)

	if err := imgio.Save("output/output.png", img_grey, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}
