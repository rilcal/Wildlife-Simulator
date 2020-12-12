package pathfinding

import (
	"math"

	"github.com/rilcal/Wildlife-Simulator/structs"
)

// Astar follows the a* search algorithm to quickly generate the fastest path to a given target point from a given start point in a matrix of given movement matrix
func Astar(start structs.Point, target structs.Point, matrix [][]int) (path []structs.Point) {
	xlen := len(matrix) - 1
	ylen := len(matrix[1]) - 1

	var current node
	current.pos = start
	current.g = 0
	current.pathto = make([]structs.Point, 0, 25)

	openList := make([]node, 0, xlen*ylen)
	openList = append(openList, current)

	for i := 1; i < xlen*ylen*10; i++ {
		if current.pos == target {
			path = current.pathto
			return
		}
		j := 0

		//Generate neighbors slice
		neighbors := make([]node, 0, 8)
		//top left (-1, -1)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y-1)
			g := current.g + matrix[current.parentPos.X-1][current.pos.Y-1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//top mid (0, -1)
		if (current.pos.X) < xlen && (current.pos.X) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X, current.pos.Y-1)
			g := current.g + matrix[current.parentPos.X][current.pos.Y-1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//top right (+1, -1)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y-1)
			g := current.g + matrix[current.parentPos.X+1][current.pos.Y-1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//mid left (-1, 0)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y) < ylen && (current.pos.Y) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y)
			g := current.g + matrix[current.parentPos.X-1][current.pos.Y]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//mid right (+1, 0)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y) < ylen && (current.pos.Y) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y)
			g := current.g + matrix[current.parentPos.X+1][current.pos.Y]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//botm left (-1, +1)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y+1)
			g := current.g + matrix[current.parentPos.X-1][current.pos.Y+1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//botm mid (0, +1)
		if (current.pos.X) < xlen && (current.pos.X) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X, current.pos.Y+1)
			g := current.g + matrix[current.parentPos.X][current.pos.Y+1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//botm right (+1, +1)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y+1)
			g := current.g + matrix[current.parentPos.X+1][current.pos.Y+1]
			neighbors[j] = newNode(p, current.pos, g, append(current.pathto, p))
			j++
		}

		//Loop through neighbors, calculate f score, and add to open list
		for i := range neighbors {
			openList = append(openList, neighbors[i])
			h := math.Min(math.Abs(float64(neighbors[i].pos.X-target.X)), math.Abs(float64(neighbors[i].pos.Y-target.Y)))
			neighbors[i].h = int(h)
			neighbors[i].f = neighbors[i].g + neighbors[i].h
		}

		//remove current from open list
		_ = findAndRemoveElement(current, &openList)

		//find min
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

func isIn(n node, slice []structs.Point) (b bool, ind int) {
	for i := range slice {
		if n.pos == slice[i] {
			b = true
			ind = i
			return
		}
	}
	b = false
	ind = 0
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
