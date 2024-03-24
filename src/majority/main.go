package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math/rand"
	"os"
	"sort"
	"sync"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: go run main.go <size> <iteration> <threshold>")
		os.Exit(1)
	}

	size := atoi(os.Args[1])
	iterations := atoi(os.Args[2])
	threshold := atoi(os.Args[3])
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
				for j := 0; j < columns; j++ {
					currentCell := originalArray[i][j]

					neighborCounts := getNeighborCounts(i, j, originalArray)

					sameNeighborCount := neighborCounts[currentCell]
					if sameNeighborCount > threshold {
						continue
					}

					turnArray[i][j] = getMajorityNeighbor(neighborCounts)
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

func getNeighborCounts(i, j int, array [][]int) map[int]int {
	neighborCounts := make(map[int]int)
	neighborRows := []int{i - 1, i, i + 1}
	neighborColumns := []int{j - 1, j, j + 1}
	for _, row := range neighborRows {
		for _, col := range neighborColumns {
			if row < 0 || row >= len(array) || col < 0 || col >= len(array) || (row == i && col == j) {
				continue
			}
			neighborCounts[array[row][col]]++
		}
	}
	return neighborCounts
}

func getMajorityNeighbor(neighborCounts map[int]int) int {
	majorityNeighborCount := 0
	for _, count := range neighborCounts {
		if count > majorityNeighborCount {
			majorityNeighborCount = count
		}
	}

	majorityNeighbors := make([]int, 0)
	for neighbor, count := range neighborCounts {
		if count == majorityNeighborCount {
			majorityNeighbors = append(majorityNeighbors, neighbor)
		}
	}

	sort.Ints(majorityNeighbors)
	return majorityNeighbors[rand.Intn(len(majorityNeighbors))]
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
