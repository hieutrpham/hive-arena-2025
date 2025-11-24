package main

import (
	"fmt"
	"math/rand"
	"os"
)

import . "hive-arena/common"

var hive_coord = Coords{Row: 0, Col: 0}
var dirs = []Direction{E, SE, SW, W, NW, NE}

func is_entity_type(hex Hex, unit_type EntityType, player int) bool {
	unit := hex.Entity
	return unit != nil && unit.Type == unit_type && unit.Player == player
}

func is_terrain_type(hex Hex, terrain_type Terrain) bool {
	terrain := hex.Terrain
	return terrain == terrain_type
}

func can_forage(hex Hex) bool {
	return is_terrain_type(hex, FIELD) && !hex.Entity.HasFlower
}

func get_hive(state *GameState, player int) {
	for coords, hex := range state.Hexes {
		if is_entity_type(*hex, HIVE, player) {
			hive_coord.Row = coords.Row
			hive_coord.Col = coords.Col
		}
	}
}

func get_hex_type(state *GameState, coord Coords, player int) []string {
	hex := state.Hexes[coord]
	var arr []string
	if hex != nil {
		entity := hex.Entity
		if entity != nil {
			arr = []string{string(hex.Terrain), string(entity.Type)}
		} else {
			arr = []string{string(hex.Terrain), "nil"}
		}
		return arr
	}
	return nil
}

func think(state *GameState, player int) []Order {
	var orders []Order

	// for coords, hex := range state.Hexes {
	// 	unit := hex.Entity
	// 	terrain := hex.Terrain
	// 	fmt.Print("coords: ", coords, " terrain type: ", terrain)
	// 	if unit != nil && unit.Player == player {
	// 		fmt.Println(" entity type: ", unit.Type)
	// 	} else {
	// 		fmt.Println()
	// 	}
	// }

	if hive_coord.Row == 0 && hive_coord.Col == 0 {
		get_hive(state, player)
	}

	for coords, hex := range state.Hexes {
		for _, v := range coords.Neighbours() {
			arr := get_hex_type(state, v, player)
			if arr != nil {
				fmt.Println(arr)
			}
		}

		if is_entity_type(*hex, BEE, player) {
			if can_forage(*hex) {
				orders = append(orders, Order{
					Type:   FORAGE,
					Coords: coords,
				})
			} else {
				orders = append(orders, Order{
					Type:      MOVE,
					Coords:    coords,
					Direction: dirs[rand.Intn(len(dirs))],
				})
			}
		}
	}

	return orders
}

func main() {
	if len(os.Args) <= 3 {
		fmt.Println("Usage: ./agent <host> <gameid> <name>")
		os.Exit(1)
	}

	host := os.Args[1]
	id := os.Args[2]
	name := os.Args[3]

	Run(host, id, name, think)
}
