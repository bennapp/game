package wg

import (
	"../dbs"
	"../gs"
	"../pnt"
	"math/rand"
)

const MAX_FILL_DISTANCE = 3

type paintMap map[gs.Coord]*pnt.Paint

func generateWorld(regionCoord gs.Coord) {
	startingX := regionCoord.X * gs.WORLD_GENERATION_DISTANCE
	startingY := regionCoord.Y * gs.WORLD_GENERATION_DISTANCE

	endingX := startingX + gs.WORLD_GENERATION_DISTANCE
	endingY := startingY + gs.WORLD_GENERATION_DISTANCE

	outerStartingX := startingX - MAX_FILL_DISTANCE
	outerStartingY := startingY - MAX_FILL_DISTANCE

	outerEndingX := startingX + gs.WORLD_GENERATION_DISTANCE + MAX_FILL_DISTANCE
	outerEndingY := startingY + gs.WORLD_GENERATION_DISTANCE + MAX_FILL_DISTANCE

	paintMapping := make(paintMap)

	for i := outerStartingX; i < outerEndingX; i++ {
		for j := outerStartingY; j < outerEndingY; j++ {
			coord := gs.NewCoord(i, j)

			if i >= startingX && i < endingX && j >= startingY && j < endingY {
				paintMapping[coord] = pnt.NewPaintWithEmptyTerrain()
			} else {
				paintMapping[coord] = dbs.LoadPaintByCoord(coord)
			}
		}
	}

	var coordToGeneratePaint []gs.Coord
	for i := startingX; i < endingX; i++ {
		for j := startingY; j < endingY; j++ {
			coord := gs.NewCoord(i, j)
			coordToGeneratePaint = append(coordToGeneratePaint, coord)
		}
	}

	// shuffle coords
	for i := range coordToGeneratePaint {
		j := rand.Intn(i + 1)
		coordToGeneratePaint[i], coordToGeneratePaint[j] = coordToGeneratePaint[j], coordToGeneratePaint[i]
	}

	for _, coord := range coordToGeneratePaint {
		paint := generatePaint(coord, paintMapping)
		dbs.SavePaintLocation(coord, paint)
	}
}

func generatePaint(coord gs.Coord, paintMapping paintMap) *pnt.Paint {
	paint := paintMapping[coord]

	// base probability of each type
	// grass 0.5
	// rock 0.2
	// mud 0.2
	// sand 0.1

	// distance coeff, lets use inverse square

	// neighbor probability of each type
	// rock 0.7
	// sand 0.5
	// grass 0.5
	// mud 0.3

	// base +
	// depth(i) : N * (distance_coeff(depth-i) * neighbor probability)

	return paint
}
