import { GRID_DISTANCE, GRID_OFFSET, NUM_CELLS } from '../constants'

class Terrain {
  constructor(game, attributes) {
    this.coord = attributes.coord;
    let location = this.spriteLocationFrom(attributes.globalPlayerLocation);

    this.type = attributes.terrainType;

    this.sprite = game.physics.add.sprite(location.x, location.y, attributes.terrainType)
      .setDisplaySize(GRID_DISTANCE, GRID_DISTANCE);
  }

  destroy() {
    this.sprite.destroy();
  }

  setNewLocation(globalPlayerLocation) {
    let location = this.spriteLocationFrom(globalPlayerLocation);

    if (this.coord.x === globalPlayerLocation.X && this.coord.y == globalPlayerLocation.Y) {
        this.destroy()
    }

    this.sprite.x = location.x;
    this.sprite.y = location.y;
  }

  spriteLocationFrom(globalPlayerLocation) {
    let distanceToTopLeft = Math.floor(NUM_CELLS / 2);
    let topLeftCoord = {
      x: globalPlayerLocation.X - distanceToTopLeft,
      y: globalPlayerLocation.Y - distanceToTopLeft,
    };

    let xDistanceFromTopLeftCoord = this.coord.x - topLeftCoord.x;
    let yDistanceFromTopLeftCoord = this.coord.y - topLeftCoord.y;

    return {
      x: (xDistanceFromTopLeftCoord * GRID_DISTANCE) + GRID_OFFSET,
      y: (yDistanceFromTopLeftCoord * GRID_DISTANCE) + GRID_OFFSET
    };
  }
}

export { Terrain };
