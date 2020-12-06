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
