package pathfinding

import "github.com/rilcal/Wildlife-Simulator/structs"

// Astar follows the a* search algorithm to quickly generate the fastest path to a given target point from a given start point in a matrix of given movement matrix
func Astar(start structs.Point, target structs.Point, matrix [][]int) (path []structs.Point) {
	openList := make([]structs.Point, 0, 50)
	openList = append(openList, start)
	//closedList := make([]structs.Point, 0, 50)
	var current node
	current.pos = start
	xlen := len(matrix) - 1
	ylen := len(matrix[1]) - 1

	for i := 1; i > 0; {
		j := 0
		neighbors := make([]node, 0, 8)
		//top left (-1, -1)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y-1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//top mid (0, -1)
		if (current.pos.X) < xlen && (current.pos.X) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X, current.pos.Y-1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//top right (+1, -1)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y-1) < ylen && (current.pos.Y-1) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y-1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//mid left (-1, 0)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y) < ylen && (current.pos.Y) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//mid right (+1, 0)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y) < ylen && (current.pos.Y) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//botm left (-1, +1)
		if (current.pos.X-1) < xlen && (current.pos.X-1) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X-1, current.pos.Y+1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//botm mid (0, +1)
		if (current.pos.X) < xlen && (current.pos.X) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X, current.pos.Y+1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}

		//botm right (+1, +1)
		if (current.pos.X+1) < xlen && (current.pos.X+1) >= 0 && (current.pos.Y+1) < ylen && (current.pos.Y+1) >= 0 {
			p := structs.NewPoint(current.parentPos.X+1, current.pos.Y+1)
			neighbors[j] = newNode(p, current.pos)
			j++
		}
	}
	return
}

type node struct {
	pos       structs.Point
	parentPos structs.Point
}

func newNode(pos, par structs.Point) (n node) {
	n.pos = pos
	n.parentPos = par
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
