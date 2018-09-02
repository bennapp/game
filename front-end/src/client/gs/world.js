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
      this.globalPlayerLocation.X = Number(jsonGameState.globalPlayerLocation.X);
      this.globalPlayerLocation.Y = Number(jsonGameState.globalPlayerLocation.Y);
    }

    this.mapStore.setState(jsonGameState, this.globalPlayerLocation);
  }

  objectFromPosition(position) {
    let x = this.globalPlayerLocation.X + position.X;
    let y = this.globalPlayerLocation.Y + position.Y;
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

  move(player, time, direction, conn) {
    if (time > (this.lastMoveTime + this.repeatMoveDelay)) {
      var nextPosition = { x: 0, y: 0 };

      switch (direction) {
        case 'up':
          nextPosition.Y -= 1;
          break;
        case 'left':
          nextPosition.X -= 1;
          break;
        case 'down':
          nextPosition.Y += 1;
          break;
        case 'right':
          nextPosition.X += 1;
          break;
      }

      let nextObject = this.objectFromPosition(nextPosition);

      if (this.isValidMove(nextObject)) {
        this.globalPlayerLocation.X += nextPosition.X;
        this.globalPlayerLocation.Y += nextPosition.Y;

        this.updateObjectRenderLocations();

        this.lastMoveTime = time;

        // This is where 'interactions' happen
        if (nextObject && nextObject.type == 'coin') {
          nextObject.destroy();
        }

        conn.send(JSON.stringify(this.globalPlayerLocation));
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
