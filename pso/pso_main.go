package main

import (
	"encoding/json"
	"fmt"
	"image"
	"image/draw"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/anthonynsimon/bild/imgio"
)

func runPso(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var psoDTO PsoDTO
	json.Unmarshal(reqBody, &psoDTO)

	var img, err = imgio.Open("../input/lena.png")
	if err != nil {
		fmt.Println(err)
		return
	}

	var rect = img.Bounds()
	var img_grey = image.NewGray(rect)
	draw.Draw(img_grey, rect, img, rect.Min, draw.Src)

	if psoDTO.Type == "parallel" {
		var start = time.Now()
		var thresholds = psoParallel(img_grey.Pix, 1, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, 100, 20, 4)
		applyThresholdsParallel(img_grey, thresholds)
		fmt.Print("Parallel: " + time.Since(start).String() + " - ")
	} else {
		var start = time.Now()
		var thresholds = psoSerial(img_grey.Pix, 1, 0.9, 0.4, 0.5, 2.5, 2.5, 0.5, 100, 20, 4)
		applyThresholds(img_grey, thresholds)
		fmt.Print("Serial: " + time.Since(start).String() + " - ")
	}

	if err := imgio.Save("../output/output.png", img_grey, imgio.PNGEncoder()); err != nil {
		fmt.Println(err)
		return
	}
}

//n_var = 1, wi = 0.9, wf = 0.4, cpi = 0.5, cpf = 2.5, cgi = 2.5, cgf = 0.5, particle_num = 10, iter_num = 10, tsallis_order = 4
func main() {
	http.HandleFunc("/pso", runPso)
	log.Fatal(http.ListenAndServe(":8082", nil))
}
