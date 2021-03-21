package pathfinding

import (
	"math"

	"github.com/rilcal/Wildlife-Simulator/structs"
)

// Astar follows the a* search algorithm to quickly generate the fastest path to a given target point from a given start point in a matrix of given movement matrix
func Astar(start structs.Point, target structs.Point, matrix map[structs.Point]int) (path []structs.Point) {
	length := len(matrix)

	var current node
	current.pos = start
	current.g = 0
	current.pathto = make([]structs.Point, 0, length)
	current.pathto = append(current.pathto, start)
	current.parentPos = start

	openList := make([]node, 0, length)
	openList = append(openList, current)
	closedList := make([]node, 0, length)

	for i := 1; i < length*10; i++ {
		if current.pos == target {
			tmp := make([]structs.Point, len(current.pathto))
			copy(tmp, current.pathto)
			pth := append(tmp, current.pos)
			path = pth[1:]
			return
		}

		//Generate neighbors slice
		var neighbors []node
		var p structs.Point
		var g int

		for x := -1; x <= 1; x++ {
			for y := -1; y <= 1; y++ {
				if x == 0 && y == 0 {
					continue
				}
				p = structs.NewPoint(current.pos.X+x, current.pos.Y+y)
				score, ok := matrix[p]

				if ok {
					if isIn(p, openList) || isIn(p, closedList) {
						continue
					}
					var pathToCurrent []structs.Point
					pathToParent := make([]structs.Point, len(current.pathto))
					g = current.g + score
					copy(pathToParent, current.pathto)
					pathToCurrent = append(pathToParent, p)
					neighborNode := newNode(p, current.pos, g, pathToCurrent)
					neighbors = append(neighbors, neighborNode)
				}
			}
		}

		//Loop through neighbors, calculate f score, and add to open list
		for i := range neighbors {
			h := math.Max(math.Abs(float64(neighbors[i].pos.X-target.X)), math.Abs(float64(neighbors[i].pos.Y-target.Y)))
			neighbors[i].h = int(h)
			neighbors[i].f = neighbors[i].g + neighbors[i].h
			openList = append(openList, neighbors[i])
		}

		//remove current from open list
		_ = findAndRemoveElement(current, &openList)
		closedList = append(closedList, current)

		//find min f
		_, minInd := findMinF(&openList)
		current = openList[minInd]
	}
	return
}

type node struct {
	pos       structs.Point
	parentPos structs.Point
	pathto    []structs.Point
	g         int
	h         int
	f         int
}

func newNode(pos structs.Point, par structs.Point, g int, parPath []structs.Point) (n node) {
	n.pos = pos
	n.parentPos = par
	n.g = g
	n.pathto = parPath
	return
}

func isIn(n structs.Point, slice []node) (b bool) {
	for i := range slice {
		if n == slice[i].pos {
			b = true
			return
		}
	}
	b = false
	return
}

func findMinFNode(l []node) (n node, ind int) {
	minF := math.Inf(1)
	ind = 0
	for i := range l {
		if float64(l[i].f) < minF {
			minF = float64(l[i].f)
			ind = i
		}
	}
	n = l[ind]
	return
}

func findElement(e node, l *[]node) (b bool, ind int) {
	for i := range *l {
		if e.pos == (*l)[i].pos {
			b = true
			ind = i
			return
		}
	}
	b = false
	ind = 0
	return
}

func findAndRemoveElement(e node, l *[]node) (b bool) {
	found, index := findElement(e, l)
	if found {
		(*l)[index] = (*l)[len(*l)-1]
		(*l) = (*l)[:len(*l)-1]
		b = true
		return
	}
	b = false
	return
}

func findMinF(l *[]node) (minF int, ind int) {
	minF = (*l)[0].f
	ind = 0
	for i := range *l {
		if (*l)[i].f < minF {
			minF = (*l)[i].f
			ind = i
		}
	}
	return
}
