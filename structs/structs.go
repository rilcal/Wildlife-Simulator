package structs

import (
	"math"
	"sync"

	"github.com/gdamore/tcell"
)

// World contains variables related to the size and terrain of the world. Also stores the Tiles matrix that stores terrain and Animal information
type World struct {
	// Make it concurrent friendly with mutex
	Mu *sync.Mutex

	// Length and width sizes for the World
	Width  int
	Length int

	// tiles stores all the terrain tiles as Tile type
	Tiles [][]Tile

	//Storing terrain Tile as points
	LandTile  []Point
	WaterTile []Point
}

// Tile contains information related to a specific cell of the Tiles matrix in the World structure
type Tile struct {
	// Contains all information needed to print to screen
	TerrainDesc  string
	TerrainSym   rune
	TerrainStyle tcell.Style
	HasAnimal    bool
	AnimalType   Animal
}

// Animals contains all animals present within the world
type Animals struct {
	Mu     *sync.Mutex
	Sheeps []Animal
	Wolves []Animal
}

// Animal stores all information related to an Animal in the Animals structure
type Animal struct {
	//path finding variables
	ToGo Point

	//descriptors
	Desc string
	Sym  rune
	Sty  tcell.Style

	//position
	Pos Point

	//stats
	Health    int
	Hunger    int
	Speed     int
	Maxhealth int
	Maxhunger int

	//states
	Fleeing bool
	Hunting bool
	Hungry  bool
	Horny   bool
	Dead    bool
	Rotten  bool

	//index in slice of animals
	Key int
}

// Point stores x and y locations
type Point struct {
	X int
	Y int
}

// DistanceTo calculates the distance from one Point to another
func (a *Point) DistanceTo(b Point) (c float32) {
	c = float32(math.Sqrt(float64((b.X-a.X)^2) + (float64((b.Y - a.Y) ^ 2))))
	return
}

//Move moves an Animal from one Point to the input Point
func (a *Animal) Move(p Point) {
	(*a).Pos = p
	return
}

//MoveAnimal updates the w.Tiles array to reflect an Animal has moved from current Point to next Point
func (w *World) MoveAnimal(a Animal, p Point) {
	(*w).Tiles[a.Pos.X][a.Pos.Y].HasAnimal = false
	(*w).Tiles[p.X][p.Y].HasAnimal = true
	(*w).Tiles[p.X][p.Y].AnimalType = a
	return
}

func GetTileType(tileDesc string) (t Tile) {
	// defining the description, symbols, and style of all tiles
	if tileDesc == "Water" { // Terrain
		t = NewTile("Water", '~', getSetStyles("Water"), false)
	} else if tileDesc == "Land" {
		t = NewTile("Land", '#', getSetStyles("Land"), false)
	} else if tileDesc == "Mountain" {
		t = NewTile("Mountain", 'M', getSetStyles("Mountain"), false)
	} else {
		t = NewTile("", '!', getSetStyles(""), false)
	}
	return
}

//Default styles. "Graphics" for everything to be printed to the screen
func getSetStyles(tileDesc string) (s tcell.Style) {
	// defining the colors and character styles for all terrains and Animals
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

func AveragePoints(s []Point) (p Point) {
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
func NewWorld(x, y int) (w World) {
	w.Width = x
	w.Length = y
	w.Tiles = make([][]Tile, x)
	w.WaterTile = make([]Point, 0)
	w.LandTile = make([]Point, 0)
	return
}

func NewTile(desc string, sym rune, style tcell.Style, occ bool) Tile {
	var t Tile
	t.TerrainDesc = desc
	t.TerrainSym = sym
	t.TerrainStyle = style
	t.HasAnimal = occ
	return t
}

func NewPoint(x, y int) (p Point) {
	p.X = x
	p.Y = y
	return
}

func NewAnimal(desc string, index int) (a Animal) {
	if desc == "Sheep" {
		//set initial state
		a.Desc = "Sheep"
		a.Sym = 'S'
		a.Sty = getSetStyles("Sheep")
		a.Key = index

	} else if desc == "Wolf" {
		//set initial state
		a.Desc = "Wolf"
		a.Sym = 'W'
		a.Sty = getSetStyles("Wolf")
		a.Key = index

	} else {
		panic("Not a valid Animal")
	}
	return
}
