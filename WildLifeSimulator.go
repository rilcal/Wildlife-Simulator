package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/aquilax/go-perlin"
	"github.com/gdamore/tcell"
	"github.com/google/go-cmp/cmp"
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
	populateWorld(50, 10)
	a.LandMaze = structs.GenerateMazes(w)
	a.OrigLandMaze = a.LandMaze
	defer s.Fini() //***
	rand.Seed(time.Now().UnixNano())
	updateScreen() //***
	time.Sleep(time.Second * 1)
	mainLoop()
}

// Functions
func generateInitialVariables() {
	var err error
	s, err = tcell.NewScreen() //***
	if err != nil {
		fmt.Println(err)
	}
	s.Init() //***
	//w = structs.NewWorld(100, 100)
	x, y := s.Size()               //***
	w = structs.NewWorld(x-1, y-1) //***
	generateTerrain()
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

			tile.HasAnimal = false
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

	for {
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
				} else if !isLand {
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
	s.Show()
}

func mainLoop() {
	for len(a.Sheeps) > 0 || len(a.Wolves) > 0 {
		sheepLogic()
		wolfLogic()
		updateScreen() //***
		time.Sleep(time.Millisecond * 10)
	}
}

// Sheep functions
func sheepLogic() {
	for i := 0; i < len(a.Sheeps); i++ {
		if a.Sheeps[i].Dead && a.Sheeps[i].DeadCount > 50 {
			removeAnimal(&a.Sheeps[i])
		} else if a.Sheeps[i].Dead && a.Sheeps[i].Health < -50 {
			removeAnimal(&a.Sheeps[i])
		} else if a.Sheeps[i].Dead {
			updateSheepState(&a.Sheeps[i])
		} else {
			if a.Sheeps[i].Fleeing {
				lookForWolves(&a.Sheeps[i])
			} else if a.Sheeps[i].Hungry {
				eat(&a.Sheeps[i])
			} else if a.Sheeps[i].Horny {
				breed(&a.Sheeps[i])
			} else {
				herd(&a.Sheeps[i])
				//roam(&a.Sheeps[i])
				lookForWolves(&a.Sheeps[i])
			}

			if a.Sheeps[i].SpeedCount <= 0 {
				moveAnimalOnPath(&a.Sheeps[i])
			}
			updateSheepState(&a.Sheeps[i])
		}
	}
}

func breed(ani *structs.Animal) {
	if ani.Desc == "Sheep" {
		var closestHornySheep *structs.Animal
		var closestDist int = 10000
		var dist int
		var found bool = false
		var otherSheep *structs.Animal
		for i := range a.Sheeps {
			otherSheep = &a.Sheeps[i]
			if !cmp.Equal(*ani, *otherSheep) {
				if otherSheep.Horny && !otherSheep.Dead {
					dist = int(ani.Pos.DistanceTo(otherSheep.Pos))
					if dist <= ani.Sight && w.Tiles[ani.Pos].IslandNumber == w.Tiles[otherSheep.Pos].IslandNumber {
						if dist < closestDist {
							found = true
							closestHornySheep = otherSheep
							closestDist = dist
						}
					}
				}
			}
		}

		if found {
			if closestDist <= 1 {
				spawnBabyAni(ani, closestHornySheep)
				ani.Horniness = 100
				ani.Horny = false
				closestHornySheep.Horniness = 100
				closestHornySheep.Horny = false
			} else {
				ani.ToGo = closestHornySheep.Pos
				ani.ToGoPath = pathfinding.Astar(ani.Pos, ani.ToGo, a.LandMaze)
			}
		} else {
			herd(ani)
		}

	} else if ani.Desc == "Wolf" {
		var closestHornyWolf *structs.Animal
		var closestDist int = 10000
		var dist int
		var found bool = false
		var otherWolf *structs.Animal
		for i := range a.Wolves {
			otherWolf = &a.Wolves[i]
			if !cmp.Equal(*ani, *otherWolf) {
				if otherWolf.Horny && !otherWolf.Dead {
					dist = int(ani.Pos.DistanceTo(otherWolf.Pos))
					if dist <= ani.Sight && w.Tiles[ani.Pos].IslandNumber == w.Tiles[otherWolf.Pos].IslandNumber {
						if dist < closestDist {
							found = true
							closestHornyWolf = otherWolf
							closestDist = dist
						}
					}
				}
			}
		}

		if found {
			if closestDist <= 1 {
				spawnBabyAni(ani, closestHornyWolf)
				ani.Horniness = 100
				ani.Horny = false
				closestHornyWolf.Horniness = 100
				closestHornyWolf.Horny = false
			} else {
				ani.ToGo = closestHornyWolf.Pos
				ani.ToGoPath = pathfinding.Astar(ani.Pos, ani.ToGo, a.LandMaze)
			}
		}
	}
}

func spawnBabyAni(ani1, ani2 *structs.Animal) {
	var loc structs.Point
	for x := -1; x < 2; x++ {
		for y := -1; y < 2; y++ {
			loc = structs.NewPoint(ani1.Pos.X+x, ani1.Pos.Y+y)
			newTile, ok := w.Tiles[loc]
			if ok { 
				if !w.Tiles[loc].HasAnimal && w.Tiles[loc].TerrainDesc != "Water" {

					ranNum := rand.Int() % 10
					var goodMut float32
					var badMut float32
					if ranNum == 1 {
						goodMut = .5
						badMut = 1.5
					} else if ranNum == 10 {
						goodMut = 1.5
						badMut = .5
					} else {
						goodMut = 1
						badMut = .5
					}

					if ani1.Desc == "Sheep" {
						var babySheep structs.Animal
						babySheep.Desc = ani1.Desc
						babySheep.Sym = ani1.Sym
						babySheep.Sty = ani1.Sty
						babySheep.Pos = loc
						babySheep.Health = int(float32((ani1.Health+ani2.Health)/2) * goodMut)
						babySheep.Hunger = 50
						babySheep.Horniness = int(float32(40) * badMut)
						babySheep.DeadCount = 0
						babySheep.Speed = int(float32((ani1.Speed+ani2.Speed)/2) * badMut)
						babySheep.SpeedCount = babySheep.Speed
						babySheep.Maxhealth = int(float32((ani1.Maxhealth+ani2.Maxhealth)/2) * goodMut)
						babySheep.Maxhunger = int(float32((ani1.Maxhunger+ani2.Maxhunger)/2) * goodMut)
						babySheep.Sight = int(float32((ani1.Sight+ani2.Sight)/2) * goodMut)
						babySheep.Fleeing = false
						babySheep.Hunting = false
						babySheep.Horny = false
						babySheep.Hungry = false
						babySheep.Dead = false
						babySheep.ToGo = babySheep.Pos
						babySheep.ToGoPath = []structs.Point{babySheep.Pos}
						a.Sheeps = append(a.Sheeps, babySheep)

						newTile.HasAnimal = true
						newTile.AnimalType = babySheep

					} else if ani1.Desc == "Wolf" {
						var babyWolf structs.Animal
						babyWolf.Desc = ani1.Desc
						babyWolf.Sym = ani1.Sym
						babyWolf.Sty = ani1.Sty
						babyWolf.Pos = loc
						babyWolf.Health = int(float32((ani1.Health+ani2.Health)/2) * goodMut)
						babyWolf.Hunger = 50
						babyWolf.Horniness = int(float32(40) * badMut)
						babyWolf.DeadCount = 0
						babyWolf.Speed = int(float32((ani1.Speed+ani2.Speed)/2) * badMut)
						babyWolf.SpeedCount = babyWolf.Speed
						babyWolf.Maxhealth = int(float32((ani1.Maxhealth+ani2.Maxhealth)/2) * goodMut)
						babyWolf.Maxhunger = int(float32((ani1.Maxhunger+ani2.Maxhunger)/2) * goodMut)
						babyWolf.Sight = int(float32((ani1.Sight+ani2.Sight)/2) * goodMut)
						babyWolf.Fleeing = false
						babyWolf.Hunting = false
						babyWolf.Horny = false
						babyWolf.Hungry = false
						babyWolf.Dead = false
						babyWolf.ToGo = loc
						babyWolf.ToGoPath = []structs.Point{structs.NewPoint(loc.X, loc.Y)}
						a.Wolves = append(a.Wolves, babyWolf)

						newTile.HasAnimal = true
						newTile.AnimalType = babyWolf
					}

					w.Tiles[loc] = newTile
					return
				}
			}
		}
	}
}

func updateSheepState(sheep *structs.Animal) {
	if sheep.Hunger > 0 {
		sheep.Hunger--
	}

	if sheep.Horniness > 0 {
		sheep.Horniness--
	}

	if sheep.Health <= 0 && !sheep.Dead {
		sheep.Dead = true
		sheep.Sty = structs.GetSetStyles("DeadSheep")
	}

	if sheep.Hunger <= 0 {
		sheep.Hungry = true
		sheep.Health--
	} else if sheep.Hunger <= 10 {
		sheep.Hungry = true
	} else if sheep.Hunger > 0 {
		sheep.Hungry = false
	}

	if sheep.Horniness <= 0 && !sheep.Horny {
		sheep.Horny = true
	}

	if sheep.Dead {
		sheep.DeadCount++
	}

	if sheep.SpeedCount <= 0 {
		sheep.SpeedCount = sheep.Speed
	}

	sheep.SpeedCount--

	tile := w.Tiles[sheep.Pos]
	tile.AnimalType = *sheep
	w.Tiles[sheep.Pos] = tile
}

func eat(sheep *structs.Animal) {
	findFood(sheep)
}

func hunt(wolf *structs.Animal) {
	wpos := wolf.Pos
	var minDist float32 = float32(math.Inf(1))
	var tarSheep structs.Animal

	for i := range a.Sheeps {
		spos := a.Sheeps[i].Pos
		dist := wpos.DistanceTo(spos)
		if w.Tiles[wpos].IslandNumber != w.Tiles[spos].IslandNumber {
			continue
		}

		if dist < minDist {
			minDist = dist
			tarSheep = a.Sheeps[i]
		}

		if dist <= 1 && a.Sheeps[i].Dead {
			a.Sheeps[i].Health -= 25
			wolf.Hunger += 100
			wolf.ToGo = a.Sheeps[i].Pos
			wolf.ToGoPath = pathfinding.Astar(wpos, spos, a.LandMaze)
			return
		} else if dist <= 1 && !a.Sheeps[i].Dead {
			a.Sheeps[i].Health -= 25
			wolf.ToGo = a.Sheeps[i].Pos
			wolf.ToGoPath = pathfinding.Astar(wpos, spos, a.LandMaze)
			return
		}
	}

	if minDist < float32(wolf.Sight) {
		wolf.ToGo = tarSheep.Pos
		wolf.ToGoPath = pathfinding.Astar(wolf.Pos, tarSheep.Pos, a.LandMaze)
		return
	}

	roam(wolf)
}

// Wolf functions
func wolfLogic() {
	for i := 0; i < len(a.Wolves); i++ {
		if a.Wolves[i].Dead && a.Wolves[i].DeadCount > 50 {
			removeAnimal(&a.Wolves[i])
		} else if a.Wolves[i].Dead {
			updateWolfState(&a.Wolves[i])
		} else {
			if a.Wolves[i].Hungry {
				hunt(&a.Wolves[i])
			} else if a.Wolves[i].Horny {
				breed(&a.Wolves[i])
			} else if a.Wolves[i].ToGoPath == nil{
				roam(&a.Wolves[i])
			}

			if a.Wolves[i].SpeedCount <= 0 {
				moveAnimalOnPath(&a.Wolves[i])
			}
			updateWolfState(&a.Wolves[i])

		}
	}
}

func updateWolfState(wolf *structs.Animal) {
	if wolf.Hunger > 0 {
		wolf.Hunger--
	}

	if wolf.Horniness > 0 {
		wolf.Horniness--
	}

	if wolf.Health <= 0 && !wolf.Dead {
		wolf.Dead = true
		wolf.Sty = structs.GetSetStyles("DeadWolf")
	}

	if wolf.Hunger <= 0 {
		wolf.Hungry = true
		wolf.Health--
	} else if wolf.Hunger <= 10 {
		wolf.Hungry = true
	} else if wolf.Hunger > 0 {
		wolf.Hungry = false
	}

	if wolf.Horniness <= 0 && !wolf.Horny {
		wolf.Horny = true
	}

	if wolf.Dead {
		wolf.DeadCount++
	}

	if wolf.SpeedCount <= 0 {
		wolf.SpeedCount = wolf.Speed
	}

	wolf.SpeedCount--

	tile := w.Tiles[wolf.Pos]
	tile.AnimalType = *wolf
	w.Tiles[wolf.Pos] = tile
}

func roam(ani *structs.Animal) {
	x, y := s.Size()
	var rx int
	var ry int
	var tile structs.Tile
	var ok bool
	var loc structs.Point
	for {
		
		rx = rand.Int() % (x-1) 
		ry = rand.Int() % (y-1)
		loc = structs.NewPoint(rx, ry)
		tile, ok = w.Tiles[loc]

		if ok {
			if !(tile.TerrainDesc == "Water") && !tile.HasAnimal {
				if tile.IslandNumber == w.Tiles[ani.Pos].IslandNumber {
					ani.ToGo = loc
					ani.ToGoPath = pathfinding.Astar(ani.Pos, loc, a.LandMaze)
					return
				}
			}
		}
	}
}

func lookForWolves(sheep *structs.Animal){
	spos := sheep.Pos
	var minDist float32 = float32(math.Inf(1))
	var tarWolf structs.Animal
	var found bool

	for i := range a.Wolves {
		wpos := a.Wolves[i].Pos
		dist := spos.DistanceTo(wpos)
		if w.Tiles[spos].IslandNumber != w.Tiles[wpos].IslandNumber {
			continue
		}

		if dist < minDist {
			tarWolf = a.Wolves[i]
			minDist = dist
			found = true
		}
	}
	
	var bestSpot structs.Point = sheep.Pos
	var bestDist float32 = 0
	sheep.Fleeing = false
	if found {
		if minDist < float32(sheep.Sight) {
			sheep.Fleeing = true
			for x := -1; x < 2; x++ {
				for y := -1; y < 2; y++{
					loc := structs.NewPoint(sheep.Pos.X + x, sheep.Pos.Y + y)
					dist := loc.DistanceTo(tarWolf.Pos)
					tile, ok := w.Tiles[loc]
					if ok {
						if dist > bestDist && !tile.HasAnimal && tile.TerrainDesc != "Water" {
							bestSpot = loc
							bestDist = dist
						}
					}
				}
			}
			sheep.ToGo = bestSpot
			sheep.ToGoPath = []structs.Point{bestSpot}
		}
	}
}

func findFood(ani *structs.Animal) {
	ani.ToGo = findClosestGrassTile(ani.Pos, w)
	ani.ToGoPath = pathfinding.Astar(ani.Pos, ani.ToGo, a.LandMaze)
	if w.Tiles[ani.Pos].TerrainDesc == "Land" {
		ani.Hunger += 15
		tile := w.Tiles[ani.Pos]
		tile.TerrainDesc = "DeadGrass"
		tile.TerrainStyle = structs.GetSetStyles("DeadGrass")
		w.Tiles[ani.Pos] = tile
	}
}

func herd(sheep *structs.Animal) {
	pointsOfSheepInHerd := make([]structs.Point, 0)
	for j := range a.Sheeps {
		if w.Tiles[a.Sheeps[j].Pos].IslandNumber != w.Tiles[sheep.Pos].IslandNumber {
			continue
		}

		if sheep.Pos.DistanceTo(a.Sheeps[j].Pos) <= float32(sheep.Sight) && !a.Sheeps[j].Dead{
			pointsOfSheepInHerd = append(pointsOfSheepInHerd, a.Sheeps[j].Pos)
		}
	}
	averagePoint := structs.AveragePoints(pointsOfSheepInHerd)
	moveToPosition := closestIslandTile(averagePoint, w.Tiles[sheep.Pos].IslandNumber)
	sheep.ToGo = moveToPosition

	if sheep.ToGoPath == nil {
		sheep.ToGoPath = pathfinding.Astar(sheep.Pos, sheep.ToGo, a.LandMaze)
	}
}

func removeAnimal(ani *structs.Animal) {
	tempTile := w.Tiles[ani.Pos]
	tempTile.HasAnimal = false
	w.Tiles[ani.Pos] = tempTile
	findAndRemoveAnimal(ani)
}

func findAndRemoveAnimal(ani *structs.Animal) {
	if ani.Desc == "Sheep" {
		for i := 0; i < len(a.Sheeps); i++ {
			if cmp.Equal(*ani, a.Sheeps[i]) {
				if len(a.Sheeps) == 1 {
					a.Sheeps = nil
				} else {
					a.Sheeps[i] = a.Sheeps[len(a.Sheeps)-1]
					a.Sheeps[i].Key = i
					a.Sheeps = a.Sheeps[:len(a.Sheeps)-1]
					return
				}
			}
		}

	} else if ani.Desc == "Wolf" {
		for i := 0; i < len(a.Wolves); i++ {
			if cmp.Equal(*ani, a.Wolves[i]) {
				if len(a.Wolves) == 1 {
					a.Wolves = nil
				} else {
					a.Wolves[i] = a.Wolves[len(a.Wolves)-1]
					a.Wolves[i].Key = i
					a.Wolves = a.Wolves[:len(a.Wolves)-1]
					return
				}
			}
		}
	}
}

func setLandSpawn(a *structs.Animal) {
	(*a).Pos = w.LandTile[rand.Intn(len(w.LandTile))]
}

func moveAnimal(ani structs.Animal, p structs.Point) (r bool) {
	tileFrom, okFrom := w.Tiles[ani.Pos]
	tileTo, okTo := w.Tiles[p]
	if okTo && okFrom {
		a.LandMaze[ani.Pos] = a.OrigLandMaze[ani.Pos]
		a.LandMaze[p] = 100
		tileFrom.HasAnimal = false
		w.Tiles[ani.Pos] = tileFrom
		tileTo.HasAnimal = true
		tileTo.AnimalType = ani
		w.Tiles[p] = tileTo
		r = true
	} else {
		r = false
	}
	return
}

func moveAnimalOnPath(ani *structs.Animal) {
	if ani.ToGoPath == nil {
		return
	}

	if ani.ToGoPath != nil {
		if w.Tiles[ani.ToGoPath[0]].HasAnimal {
			ani.ToGoPath = nil
			return
		}

		moveAnimal(*ani, ani.ToGoPath[0])
		ani.Pos = ani.ToGoPath[0]
		if len(ani.ToGoPath) == 1 {
			ani.ToGoPath = nil
		} else {
			ani.ToGoPath = ani.ToGoPath[1:]
		}
	}
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

func closestIslandTile(p structs.Point, islNum int) (r structs.Point) {
	minDist := float32(math.Inf(1))
	for _, t := range w.LandTile {
		if t == p && w.Tiles[t].IslandNumber == islNum {
			r = p
			return
		}

		dist := t.DistanceTo(p)
		if dist < minDist && w.Tiles[t].IslandNumber == islNum {
			minDist = dist
			r = t
		}

	}
	return
}

func findClosestLandTile(p structs.Point, wor structs.World) (r structs.Point) {
	min := float32(math.Inf(1))
	for _, grass := range wor.LandTile {
		dist := p.DistanceTo(grass)
		if dist < min {
			min = dist
			r = grass
		}
	}
	return
}

func findClosestGrassTile(p structs.Point, wor structs.World) (r structs.Point) {
	min := float32(math.Inf(1))
	r = p
	for _, grass := range wor.LandTile {
		if w.Tiles[grass].IslandNumber == w.Tiles[p].IslandNumber && w.Tiles[grass].TerrainDesc == "Land" && !w.Tiles[grass].HasAnimal {
			dist := p.DistanceTo(grass)
			if dist < min {
				min = dist
				r = grass
			}
		}
	}
	return
}
