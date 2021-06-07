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

	var img = loadImage(superpixelsDTO.InputPath)

	var bounds = img.Bounds()
	var width, height = bounds.Max.X, bounds.Max.Y
	var pixels = repackPixels(img)

	var outputLocation = superpixelsDTO.OutputPath

	if superpixelsDTO.Type == "true" {
		fmt.Println("Running parallel Superpixels...")

		var n_cpu = 4
		var sp = SuperpixelsProcessorParallel{image: pixels, img_w: width, img_h: height, K: 2000, M: 20}

		var start = time.Now()
		sp.initialize()
		sp.initClusters()
		sp.moveClusters()

		sp.assign(n_cpu)

		fmt.Println("Completed in: " + time.Since(start).String())
		sp.saveImage(outputLocation, img)
	} else {
		fmt.Println("Running serial Superpixels...")

		var sp = SuperpixelsProcessor{image: pixels, img_w: width, img_h: height, K: 2000, M: 20}

		var start = time.Now()
		sp.initialize()
		sp.initClusters()
		sp.moveClusters()

		for i := 0; i < 10; i++ {
			sp.assign()
			sp.updateCluster()
		}
		fmt.Println("Completed in: " + time.Since(start).String())
		sp.saveImage(outputLocation, img)
	}
	json.NewEncoder(w).Encode(outputLocation)
}

func main() {
	http.HandleFunc("/superpixels", applySuperPixels)
	log.Fatal(http.ListenAndServe(":8081", nil))
}
