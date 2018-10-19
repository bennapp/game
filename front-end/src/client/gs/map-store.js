import {GRID_DISTANCE, NUM_CELLS} from "../constants";
import {Terrain} from "../el/terrain";
import {Coin} from "../el/coin";
import {Player} from "../el/player";

class MapStore {
  constructor(game){
    this.store = {};
    this.game = game;
  }

  // rawGame state is a string, representing the state of the world the player has loaded
  // the string is JSON formatted
  // The player will see a smaller subset of what they have loaded; their 'vision' will be smaller
  // than their `loaded` vision.
  // I think it would be nice to have coordinates be there 'true' coordinate position, as in ignoring subworlds.
  // The coordinates are key values pretty similar to our redis-store.
  // The objects can be keyed by their type. I think this will allow the client to quickly iterate over objects and be
  // able to easily instantiate them in the client's state.
  // I think the client should update the server / backend with the state of its player but should only be notified form
  // the server of the state of the player when the server detects the client's state of the player is wrong / needs to be corrected.
  // This way, the server is sending gameStates to the client that would 'undo' the client's previous move.

  // merge all objects.
  // Create new objects from `objects` key.
  // Update existing objects from `objects` key.
  // Delete existing objects that are now missing from `objects` key.
  // update all coordinates with new objects

  // Create new objects from `objects` key.
  // Update existing objects from `objects` key.
  // Delete existing objects that are now missing from `objects` key.
  // update all coordinates with new objects

  setState(jsonGameState, globalPlayerLocation) {
    let coordinateState = jsonGameState.coordinates;

    console.log(coordinateState);

    for (var coordString in coordinateState) {
      this.buildObjectFromCoord(coordString, coordinateState, globalPlayerLocation)
    }
  }

  buildObjectFromCoord(coordString, coordinateState, globalPlayerLocation) {
    // FIXME: this needs to do some sort of diffing, right now it is assuming no changes
    // FIXME: There should be some strategy to clean up old coordinates

    if (!this.store[coordString]) {
      this.store[coordString] = this.buildObject(coordString, coordinateState, globalPlayerLocation)
    }
  }

  buildObject(coordString, coordinateState, globalPlayerLocation) {
    let cell = coordinateState[coordString];

    let coordArray = coordString.split(',').map(Number);
    let coord = { x: coordArray[0], y: coordArray[1] };
    let object;

    // FIXME: This prevents writing terrain over user sprite
    if (coord.x === globalPlayerLocation.X && coord.y == globalPlayerLocation.Y) {
        return
    }

    // FIXME: this assumes all items are coins
    if (cell.Items) {
        object = new Coin(this.game, { coord: coord, globalPlayerLocation: globalPlayerLocation });
    }

    if (cell.Paint) {
      if (cell.Paint.TerrainType !== "") {
        object = new Terrain(this.game, { coord: coord, globalPlayerLocation: globalPlayerLocation, terrainType: cell.Paint.TerrainType });
      }
    }

    if (cell.Object) {
      object = new Player(this.game, { coord: coord, globalPlayerLocation: globalPlayerLocation })
    }

    return object
  }
}

export { MapStore }
