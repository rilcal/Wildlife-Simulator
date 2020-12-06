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
	w = newWorld(x-1, y-1)
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
				w.Tiles[x][y] = getTileType("Water")
				wp := newPoint(x, y)
				w.WaterTile = append(w.WaterTile, wp)
				waterCount++
			} else if terrain > -0.12 && terrain < 0.3 {
				w.Tiles[x][y] = getTileType("Land")
				lp := newPoint(x, y)
				w.LandTile = append(w.LandTile, lp)
				landCount++
			} else {
				w.Tiles[x][y] = getTileType("Mountain")
				lp := newPoint(x, y)
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
		a.Sheeps[i] = newAnimal("Sheep", i)
		w.Tiles[a.Sheeps[i].Pos.X][a.Sheeps[i].Pos.Y].HasAnimal = true
		w.Tiles[a.Sheeps[i].Pos.X][a.Sheeps[i].Pos.Y].AnimalType = a.Sheeps[i]
	}

	//Generate wolf population
	a.Wolves = make([]structs.Animal, numWolves)
	for j := 0; j < numWolves; j++ {
		a.Wolves[j] = newAnimal("Wolf", j)
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
	moveToPosition := averagePoints(pointsOfSheepInHerd)
	w.MoveAnimal(*sheep, moveToPosition)
	sheep.Move(moveToPosition)
	return
}

func getTileType(tileDesc string) (t structs.Tile) {
	// defining the description, symbols, and style of all tiles
	if tileDesc == "Water" { // Terrain
		t = newTile("Water", '~', getSetStyles("Water"), false)
	} else if tileDesc == "Land" {
		t = newTile("Land", '#', getSetStyles("Land"), false)
	} else if tileDesc == "Mountain" {
		t = newTile("Mountain", 'M', getSetStyles("Mountain"), false)
	} else {
		t = newTile("", '!', getSetStyles(""), false)
	}
	return
}

//Default styles. "Graphics" for everything to be printed to the screen
func getSetStyles(tileDesc string) (s tcell.Style) {
	// defining the colors and character styles for all terrains and structs.Animals
	var color tcell.Color
	var bold bool
	if tileDesc == "Water" { // Terrain
		color = tcell.NewRGBColor(0, 153, 153)
	} else if tileDesc == "Land" {
		color = tcell.NewRGBColor(0, 153, 76)
	} else if tileDesc == "Mountain" {
		color = tcell.NewRGBColor(160, 160, 160)
	} else if tileDesc == "Wolf" { // Animals
		color = tcell.NewRGBColor(255, 51, 51)
		bold = true
	} else if tileDesc == "Sheep" {
		color = tcell.NewRGBColor(255, 255, 255)
		bold = true
	} else {
		color = tcell.NewRGBColor(255, 0, 0)
	}
	s = s.Foreground(color)
	s = s.Bold(bold)
	return
}

func averagePoints(s []structs.Point) (p structs.Point) {
	var xvalues int = 0
	var yvalues int = 0
	var count int = 0
	for i := range s {
		xvalues += s[i].X
		yvalues += s[i].Y
		count++
	}
	p.X = xvalues / count
	p.Y = yvalues / count
	return
}

// Factory Functions
func newWorld(x, y int) (w structs.World) {
	w.Width = x
	w.Length = y
	w.Tiles = make([][]structs.Tile, x)
	w.WaterTile = make([]structs.Point, 0)
	w.LandTile = make([]structs.Point, 0)
	return
}

func newTile(desc string, sym rune, style tcell.Style, occ bool) structs.Tile {
	var t structs.Tile
	t.TerrainDesc = desc
	t.TerrainSym = sym
	t.TerrainStyle = style
	t.HasAnimal = occ
	return t
}

func newPoint(x, y int) (p structs.Point) {
	p.X = x
	p.Y = y
	return
}

func newAnimal(desc string, index int) (a structs.Animal) {
	if desc == "Sheep" {
		//set initial state
		a.Desc = "Sheep"
		a.Sym = 'S'
		a.Sty = getSetStyles("Sheep")
		a.Key = index

		//set spawn Point
		r := rand.Intn(len(w.LandTile))
		a.Pos = w.LandTile[r]

	} else if desc == "Wolf" {
		//set initial state
		a.Desc = "Wolf"
		a.Sym = 'W'
		a.Sty = getSetStyles("Wolf")
		a.Key = index

		//set spawn Point
		r := rand.Intn(len(w.LandTile))
		a.Pos = w.LandTile[r]

	} else {
		panic("Not a valid Animal")
	}
	return
}
