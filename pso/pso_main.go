package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

func runPso(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var psoDTO PsoDTO
	json.Unmarshal(reqBody, &psoDTO)

	var img, err = imgio.Open(psoDTO.InputPath)
	if err != nil {
		fmt.Println(err)
		return
	}

	var rect = img.Bounds()
	var img_grey = image.NewGray(rect)
	draw.Draw(img_grey, rect, img, rect.Min, draw.Src)

	var outputLocation = psoDTO.OutputPath

	if psoDTO.Type == "true" {
		for k := 0; k < 1; k++ {
			rand.Seed(42)

			// var img_grey = image.NewGray(rect)
			// draw.Draw(img_grey, rect, img, rect.Min, draw.Src)

			var posSaveLocation = "../output/all_positions" + fmt.Sprint(k) + ".csv"

			fmt.Println("Running parallel PSO...")

			var start = time.Now()
			var thresholds, all_positions, all_values = psoParallel(img_grey.Pix, psoDTO.NumThresholds, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, psoDTO.MaxIter, psoDTO.NumParticles, 4)
			fmt.Println("Completed in: " + time.Since(start).String())

			writePositionLog(all_positions, all_values, posSaveLocation)
			applyThresholdsParallel(img_grey, thresholds)
		}
	} else {
		const posSaveLocation = "../output/all_positions.csv"

		fmt.Println("Running serial PSO...")

		var start = time.Now()
		var thresholds, all_positions, all_values = psoSerial(img_grey.Pix, psoDTO.NumThresholds, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, psoDTO.MaxIter, psoDTO.NumParticles, 4)
		fmt.Println("Completed in: " + time.Since(start).String())

		writePositionLog(all_positions, all_values, posSaveLocation)
		applyThresholds(img_grey, thresholds)
	}

	if err := imgio.Save(outputLocation, img_grey, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode("file:///D:/Fakultet/4.%20godina/NTP/Projekat/NTP_Projekat/output/all_positions.csv")
}

//thresh_num = 1, wi = 0.9, wf = 0.4, cpi = 0.5, cpf = 2.5, cgi = 2.5, cgf = 0.5, particle_num = 10, iter_num = 10, tsallis_order = 4
func main() {
	http.HandleFunc("/pso", runPso)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
