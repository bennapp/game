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

func GenerateWorld(regionCoord gs.Coord) {
	startingX := regionCoord.X
	startingY := regionCoord.Y

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
		paint := paintMapping[coord]

		terrainType := generatePaintType(coord, paintMapping)
		paint.SetTerrainByType(terrainType)

		dbs.SavePaintLocation(coord, paint)

		fmt.Printf("saving coord: %v, paint: %v\n", coord, paint)
	}
}

type weightMap map[string]float64

func coordDistance(firstCoord gs.Coord, secondCoord gs.Coord) int {
	diffX := firstCoord.X - secondCoord.X
	diffY := firstCoord.Y - secondCoord.Y

	return math_util.Max(math_util.Abs(diffX), math_util.Abs(diffY))
}

func generatePaintType(coord gs.Coord, paintMapping paintMap) string {
	v := gs.NewVector(-MAX_FILL_DISTANCE, -MAX_FILL_DISTANCE)
	weightedDistanceMap := make(map[int]weightMap)

	for i := 0; i < (MAX_FILL_DISTANCE*2)+1; i++ {
		for j := 0; j < (MAX_FILL_DISTANCE*2)+1; j++ {
			otherCoord := coord.AddVector(v)
			paint := paintMapping[otherCoord]

			distance := coordDistance(coord, otherCoord)

			if weightedDistanceMap[distance] == nil {
				weightedDistanceMap[distance] = make(weightMap)
			}

			if paint != nil && !paint.TerrainEmpty() {
				weightedDistanceMap[distance][paint.TerrainType] += 1
			}

			v.X += 1
		}
		v.X = -MAX_FILL_DISTANCE
		v.Y += 1
	}

	probMap := getBaseWeightProbMap()
	probMap = scaleProbMap(probMap, 0.5)
	neighborWeightProbMap := getNeighborWeightProbMap()

	for distance, wMap := range weightedDistanceMap {
		if distance == 0 {
			continue
		}

		weightedNeighborMap := reWeight(wMap, neighborWeightProbMap)
		dWeight := distance + 1
		scaleWeight := 1 / (math.Pow(2, float64(dWeight)))
		scaledWeightedNeighborMap := scaleProbMap(weightedNeighborMap, scaleWeight)
		probMap = addProbMap(probMap, scaledWeightedNeighborMap)
	}

	terrainType := randomValueFromProbMap(probMap)

	return terrainType
}

var baseWeightProbMap map[string]float64

func getBaseWeightProbMap() map[string]float64 {
	if baseWeightProbMap == nil {
		baseWeightMap := make(map[string]float64)
		baseWeightMap["rock"] = 0
		baseWeightMap["grass"] = 10
		baseWeightMap["sand"] = 2
		baseWeightMap["mud"] = 2

		baseWeightProbMap = weightsToProbMap(baseWeightMap)
	}

	return baseWeightProbMap
}

var neighborWeightProbMap map[string]float64

func getNeighborWeightProbMap() map[string]float64 {
	if neighborWeightProbMap == nil {
		neighborWeightMap := make(map[string]float64)
		neighborWeightMap["rock"] = 0
		neighborWeightMap["grass"] = 3
		neighborWeightMap["sand"] = 3
		neighborWeightMap["mud"] = 3

		neighborWeightProbMap = weightsToProbMap(neighborWeightMap)
	}

	return neighborWeightProbMap
}
