package structs

import (
	"sync"

	"github.com/gdamore/tcell"
)

// World contains variables related to the size and terrain of the world. Also stores the Tiles matrix that stores terrain and animal information
type World struct {
	// Make it concurrent friendly with mutex
	mu *sync.Mutex

	// Length and width sizes for theWworld
	width  int
	length int

	// tiles stores all the terrain tiles as tile type
	tiles [][]Tile

	//Storing terrain tile as points
	landTile  []Point
	waterTile []Point
}

// Tile contains information related to a specific cell of the Tiles matrix in the World structure
type Tile struct {
	// Contains all information needed to print to screen
	terrainDesc  string
	terrainSym   rune
	terrainStyle tcell.Style
	hasAnimal    bool
	animal       Animal
}

// Animals contains all animals present within the world
type Animals struct {
	mu     *sync.Mutex
	sheeps []Animal
	wolves []Animal
}

// Animal stores all information related to an animal in the Animals structure
type Animal struct {
	//path finding variables
	toGo Point

	//descriptors
	desc string
	sym  rune
	sty  tcell.Style

	//position
	pos Point

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

// Point stores x and y locations
type Point struct {
	x int
	y int
}
