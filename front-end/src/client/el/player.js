import {GRID_DISTANCE, GRID_OFFSET, NUM_CELLS} from '../constants'

class Player {
  constructor(game) {
    let location = this.spriteLocationFrom();

    // this class is the user's player
    // we will probably want to make another one for other players
    this.type = 'client-player';

    this.sprite = game.physics.add.sprite(location.x, location.y, 'ship')
      .setDisplaySize(GRID_DISTANCE, GRID_DISTANCE);
  }

  destroy() {
    this.sprite.destroy();
  }

  spriteLocationFrom() {
    let distanceFromCenter = Math.floor(NUM_CELLS / 2);

    return {
      x: (distanceFromCenter * GRID_DISTANCE) + GRID_OFFSET,
      y: (distanceFromCenter * GRID_DISTANCE) + GRID_OFFSET
    };
  }
}

export { Player };
