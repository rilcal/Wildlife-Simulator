package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell"
)

//Global Variables
var w, s = generateInitialVariables()
var a = populateWorld(25, 5)

// MAIN FUNCTION
func main() {
	defer s.Fini()
	rand.Seed(time.Now().UnixNano())
	updateScreen()
	time.Sleep(time.Second * 5)
	mainLoop()
}

// Funtions
func generateInitialVariables() (w world, s tcell.Screen) {
	s, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	s.Init()
	x, y := s.Size()
	w = newWorld(x-1, y-1)
	w.generateTerrain()
	return
}

func populateWorld(numSheep int, numWolves int) animals {
	var a animals

	//Generate sheep population
	a.sheeps = make([]animal, numSheep)
	for i := 0; i < numSheep; i++ {
		a.sheeps[i] = newAnimal("Sheep", i)
		w.tiles[a.sheeps[i].pos.x][a.sheeps[i].pos.y].hasAnimal = true
		w.tiles[a.sheeps[i].pos.x][a.sheeps[i].pos.y].animal = a.sheeps[i]
	}

	//Generate wolf population
	a.wolves = make([]animal, numWolves)
	for j := 0; j < numWolves; j++ {
		a.wolves[j] = newAnimal("Wolf", j)
		w.tiles[a.wolves[j].pos.x][a.wolves[j].pos.y].hasAnimal = true
		w.tiles[a.wolves[j].pos.x][a.wolves[j].pos.y].animal = a.wolves[j]
	}
	return a
}

func updateScreen() {
	for i, ii := range w.tiles {
		for j := range ii {
			if ii[j].hasAnimal {
				(s).SetContent(i, j, rune(ii[j].animal.sym), []rune(""), ii[j].animal.sty)
			} else {
				(s).SetContent(i, j, rune(ii[j].terrainSym), []rune(""), ii[j].terrainStyle)

			}
		}
	}
	(s).Show()
}

func mainLoop() {
	for t := 0; t < 5; t++ {
		// Sheep Logic
		var wg1 sync.WaitGroup
		for i := range a.sheeps {
			a.sheeps[i].sheepLogic(&wg1)
			wg1.Wait()
		}
		updateScreen()
		time.Sleep(time.Second * 1)
	}
}

func (sheep *animal) sheepLogic(wg *sync.WaitGroup) {
	defer wg.Done()
	wg.Add(1)

	// herding
	pointsOfSheepInHerd := make([]point, 0)
	for i := range a.sheeps {
		if a.sheeps[i].key == sheep.key {
			continue
		}
		pointsOfSheepInHerd = append(pointsOfSheepInHerd, a.sheeps[i].pos)
	}
	moveToPosition := averagePoints(pointsOfSheepInHerd)
	w.moveAnimal(*sheep, moveToPosition)
	sheep.move(moveToPosition)
	return
}

func getTileType(tileDesc string) (t tile) {
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
	// defining the colors and character styles for all terrains and animals
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

func averagePoints(s []point) (p point) {
	var xvalues int = 0
	var yvalues int = 0
	var count int = 0
	for i := range s {
		xvalues += s[i].x
		yvalues += s[i].y
		count++
	}
	p.x = xvalues / count
	p.y = yvalues / count
	return
}

// Structures
type world struct {
	// Make it concurrent friendly with mutex
	mu *sync.Mutex

	// Length and width sizes for the world
	width  int
	length int

	// tiles stores all the terrain tiles as tile type
	tiles [][]tile

	//Storing terrain tile as points
	landTile  []point
	waterTile []point
}

type tile struct {
	// Contains all information needed to print to screen
	terrainDesc  string
	terrainSym   rune
	terrainStyle tcell.Style
	hasAnimal    bool
	animal       animal
}

type animals struct {
	mu     *sync.Mutex
	sheeps []animal
	wolves []animal
}

type animal struct {
	//path finding variables
	toGo point

	//descriptors
	desc string
	sym  rune
	sty  tcell.Style

	//position
	pos point

	//stats
	health    int
	hunger    int
	speed     int
	maxhealth int
	maxhunger int

	//states
	fleeing bool
	hunting bool
	hungry  bool
	horny   bool
	dead    bool
	rotten  bool

	//index in slice of animals
	key int
}

type point struct {
	x int
	y int
}

//Methods

//generates initial terrain of the world
func (w *world) generateTerrain() {
	rand.Seed(time.Now().UnixNano())
	p := perlin.NewPerlin(2, 2, 10, int64(rand.Int()))
	(*w).tiles = make([][]tile, (*w).width)
	// Initializing the tiles
	for i := range (*w).tiles {
		(*w).tiles[i] = make([]tile, (*w).length)
	}

	waterCount := 0
	landCount := 0
	for x := 0; x < (*w).width; x++ {
		for y := 0; y < (*w).length; y++ {
			terrain := p.Noise2D(float64(x)/10, float64(y)/10)
			if terrain <= -0.12 {
				(*w).tiles[x][y] = getTileType("Water")
				wp := newPoint(x, y)
				(*w).waterTile = append((*w).waterTile, wp)
				waterCount++
			} else if terrain > -0.12 && terrain < 0.3 {
				(*w).tiles[x][y] = getTileType("Land")
				lp := newPoint(x, y)
				(*w).landTile = append((*w).landTile, lp)
				landCount++
			} else {
				(*w).tiles[x][y] = getTileType("Mountain")
				lp := newPoint(x, y)
				(*w).landTile = append((*w).landTile, lp)
				landCount++
			}
		}
	}
}

//calculates the distance from one point to another
func (a *point) distanceTo(b point) (c float32) {
	c = float32(math.Sqrt(float64((b.x-a.x)^2) + (float64((b.y - a.y) ^ 2))))
	return
}

//moves an animal from one point to the input point
func (a *animal) move(p point) {
	(*a).pos = p
	return
}

//updates the w.Tiles array to reflect an animal has moved from current point to next point
func (w *world) moveAnimal(a animal, p point) {
	(*w).tiles[a.pos.x][a.pos.y].hasAnimal = false
	(*w).tiles[p.x][p.y].hasAnimal = true
	(*w).tiles[p.x][p.y].animal = a
	return
}

// Factory Functions
func newWorld(x, y int) (w world) {
	w.width = x
	w.length = y
	w.tiles = make([][]tile, x)
	w.waterTile = make([]point, 0)
	w.landTile = make([]point, 0)
	return
}

func newTile(desc string, sym rune, style tcell.Style, occ bool) tile {
	var t tile
	t.terrainDesc = desc
	t.terrainSym = sym
	t.terrainStyle = style
	t.hasAnimal = occ
	return t
}

func newPoint(x, y int) (p point) {
	p.x = x
	p.y = y
	return
}

func newAnimal(desc string, index int) (a animal) {
	if desc == "Sheep" {
		//set initial state
		a.desc = "Sheep"
		a.sym = 'S'
		a.sty = getSetStyles("Sheep")
		a.key = index

		//set spawn point
		r := rand.Intn(len(w.landTile))
		a.pos = w.landTile[r]

	} else if desc == "Wolf" {
		//set initial state
		a.desc = "Wolf"
		a.sym = 'W'
		a.sty = getSetStyles("Wolf")
		a.key = index

		//set spawn point
		r := rand.Intn(len(w.landTile))
		a.pos = w.landTile[r]

	} else {
		panic("Not a valid animal")
	}
	return
}
