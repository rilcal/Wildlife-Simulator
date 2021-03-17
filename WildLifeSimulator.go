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
	populateWorld(30, 5)
	a.SheepMaze, a.WolfMaze = structs.GenerateMazes(w)
	defer s.Fini() //***
	rand.Seed(time.Now().UnixNano())
	updateScreen() //***
	time.Sleep(time.Second * 1)
	mainLoop()
	time.Sleep(time.Second * 15)
}

// Functions
func generateInitialVariables() {
	var err error
	s, err = tcell.NewTerminfoScreen() //***
	if err != nil {
		fmt.Println(err)
	}
	s.Init() //***
	x, y := s.Size() //***
	w = structs.NewWorld(x-1, y-1) //***
	//w = structs.NewWorld(40, 40)
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
	for t := 0; t < 1000; t++ {
		// Sheep Logic
		sheepLogic()
		wolfLogic()
		updateScreen() //***
		time.Sleep(time.Millisecond * 500)
	}
}

// Sheep functions
func sheepLogic() {
	for i := 0; i < len(a.Sheeps); i++ {
		if a.Sheeps[i].Dead && a.Sheeps[i].DeadCount > 50 {
			removeAnimal(&a.Sheeps[i])
		} else if a.Sheeps[i].Dead && a.Sheeps[i].Health < 50{
			removeAnimal(&a.Sheeps[i])
		} else if a.Sheeps[i].Dead{
			updateSheepState(&a.Sheeps[i])
		} else {
			if a.Sheeps[i].Fleeing {
				lookForWolves(&a.Sheeps[i])
			} else if a.Sheeps[i].Hungry  {
				eat(&a.Sheeps[i])
			} else if !a.Sheeps[i].Fleeing && !a.Sheeps[i].Hungry {
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

func updateSheepState(sheep *structs.Animal){
	if sheep.Hunger > 0 {
		sheep.Hunger --
	}

	if sheep.Horniness > 0 {
		sheep.Horniness--
	}
	
	if sheep.Health <=0 && !sheep.Dead{
		sheep.Dead = true
		sheep.Sty = structs.GetSetStyles("DeadSheep")
	}

	if sheep.Hunger <= 7 {
		sheep.Hungry = true
	}

	if sheep.Hunger < 0{
		sheep.Health--
	}

	if sheep.Hunger > 7 {
		sheep.Hungry = false
	}

	if sheep.Horniness <= 0 && !sheep.Horny{
		sheep.Horny = true
	}

	if sheep.Dead {
		sheep.DeadCount++
	}

	if sheep.SpeedCount <=0 {
		sheep.SpeedCount = sheep.Speed
	}

	sheep.SpeedCount --

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
		if w.Tiles[wpos].IslandNumber != w.Tiles[spos].IslandNumber{
			continue
		}

		if dist < minDist {
			minDist = dist
			tarSheep = a.Sheeps[i]
		}
		
		if dist <= 1 && a.Sheeps[i].Dead{
			a.Sheeps[i].Health -= 25
			wolf.Hunger += 25
			wolf.ToGo = a.Sheeps[i].Pos
			wolf.ToGoPath = pathfinding.Astar(wpos, spos, a.WolfMaze)
			return
		} else if dist <= 1 && !a.Sheeps[i].Dead{
			a.Sheeps[i].Health -= 25
			wolf.ToGo = a.Sheeps[i].Pos
			wolf.ToGoPath = pathfinding.Astar(wpos, spos, a.WolfMaze)
			return
		}
	}

	if minDist < float32(wolf.Sight) {
		wolf.ToGo = tarSheep.Pos
		wolf.ToGoPath = pathfinding.Astar(wolf.Pos, tarSheep.Pos, a.WolfMaze)
		return
	}

	roam(wolf)
}

// Wolf functions
func wolfLogic() {
	for i := 0; i < len(a.Wolves); i++ {
		if a.Wolves[i].Dead && a.Wolves[i].DeadCount > 10 {
			removeAnimal(&a.Wolves[i])
		} else if a.Wolves[i].Dead{
			updateWolfState(&a.Wolves[i])
		} else {
			if a.Wolves[i].Hungry {
				hunt(&a.Wolves[i])
			} else {
				roam(&a.Wolves[i])
			}

			if a.Wolves[i].SpeedCount <= 0 {
				moveAnimalOnPath(&a.Wolves[i])
			}
			updateWolfState(&a.Wolves[i])

		}
	}
}

func updateWolfState(wolf *structs.Animal){
	if wolf.Hunger > 0 {
		wolf.Hunger--
	}

	if wolf.Horniness > 0 {
		wolf.Horniness--
	}
	
	if wolf.Health <=0 && !wolf.Dead{
		wolf.Dead = true
		wolf.Sty = structs.GetSetStyles("DeadSheep")
	}

	if wolf.Hunger <=0 {
		wolf.Hungry = true
		wolf.Health--
	}else if wolf.Hunger <= 15 {
		wolf.Hungry = true
	}else if wolf.Hunger > 0 {
		wolf.Hungry = false
	}

	if wolf.Horniness <= 0 && !wolf.Horny{
		wolf.Horny = true
	}

	if wolf.Dead {
		wolf.DeadCount++
	}

	if wolf.SpeedCount <=0 {
		wolf.SpeedCount = wolf.Speed
	}

	wolf.SpeedCount --

	tile := w.Tiles[wolf.Pos]
	tile.AnimalType = *wolf
	w.Tiles[wolf.Pos] = tile
}

func roam(ani *structs.Animal){
	for {
		x := rand.Int()%(ani.Sight*2)-ani.Sight
		y := rand.Int()%(ani.Sight*2)-ani.Sight
		loc := structs.NewPoint(ani.Pos.X+x, ani.Pos.Y+y)
		ok,_  := isIn(loc, w.LandTile)
		
		if ok {
			if w.Tiles[loc].IslandNumber == w.Tiles[ani.Pos].IslandNumber{
				ani.ToGo = loc
				ani.ToGoPath = pathfinding.Astar(ani.Pos, loc, a.SheepMaze)
				return
			}
		}
	}
}

func lookForWolves(sheep *structs.Animal){
	sheep.Fleeing = false
	var loc structs.Point
	for x := -sheep.Sight; x <= sheep.Sight; x ++ {
		for y := -sheep.Sight; y <= sheep.Sight; y ++{
			loc = structs.NewPoint(sheep.Pos.X + x, sheep.Pos.Y+ y)
			tile, ok := w.Tiles[loc]
			if ok {
				if tile.HasAnimal {
					if tile.AnimalType.Desc == "Wolf" {
						sheep.Fleeing = true
						var runToPos structs.Point
						var runToDis float32
						runToPos = sheep.Pos
						runToDis = 0
							for yy := sheep.Pos.Y-1; yy <=sheep.Pos.Y+1; yy++ {
								for xx := sheep.Pos.X-1; xx <=sheep.Pos.X+1; xx++ {
									checkPos := structs.NewPoint(xx,yy)
									okk, _ := isIn(checkPos, w.LandTile)
									if okk {
										if checkPos.DistanceTo(tile.AnimalType.Pos) > runToDis {
											runToDis = checkPos.DistanceTo(tile.AnimalType.Pos)
											runToPos = checkPos
										}
									}
								}
							}
						sheep.ToGo = runToPos
						var p []structs.Point = make([]structs.Point, 1)
						p[0] = runToPos
						sheep.ToGoPath = p
						return
					}
				}
			}
		}
	}
}

func findFood(ani *structs.Animal){
	ani.ToGo = findClosestGrassTile(ani.Pos, w)
	ani.ToGoPath = pathfinding.Astar(ani.Pos, ani.ToGo, a.SheepMaze)
	if w.Tiles[ani.Pos].TerrainDesc == "Land" {
		ani.Hunger += 5
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

		if sheep.Pos.DistanceTo(a.Sheeps[j].Pos) <= float32(sheep.Sight) {
			pointsOfSheepInHerd = append(pointsOfSheepInHerd, a.Sheeps[j].Pos)
		}
	}
	averagePoint := structs.AveragePoints(pointsOfSheepInHerd)
	moveToPosition := closestIslandTile(averagePoint, w.Tiles[sheep.Pos].IslandNumber)
	sheep.ToGo = moveToPosition

	if sheep.ToGoPath == nil {
		sheep.ToGoPath = pathfinding.Astar(sheep.Pos, sheep.ToGo, a.SheepMaze)
	}
}

func removeAnimal(ani *structs.Animal){
	tempTile := w.Tiles[ani.Pos]
	tempTile.HasAnimal = false
	w.Tiles[ani.Pos] = tempTile
	findAndRemoveAnimal(ani)
}

func findAndRemoveAnimal(ani *structs.Animal){
	if ani.Desc == "Sheep"{
		for i:= 0; i < len(a.Sheeps); i++ {
			if cmp.Equal(*ani, a.Sheeps[i]){
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

	} else if ani.Desc == "Wolf"{
		if len(a.Wolves) == 1 {
			a.Wolves = nil
		}
		for i:= 0; i < len(a.Wolves); i++ {
			if cmp.Equal(ani, a.Wolves[i]){
				if len(a.Wolves) == 1 {
					a.Wolves = nil
				} else {
					a.Wolves[len(a.Wolves)-1] = a.Wolves[i]
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

func moveAnimal(ani structs.Animal, p structs.Point) (r bool){
	tileFrom, okFrom := w.Tiles[ani.Pos]
	tileTo, okTo := w.Tiles[p]
	if okTo && okFrom {
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

func findClosestLandTile(p structs.Point, wor structs.World) (r structs.Point){
	min := float32(math.Inf(1))
	for _, grass := range wor.LandTile {
		dist := p.DistanceTo(grass)
		if (dist < min) {
			min = dist
			r = grass
		}
	}
	return
}

func findClosestGrassTile(p structs.Point, wor structs.World) (r structs.Point){
	min := float32(math.Inf(1))
	r = p
	for _, grass := range wor.LandTile {
		if w.Tiles[grass].IslandNumber == w.Tiles[p].IslandNumber && w.Tiles[grass].TerrainDesc == "Land"{
			dist := p.DistanceTo(grass)
			if dist < min{
					min = dist
					r = grass
			}
		}
	}
	return
}