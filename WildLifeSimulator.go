package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell"
	"github.com/rilcal/Wildlife-Simulator/pathfinding"
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
	a.SheepMaze, a.WolfMaze = structs.GenerateMazes(w)
	defer s.Fini()
	rand.Seed(time.Now().UnixNano())
	updateScreen()
	time.Sleep(time.Second * 5)
	mainLoop()
	time.Sleep(time.Second * 15)
}

// Functions
func generateInitialVariables() {
	var err error
	s, err = tcell.NewScreen()
	if err != nil {
		fmt.Println(err)
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

	waterCount := 0
	landCount := 0
	for x := 0; x < w.Width; x++ {
		for y := 0; y < w.Length; y++ {
			location := structs.NewPoint(x, y)
			terrain := p.Noise2D(float64(x)/10, float64(y)/10)
			tile := w.Tiles[location]

			if terrain <= -0.12 {
				tile = structs.GetTileType("Water")
				wp := structs.NewPoint(x, y)
				w.WaterTile = append(w.WaterTile, wp)
				waterCount++
			} else if terrain > -0.12 && terrain < 0.3 {
				tile = structs.GetTileType("Land")
				lp := structs.NewPoint(x, y)
				w.LandTile = append(w.LandTile, lp)
				landCount++
			} else {
				tile = structs.GetTileType("Mountain")
				lp := structs.NewPoint(x, y)
				w.LandTile = append(w.LandTile, lp)
				landCount++
			}

			w.Tiles[location] = tile
		}
	}

	// Generate island numbers
	var currentIslandCount int
	var currentPoint structs.Point
	var lookedAtPoints []structs.Point
	tempLandTiles := make([]structs.Point, len(w.LandTile))
	var toLookAtPoints []structs.Point
	copy(tempLandTiles, w.LandTile)
	xlen := w.Width
	ylen := w.Length

	for true {
		if len(toLookAtPoints) == 0 && len(tempLandTiles) != 0 {
			currentIslandCount++
			currentPoint = tempLandTiles[0]
		} else if toLookAtPoints == nil && tempLandTiles == nil {
			break
		} else {
			currentPoint = toLookAtPoints[0]
		}
		xcoor := currentPoint.X
		ycoor := currentPoint.Y
		location := structs.NewPoint(xcoor, ycoor)
		tile := w.Tiles[location]
		tile.IslandNumber = currentIslandCount
		w.Tiles[location] = tile

		for x := -1; x < 2; x++ {
			for y := -1; y < 2; y++ {
				neighborX := xcoor + x
				neighborY := ycoor + y
				neighborPoint := structs.NewPoint(neighborX, neighborY)
				beenLookedAt, _ := isIn(neighborPoint, lookedAtPoints)
				isLand, _ := isIn(neighborPoint, w.LandTile)
				pendingLook, _ := isIn(neighborPoint, toLookAtPoints)
				if x == 0 && y == 0 {
					continue
				} else if neighborX < 0 || neighborX >= xlen || neighborY < 0 || neighborY >= ylen {
					continue
				} else if w.Tiles[neighborPoint].IslandNumber != 0 {
					continue
				} else if beenLookedAt {
					continue
				} else if isLand != true {
					continue
				} else if pendingLook {
					continue
				} else {
					toLookAtPoints = append(toLookAtPoints, neighborPoint)
				}
			}
		}
		findAndRemovePoint(currentPoint, &toLookAtPoints)
		findAndRemovePoint(currentPoint, &tempLandTiles)
	}
}

func populateWorld(numSheep int, numWolves int) {
	//Generate sheep population
	a.Sheeps = make([]structs.Animal, numSheep)
	for i := 0; i < numSheep; i++ {
		a.Sheeps[i] = structs.NewAnimal("Sheep", i)
		setLandSpawn(&a.Sheeps[i])

		var sheepLoc structs.Point
		sheepLoc.X = a.Sheeps[i].Pos.X
		sheepLoc.Y = a.Sheeps[i].Pos.Y
		tile, ok := w.Tiles[sheepLoc]
		if ok {
			tile.HasAnimal = true
			tile.AnimalType = a.Sheeps[i]
			w.Tiles[sheepLoc] = tile
		} else {
			panic("Something wrong with the sheep spawn")
		}
	}

	//Generate wolf population
	a.Wolves = make([]structs.Animal, numWolves)
	for j := 0; j < numWolves; j++ {
		a.Wolves[j] = structs.NewAnimal("Wolf", j)
		setLandSpawn(&a.Wolves[j])

		var wolfLoc structs.Point
		wolfLoc.X = a.Wolves[j].Pos.X
		wolfLoc.Y = a.Wolves[j].Pos.Y
		tile, ok := w.Tiles[wolfLoc]
		if ok {
			tile.HasAnimal = true
			tile.AnimalType = a.Wolves[j]
			w.Tiles[wolfLoc] = tile
		} else {
			panic("Something wrong with the sheep spawn")
		}
	}
	return
}

func updateScreen() {

	for point := range w.Tiles {
		tile, ok := w.Tiles[point]
		if ok {
			if tile.HasAnimal {
				s.SetContent(point.X, point.Y, rune(tile.AnimalType.Sym), []rune(""), tile.AnimalType.Sty)
			} else {
				s.SetContent(point.X, point.Y, rune(tile.TerrainSym), []rune(""), tile.TerrainStyle)
			}

		} else {
			panic("Something wrong in updateScreen")
		}
	}
	// for i, ii := range w.Tiles {
	// 	for j := range ii {
	// 		if ii[j].HasAnimal {
	// 			(s).SetContent(i, j, rune(ii[j].AnimalType.Sym), []rune(""), ii[j].AnimalType.Sty)
	// 		} else {
	// 			//rrune := rune(ii[j].IslandNumber)
	// 			//x := []rune{'0', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
	// 			//rrune := x[ii[j].IslandNumber]
	// 			//(s).SetContent(i, j, rune(rrune), []rune(""), ii[j].TerrainStyle)
	// 			(s).SetContent(i, j, rune(ii[j].TerrainSym), []rune(""), ii[j].TerrainStyle)

	// 		}
	// 	}
	// }
	s.Show()
	return
}

func mainLoop() {
	for t := 0; t < 1000; t++ {
		// Sheep Logic
		sheepLogic()
		updateScreen()
		time.Sleep(time.Millisecond * 500)
	}
	return
}

func sheepLogic() {
	for i := range a.Sheeps {
		// herding
		pointsOfSheepInHerd := make([]structs.Point, 0)
		for j := range a.Sheeps {
			if a.Sheeps[i].Pos.DistanceTo(a.Sheeps[j].Pos) <= 10 {
				pointsOfSheepInHerd = append(pointsOfSheepInHerd, a.Sheeps[j].Pos)
			}
			if a.Sheeps[j].Key == a.Sheeps[i].Key {
				continue
			}
		}
		moveToPosition := structs.AveragePoints(pointsOfSheepInHerd)
		a.Sheeps[i].ToGo = moveToPosition

		if a.Sheeps[i].ToGoPath == nil {
			a.Sheeps[i].ToGoPath = pathfinding.Astar(a.Sheeps[i].Pos, a.Sheeps[i].ToGo, a.SheepMaze)
		}

		moveAnimal(a.Sheeps[i], a.Sheeps[i].ToGoPath[0])
		a.Sheeps[i].Pos, a.Sheeps[i].ToGoPath = move(a.Sheeps[i])
	}
	return
}

func setLandSpawn(a *structs.Animal) {
	(*a).Pos = w.LandTile[rand.Intn(len(w.LandTile))]
	return
}

func moveAnimal(a structs.Animal, p structs.Point) {
	tile, ok := w.Tiles[a.Pos]
	if ok {
		tile.HasAnimal = false
		w.Tiles[a.Pos] = tile
	}

	tile, ok = w.Tiles[p]
	if ok {
		tile.HasAnimal = true
		tile.AnimalType = a
		w.Tiles[p] = tile
	}
	return
}

//Move moves an Animal from one Point to the input Point
func move(animal structs.Animal) (p structs.Point, path []structs.Point) {
	if animal.ToGoPath != nil {
		p = animal.ToGoPath[0]
	} else {
		p = animal.Pos
	}

	if len(animal.ToGoPath) == 1 {
		path = nil
	} else {
		path = animal.ToGoPath[1:]
	}
	return
}

func isIn(n structs.Point, slice []structs.Point) (b bool, ind int) {
	for i := range slice {
		if n == slice[i] {
			b = true
			ind = i
			return
		}
	}
	b = false
	ind = 0
	return
}

func findPoint(e structs.Point, l []structs.Point) (b bool, ind int) {
	for i := range l {
		if e == l[i] {
			b = true
			ind = i
			return
		}
	}
	b = false
	ind = 0
	return
}

func findAndRemovePoint(e structs.Point, l *[]structs.Point) (b bool) {
	found, index := findPoint(e, *l)
	if found {
		if len(*l) == 1 {
			*l = nil
		} else {
			(*l)[index] = (*l)[len(*l)-1]
			(*l) = (*l)[:len(*l)-1]
			b = true
			return
		}
	}
	b = false
	return
}
