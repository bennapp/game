package wg

import (
	"../dbs"
	"../gs"
	"../math_util"
	"../pnt"
	"fmt"
	"math"
	"math/rand"
)

const MAX_FILL_DISTANCE = 3

type paintMap map[gs.Coord]*pnt.Paint

var baseWeightProbMap map[string]float64
var neighborWeightProbMap map[string]float64

func init() {
	baseWeightMap := make(map[string]float64)
	baseWeightMap["rock"] = 3
	baseWeightMap["grass"] = 7
	baseWeightMap["sand"] = 1
	baseWeightMap["mud"] = 2

	baseWeightProbMap := weightsToProbMap(baseWeightMap)
	baseWeightProbMap = scaleProbMap(baseWeightProbMap, 0.5)

	neighborWeightMap := make(map[string]float64)
	neighborWeightMap["rock"] = 7
	neighborWeightMap["grass"] = 4
	neighborWeightMap["sand"] = 3
	neighborWeightMap["mud"] = 2

	neighborWeightProbMap = weightsToProbMap(neighborWeightMap)
}

func GenerateWorld(regionCoord gs.Coord) {
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

	fmt.Println(paintMapping)
	for _, coord := range coordToGeneratePaint {
		paint := generatePaint(coord, paintMapping)
		paintMapping[coord] = paint
		dbs.SavePaintLocation(coord, paint)
	}
}

type weightMap map[string]float64

func coordDistance(firstCoord gs.Coord, secondCoord gs.Coord) int {
	diffX := firstCoord.X - secondCoord.X
	diffY := firstCoord.Y - secondCoord.Y

	return math_util.Max(math_util.Abs(diffX), math_util.Abs(diffY))
}

func generatePaint(coord gs.Coord, paintMapping paintMap) *pnt.Paint {
	paint := paintMapping[coord]

	v := gs.NewVector(-MAX_FILL_DISTANCE, -MAX_FILL_DISTANCE)
	weightedDistanceMap := make(map[int]weightMap)

	for i := 0; i < (MAX_FILL_DISTANCE*2)+1; i++ {
		for j := 0; j < (MAX_FILL_DISTANCE*2)+1; j++ {
			otherCoord := coord.AddVector(v)
			paint := paintMapping[otherCoord]

			distance := coordDistance(coord, otherCoord)
			fmt.Println(coord, otherCoord, distance)

			if weightedDistanceMap[distance] == nil {
				weightedDistanceMap[distance] = make(weightMap)
			}

			if paint != nil && !paint.TerrainEmpty() {
				weightedDistanceMap[distance][paint.TerrainType] += 1
			}

			v.X += 1
		}
		v.X = -(MAX_FILL_DISTANCE * 2)
		v.Y += 1
	}

	fmt.Println(weightedDistanceMap)

	for distance, wMap := range weightedDistanceMap {
		if distance == 0 {
			continue
		}

		weightedNeighborMap := reWeight(wMap, neighborWeightProbMap)
		scaleWeight := 1 / (math.Pow(2, float64(distance)))

		scaledWeightedNeighborMap := scaleProbMap(weightedNeighborMap, scaleWeight)

		baseWeightProbMap = addProbMap(baseWeightProbMap, scaledWeightedNeighborMap)
	}

	terrainType := randomValueFromProbMap(baseWeightProbMap)
	paint.SetTerrainByType(terrainType)

	return paint
}
