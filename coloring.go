// A basic HTTP server.
// By default, it serves the current working directory on port 8080.
package main

import (
	"math/rand"
	"sort"
	"strings"
	"time"
)

type slice struct {
	sort.IntSlice
	solutions [][]string
}

func (s slice) Swap(i, j int) {
	s.IntSlice.Swap(i, j)
	s.solutions[i], s.solutions[j] = s.solutions[j], s.solutions[i]
}

func contains(s []int, e int) (bool, int) {
	for i, a := range s {
		if a == e {
			return true, i
		}
	}
	return false, -1
}

func generateRandomSolution(length int, colors []string) []string {
	states := []string{}
	for i := 0; i < length; i++ {
		randIndex := rand.Intn(len(colors))
		color := colors[randIndex]
		states = append(states, color)
	}

	return states
}

func fitnessFunction(neighbors [][]int, solution []string) int {

	fitness := 0
	for i := 0; i < len(neighbors); i++ {
		for j := 0; j < len(neighbors[i]); j++ {
			index := neighbors[i][j]
			if strings.Compare(solution[i], solution[index]) == 0 {
				fitness++
			}
		}
	}

	return fitness
}

func getBestHalf(solutions [][]string, neighbors [][]int) [][]string {
	length := len(solutions)

	fitnessArray := make([]int, 0, length)

	for i := 0; i < length; i++ {
		fitnessArray = append(fitnessArray, fitnessFunction(neighbors, solutions[i]))
	}

	s := &slice{IntSlice: sort.IntSlice(fitnessArray), solutions: solutions}
	sort.Sort(s)

	return s.solutions[:length/2]
}

func getBestSolution(solutions [][]string, neighbors [][]int) ([]string, int) {
	best := 99999999
	index := 0
	for i := 0; i < len(solutions); i++ {
		res := fitnessFunction(neighbors, solutions[i])
		if res < best {
			best = res
			index = i
		}
	}
	return solutions[index], best
}
func onePointCrossover(solutions [][]string, numNodes int, numGenerations int, neighbors [][]int, colors []string) ([][]string, time.Duration) {
	start := time.Now()

	mutationProbability := 0.1
	newSolutions := make([][]string, len(solutions))

	for i := range newSolutions {
		newSolutions[i] = make([]string, len(solutions[i]))
		copy(newSolutions[i], solutions[i])
	}

	for i := 0; i < numGenerations; i++ {
		solutionsLength := len(newSolutions)
		for j := 0; j < solutionsLength; j += 2 {
			index := rand.Intn(numNodes)
			parent1, parent2 := newSolutions[j], newSolutions[j+1]
			child1, child2 := make([]string, numNodes), make([]string, numNodes)
			copy(child1, parent1)
			copy(child2, parent2)
			child1[index], child2[index] = parent2[index], parent1[index]

			child1Mutation := rand.Float64()
			if child1Mutation < mutationProbability {
				randIndex := rand.Intn(numNodes)
				randColorIndex := rand.Intn(len(colors))
				child1[randIndex] = colors[randColorIndex]
			}

			child2Mutation := rand.Float64()
			if child2Mutation < mutationProbability {
				randIndex := rand.Intn(numNodes)
				randColorIndex := rand.Intn(len(colors))
				child2[randIndex] = colors[randColorIndex]
			}

			newSolutions = append(newSolutions, child1)
			newSolutions = append(newSolutions, child2)
		}

		newSolutions = getBestHalf(newSolutions, neighbors)
	}
	return newSolutions, time.Since(start)

}

func twoPointCrossover(solutions [][]string, numNodes int, numGenerations int, neighbors [][]int, colors []string) ([][]string, time.Duration) {
	start := time.Now()
	mutationProbability := 0.1

	newSolutions := make([][]string, len(solutions))
	for i := range newSolutions {
		newSolutions[i] = make([]string, len(solutions[i]))
		copy(newSolutions[i], solutions[i])
	}

	for i := 0; i < numGenerations; i++ {
		solutionsLength := len(newSolutions)
		for j := 0; j < solutionsLength; j += 2 {
			fromIndex := rand.Intn(numNodes)
			toIndex := fromIndex + rand.Intn(numNodes-fromIndex)
			parent1, parent2 := newSolutions[j], newSolutions[j+1]
			child1, child2 := make([]string, numNodes), make([]string, numNodes)
			copy(child1, parent1)
			copy(child2, parent2)

			for index := fromIndex; index < toIndex; index++ {
				child1[index], child2[index] = parent2[index], parent1[index]
			}

			child1Mutation := rand.Float64()
			if child1Mutation < mutationProbability {
				randIndex := rand.Intn(numNodes)
				randColorIndex := rand.Intn(len(colors))
				child1[randIndex] = colors[randColorIndex]
			}

			child2Mutation := rand.Float64()
			if child2Mutation < mutationProbability {
				randIndex := rand.Intn(numNodes)
				randColorIndex := rand.Intn(len(colors))
				child2[randIndex] = colors[randColorIndex]
			}

			newSolutions = append(newSolutions, child1)
			newSolutions = append(newSolutions, child2)
		}

		newSolutions = getBestHalf(newSolutions, neighbors)
	}

	return newSolutions, time.Since(start)
}

func getRandomUniqueNeighborSlice(index int, numNodes int) []int {
	neighbors := []int{}

	for i := 0; i < numNodes; i++ {
		rand := rand.Float64()
		if rand < 0.05 && i != index {
			neighbors = append(neighbors, i)
		}
	}

	return neighbors
}

func generateRandomGraph(numNodes int) [][]int {
	neighbors := make([][]int, numNodes)
	for i := 0; i < numNodes; i++ {
		neighbors[i] = getRandomUniqueNeighborSlice(i, numNodes)
	}

	return neighbors
}

var onePointSolutions = [][]string{}
var twoPointSolutions = [][]string{}
var randSolutions [][]string
var neighbors [][]int
var numNodes = 100
var numGenerations = 20
var colors = []string{"b", "w", "r"}

// Init the algorithm
func InitColoring(nodes int, generations int) [][]int {
	numNodes = nodes
	numGenerations = generations
	rand.Seed(time.Now().UTC().UnixNano())
	neighbors = generateRandomGraph(numNodes)
	randSolutions = [][]string{}
	for i := 0; i < numNodes; i++ {
		randSolutions = append(randSolutions, generateRandomSolution(numNodes, colors))
	}
	return neighbors
}

// RUN ONE POINT
func RunOnePoint() ([]string, int, time.Duration) {
	solutions, time := onePointCrossover(randSolutions, numNodes, numGenerations, neighbors, colors)
	solution, fitness := getBestSolution(solutions, neighbors)
	return solution, fitness, time
}

// RUN TWO POINT
func RunTwoPoint() ([]string, int, time.Duration) {
	solutions, time := twoPointCrossover(randSolutions, numNodes, numGenerations, neighbors, colors)
	solution, fitness := getBestSolution(solutions, neighbors)

	return solution, fitness, time
}
