package main

import (
	"fmt"
	"time"
)

func main() {
	for j := 0; j < 10; j++ {
		var img = loadImage("../input/lena.png")

		var bounds = img.Bounds()
		var width, height = bounds.Max.X, bounds.Max.Y
		var pixels = repackPixels(img)
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
	}
}
