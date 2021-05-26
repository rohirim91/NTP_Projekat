package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func applySuperPixels(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var superpixelsDTO SuperpixelsDTO
	json.Unmarshal(reqBody, &superpixelsDTO)

	var img = loadImage("../input/lena.png")

	var bounds = img.Bounds()
	var width, height = bounds.Max.X, bounds.Max.Y
	var pixels = repackPixels(img)

	if superpixelsDTO.Type == "parallel" {
		var sp = SuperpixelsProcessorParallel{image: pixels, img_w: width, img_h: height, K: 2000, M: 20}
		sp.initialize()
		sp.initClusters()
		sp.moveClusters()

		var start = time.Now()
		for i := 0; i < 10; i++ {
			sp.assign()
			sp.updateCluster()
		}
		sp.saveImage("../output/output.png", img)
		fmt.Println("Parallel: " + time.Since(start).String() + " - ")
	} else {
		var sp = SuperpixelsProcessor{image: pixels, img_w: width, img_h: height, K: 2000, M: 20}
		sp.initialize()
		sp.initClusters()
		sp.moveClusters()

		var start = time.Now()
		for i := 0; i < 10; i++ {
			sp.assign()
			sp.updateCluster()
		}
		sp.saveImage("../output/output.png", img)
		fmt.Println("Serial: " + time.Since(start).String() + " - ")
	}
}

func main() {
	http.HandleFunc("/superpixels", applySuperPixels)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
