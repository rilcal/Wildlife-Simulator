package main

import (
	"math/rand"
	"sync"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell"
	"github.com/rilcal/Wildlife-Simulator/structs"
)

//Global Variables
var w structs.World
var s tcell.Screen
var a structs.Animals

// MAIN FUNCTION
func main() {
	generateInitialVariables()
	populateWorld(25, 5)
	defer s.Fini()
	rand.Seed(time.Now().UnixNano())
	updateScreen()
	time.Sleep(time.Second * 5)
	mainLoop()
}

// Functions
func generateInitialVariables() {
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	s.Init()
	x, y := s.Size()
	w = structs.NewWorld(x-1, y-1)
	generateTerrain()
	return
}

//generates initial terrain of the structs.World
func generateTerrain() {
	rand.Seed(time.Now().UnixNano())
	p := perlin.NewPerlin(2, 2, 10, int64(rand.Int()))
	w.Tiles = make([][]structs.Tile, w.Width)
	// Initializing the tiles
	for i := range w.Tiles {
		w.Tiles[i] = make([]structs.Tile, w.Length)
	}

	waterCount := 0
	landCount := 0
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Length; y++ {
			terrain := p.Noise2D(float64(x)/10, float64(y)/10)
			if terrain <= -0.12 {
				w.Tiles[x][y] = structs.GetTileType("Water")
				wp := structs.NewPoint(x, y)
				w.WaterTile = append(w.WaterTile, wp)
				waterCount++
			} else if terrain > -0.12 && terrain < 0.3 {
				w.Tiles[x][y] = structs.GetTileType("Land")
				lp := structs.NewPoint(x, y)
				w.LandTile = append(w.LandTile, lp)
				landCount++
			} else {
				w.Tiles[x][y] = structs.GetTileType("Mountain")
				lp := structs.NewPoint(x, y)
				w.LandTile = append(w.LandTile, lp)
				landCount++
			}
		}
	}
}

func populateWorld(numSheep int, numWolves int) structs.Animals {
	var a structs.Animals

	//Generate sheep population
	a.Sheeps = make([]structs.Animal, numSheep)
	for i := 0; i < numSheep; i++ {
		a.Sheeps[i] = structs.NewAnimal("Sheep", i)
		w.Tiles[a.Sheeps[i].Pos.X][a.Sheeps[i].Pos.Y].HasAnimal = true
		w.Tiles[a.Sheeps[i].Pos.X][a.Sheeps[i].Pos.Y].AnimalType = a.Sheeps[i]
	}

	//Generate wolf population
	a.Wolves = make([]structs.Animal, numWolves)
	for j := 0; j < numWolves; j++ {
		a.Wolves[j] = structs.NewAnimal("Wolf", j)
		w.Tiles[a.Wolves[j].Pos.X][a.Wolves[j].Pos.Y].HasAnimal = true
		w.Tiles[a.Wolves[j].Pos.X][a.Wolves[j].Pos.Y].AnimalType = a.Wolves[j]
	}
	return a
}

func updateScreen() {
	for i, ii := range w.Tiles {
		for j := range ii {
			if ii[j].HasAnimal {
				(s).SetContent(i, j, rune(ii[j].AnimalType.Sym), []rune(""), ii[j].AnimalType.Sty)
			} else {
				(s).SetContent(i, j, rune(ii[j].TerrainSym), []rune(""), ii[j].TerrainStyle)

			}
		}
	}
	(s).Show()
}

func mainLoop() {
	for t := 0; t < 5; t++ {
		// Sheep Logic
		var wg1 sync.WaitGroup
		for i := range a.Sheeps {
			sheepLogic(&a.Sheeps[i], &wg1)
			wg1.Wait()
		}
		updateScreen()
		time.Sleep(time.Second * 1)
	}
}

func sheepLogic(sheep *structs.Animal, wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	// herding
	pointsOfSheepInHerd := make([]structs.Point, 0)
	for i := range a.Sheeps {
		if a.Sheeps[i].Key == sheep.Key {
			continue
		}
		pointsOfSheepInHerd = append(pointsOfSheepInHerd, a.Sheeps[i].Pos)
	}
	moveToPosition := structs.AveragePoints(pointsOfSheepInHerd)
	w.MoveAnimal(*sheep, moveToPosition)
	sheep.Move(moveToPosition)
	return
}
