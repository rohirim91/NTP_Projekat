package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"sort"
	"strconv"
)

type PsoDTO struct {
	InputPath     string `json:"InputPath"`
	Type          string `json:"Type"`
	NumParticles  int    `json:"NumParticles"`
	MaxIter       int    `json:"MaxIter"`
	NumThresholds int    `json:"NumThresholds"`
}

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

func checkPositionInvalid(position *[]float64) []float64 {
	for i := 0; i < len(*position); i++ {
		if !((*position)[i] > 0 && (*position)[i] < 256) {
			(*position)[i] = rand.Float64() * 255
		}
	}
	return *position
}

func checkWritten(position *[]uint8, writtenPos *[][]uint8) bool {
	if len(*writtenPos) == 0 {
		return false
	}
	for i := 0; i < len(*writtenPos); i++ {
		var found = true
		for j := 0; j < len(*position); j++ {
			if (*writtenPos)[i][j] != (*position)[j] {
				found = false
			}
		}
		if found {
			return true
		}
	}
	return false
}

func writePositionLog(all_positions [][]uint8, all_values []float64, posSaveLocation string) {
	f, err := os.Create(posSaveLocation)
	if err != nil {
		fmt.Println(err)
		return
	}

	var header string
	for i := 0; i < len(all_positions[0]); i++ {
		header += "t" + strconv.Itoa(i+1) + ","
	}

	header += "v,s\n"

	_, err = f.WriteString(header)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}

	var writtenPos [][]uint8
	for i := 0; i < len(all_positions); i++ {
		if checkWritten(&(all_positions[i]), &writtenPos) && (i != len(all_positions)-1) {
			continue
		}
		writtenPos = append(writtenPos, all_positions[i])

		var position string

		for j := 0; j < len(all_positions[i]); j++ {
			position += strconv.Itoa(int(all_positions[i][j]))
			position += ","
		}

		position += strconv.FormatFloat(all_values[i], 'f', 5, 64) + ","

		if i == len(all_positions)-1 {
			position += "1\n"
		} else {
			position += "0.1\n"
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
