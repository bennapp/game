import { MapStore } from './map-store'

class World {
  constructor(game) {
    // This will be refactored later when we have game state passed by websockets
    this.lastMoveTime = 0;
    this.repeatMoveDelay = 100;

    this.globalPlayerLocation = {};
    this.mapStore = new MapStore(game);
  }

  setState(jsonGameState) {
    if (jsonGameState.globalPlayerLocation){
      this.globalPlayerLocation.x = Number(jsonGameState.globalPlayerLocation.x);
      this.globalPlayerLocation.y = Number(jsonGameState.globalPlayerLocation.y);
    }

    this.mapStore.setState(jsonGameState, this.globalPlayerLocation);
  }

  objectFromPosition(position) {
    let x = this.globalPlayerLocation.x + position.x;
    let y = this.globalPlayerLocation.y + position.y;
    return this.mapStore.store[`${x},${y}`];
  }

  isValidMove(nextObject) {
    // no object, next move is empty
    if (!nextObject) {
      return true;
    }

    if (nextObject.type === 'rock') {
      return false;
    }

    // defaults to true
    return true;
  }

  move(player, time, direction, socket) {
    if (time > (this.lastMoveTime + this.repeatMoveDelay)) {
      var nextPosition = { x: 0, y: 0 };

      switch (direction) {
        case 'up':
          nextPosition.y -= 1;
          break;
        case 'left':
          nextPosition.x -= 1;
          break;
        case 'down':
          nextPosition.y += 1;
          break;
        case 'right':
          nextPosition.x += 1;
          break;
      }

      let nextObject = this.objectFromPosition(nextPosition);

      if (this.isValidMove(nextObject)) {
        this.globalPlayerLocation.x += nextPosition.x;
        this.globalPlayerLocation.y += nextPosition.y;

        this.updateObjectRenderLocations();

        this.lastMoveTime = time;

        // This is where 'interactions' happen
        if (nextObject && nextObject.type == 'coin') {
          nextObject.destroy();
        }

        socket.emit('playerMovement', this.globalPlayerLocation);
      }
    }
  }

  updateObjectRenderLocations() {
    for (let coordString in this.mapStore.store) {
      let object = this.mapStore.store[coordString];
      object.setNewLocation(this.globalPlayerLocation);
    }
  }
}

export { World }
