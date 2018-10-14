// +build js,wasm

package main

import (
	"fmt"
	"strconv"
	"strings"
	"syscall/js"
)

func runAlgorithms(i []js.Value) {

	numNodesString := js.Global().Get("document").Call("getElementById", "numNodes").Get("value").String()
	numNodes, _ := strconv.Atoi(numNodesString)

	numGenerationsString := js.Global().Get("document").Call("getElementById", "numGenerations").Get("value").String()
	numGenerations, _ := strconv.Atoi(numGenerationsString)

	neighborsSlice := InitColoring(numNodes, numGenerations)

	oneSolution, oneFitness, onePointTime := RunOnePoint()
	twoSolution, twoFitness, twoPointTime := RunTwoPoint()
	js.Global().Get("document").Call("getElementById", "onePointTime").Set("value", onePointTime.String())
	js.Global().Get("document").Call("getElementById", "onePointFitness").Set("value", oneFitness)

	js.Global().Get("document").Call("getElementById", "twoPointTime").Set("value", twoPointTime.String())
	js.Global().Get("document").Call("getElementById", "twoPointFitness").Set("value", twoFitness)
	fmt.Println(onePointTime, twoPointTime)

	renderCanvas("canvas1", numNodes, neighborsSlice, oneSolution)
	renderCanvas("canvas2", numNodes, neighborsSlice, twoSolution)

}

func runAlgorithmsAvg(i []js.Value) {
	numNodesString := js.Global().Get("document").Call("getElementById", "numNodesAvg").Get("value").String()
	numNodes, _ := strconv.Atoi(numNodesString)

	numGenerationsString := js.Global().Get("document").Call("getElementById", "numGenerationsAvg").Get("value").String()
	numGenerations, _ := strconv.Atoi(numGenerationsString)

	iterationsString := js.Global().Get("document").Call("getElementById", "iterations").Get("value").String()
	iterations, _ := strconv.Atoi(iterationsString)

	InitColoring(numNodes, numGenerations)

	totalTimeOne := 0.0
	totalTimeTwo := 0.0

	totalFitnessOne := 0
	totalFitnessTwo := 0
	for i := 0; i < iterations; i++ {
		_, oneFitness, onePointTime := RunOnePoint()
		_, twoFitness, twoPointTime := RunTwoPoint()
		totalTimeOne += onePointTime.Seconds()
		totalTimeTwo += twoPointTime.Seconds()

		totalFitnessOne += oneFitness
		totalFitnessTwo += twoFitness
	}

	fmt.Println(totalFitnessOne, totalFitnessTwo, totalFitnessOne, totalFitnessTwo)
	js.Global().Get("document").Call("getElementById", "avgTimeOne").Set("innerHTML", totalTimeOne/float64(iterations))
	js.Global().Get("document").Call("getElementById", "avgTimeTwo").Set("innerHTML", totalTimeTwo/float64(iterations))
	js.Global().Get("document").Call("getElementById", "avgFitnessOne").Set("innerHTML", totalFitnessOne/iterations)
	js.Global().Get("document").Call("getElementById", "avgFitnessTwo").Set("innerHTML", totalFitnessTwo/iterations)

}

func main() {
	href := js.Global().Get("location").Get("href")
	fmt.Println(href)

	c := make(chan struct{}, 0)

	js.Global().Set("runAlgorithms", js.NewCallback(runAlgorithms))
	js.Global().Set("runAlgorithmsAvg", js.NewCallback(runAlgorithmsAvg))
	<-c
}

func renderCanvas(id string, numNodes int, neighborsSlice [][]int, solution []string) {

	xWidth := 120
	yWidth := 130
	dotsPerRow := 15

	canvas := js.Global().Get("document").Call("getElementById", id)
	ctx := canvas.Call("getContext", "2d")
	ctx.Call("clearRect", 0, 0, 1700, 1000)
	xPos := 0
	yPos := 0

	for i := 1; i <= numNodes; i++ {
		colorLetter := solution[i-1]
		var color string
		switch colorLetter {
		case "r":
			color = "#FF0000"
		case "b":
			color = "#000000"
		case "w":
			color = "#0000FF"
		}

		ctx.Set("fillStyle", color)
		ctx.Call("fillRect", xPos, yPos, 15, 15)

		xPos += xWidth
		if i%dotsPerRow == 0 {
			yPos += yWidth
			xPos = 0
		}
	}

	for i := 0; i < numNodes; i++ {
		neighbors := neighborsSlice[i]

		yIndexPos1 := int(i / dotsPerRow)
		yPos1 := yWidth * yIndexPos1
		xIndexPos1 := int(i % dotsPerRow)
		xPos1 := xWidth * xIndexPos1
		for _, index := range neighbors {

			yIndexPos2 := int(index / dotsPerRow)
			yPos2 := yWidth * yIndexPos2
			xIndexPos2 := int(index % dotsPerRow)

			xPos2 := xWidth * xIndexPos2

			strokeColor := fitnessLineColor(solution[i], solution[index])
			ctx.Set("strokeStyle", strokeColor)
			ctx.Call("beginPath", xPos1, yPos1)
			ctx.Call("moveTo", xPos1, yPos1)
			ctx.Call("lineTo", xPos2, yPos2)
			ctx.Call("stroke")
		}
	}
}

func fitnessLineColor(color1 string, color2 string) string {
	if strings.Compare(color1, color2) != 0 {
		return "#00FF00"
	} else {
		return "#FF0000"
	}

}
