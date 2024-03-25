package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"math/rand"
	"os"
	"sync"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run main.go <size> <iteration>")
		os.Exit(1)
	}

	size := atoi(os.Args[1])
	iterations := atoi(os.Args[2])
	threshold := math.Log2(math.Sqrt(float64(size * iterations)))
	rows := size
	columns := size

	originalArray := make([][]int, rows)
	for i := range originalArray {
		originalArray[i] = make([]int, columns)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < columns; j++ {
			originalArray[i][j] = rand.Intn(4) + 1
		}
	}

	turnArray := make([][]int, rows)
	for i := range turnArray {
		turnArray[i] = make([]int, columns)
		copy(turnArray[i], originalArray[i])
	}

	for iter := 0; iter <= iterations; iter++ {
		saveImage(iter, originalArray)
		var wg sync.WaitGroup
		for i := 0; i < rows; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				for j := 0; j < columns; j++ {
					currentCell := originalArray[i][j]

					neighborsWeights := getNeighborWeights(i, j, threshold, originalArray)

					sameNeighborWeight := neighborsWeights[currentCell]
					if math.Sqrt(sameNeighborWeight) > threshold {
						continue
					}

					turnArray[i][j] = pickPixelValue(neighborsWeights)
				}
			}(i)
		}
		copy(originalArray, turnArray)
	}
}

func atoi(s string) int {
	i := 0
	for _, r := range s {
		i = i*10 + int(r-'0')
	}
	return i
}

func getNeighborWeights(i, j int, threshold float64, array [][]int) map[int]float64 {
	neighborWeights := make(map[int]float64)
	neighborRows := []int{i - 2, i - 1, i, i + 1, i + 2}
	neighborColumns := []int{j - 2, j - 1, j, j + 1, j + 2}
	for _, row := range neighborRows {
		for _, col := range neighborColumns {
			if row < 0 || row >= len(array) || col < 0 || col >= len(array) || (row == i && col == j) {
				continue
			}

			weight := threshold / 10
			if row == i-2 || row == i+2 || col == j-2 || col == j+2 {
				neighborWeights[array[row][col]] += weight
			} else {
				neighborWeights[array[row][col]] += weight * 1.67
			}
		}
	}
	return neighborWeights
}

func pickPixelValue(neighborWeights map[int]float64) int {
	var biggestWeight float64 = 0
	index := 0
	for i, weight := range neighborWeights {
		if weight > biggestWeight {
			biggestWeight = weight
			index = i
		}
	}
	return index
}

func saveImage(iter int, array [][]int) {
	img := image.NewRGBA(image.Rect(0, 0, len(array), len(array[0])))
	for i := range array {
		for j := range array[i] {
			value := array[i][j]
			var pixel color.RGBA
			switch value {
			case 1:
				pixel = color.RGBA{155, 118, 83, 255}
			case 2:
				pixel = color.RGBA{194, 178, 128, 255}
			case 3:
				pixel = color.RGBA{27, 141, 58, 255}
			case 4:
				pixel = color.RGBA{0, 0, 255, 255}
			}
			img.Set(i, j, pixel)
		}
	}

	f, err := os.Create(fmt.Sprintf("%03d.png", iter))
	if err != nil {
		fmt.Println("Error creating image file:", err)
		os.Exit(1)
	}
	defer f.Close()

	err = png.Encode(f, img)
	if err != nil {
		fmt.Println("Error encoding image:", err)
		os.Exit(1)
	}
}
