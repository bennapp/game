import 'phaser'
import { WIDTH, HEIGHT } from './constants'

import { World } from './gs/world'
import { Player } from "./el/player";

var config = {
  type: Phaser.AUTO,
  parent: 'phaser-example',
  width: WIDTH,
  height: HEIGHT,
  physics: {
    default: 'arcade',
    arcade: {
      debug: true,
      gravity: {
        x: 0,
        y: 0
      }
    }
  },
  scene: {
    preload: preload,
    create: create,
    update: update
  } 
};

new Phaser.Game(config);
var world;
var conn;

function preload() {
  // TODO WEBPACK ASSETS
  this.load.image('ship', 'assets/spaceShips_001.png');
  this.load.image('otherPlayer', 'assets/enemyBlack5.png');
  this.load.image('star', 'assets/star_gold.png');
  this.load.image('rocks', 'assets/sprites/Rock Pile.png');
}

function create() {
  var self = this;
  self.world = new World(self);

  self.gameStateUpdate = (rawGameState) => {
      let jsonGameState = JSON.parse(rawGameState);
      this.world.setState(jsonGameState);
  };

  if (window["WebSocket"]) {
    console.log('websockts!');
    conn = new WebSocket("ws://" + "localhost:8081" + "/ws");

    conn.onopen = function (event) {
      self.ship = new Player(self);

      window.onbeforeunload = function() {
        conn.onclose = function () {};
        conn.close()
      };
    };

    conn.onclose = function (event) {
      console.log("Connection closed.");
    };
    conn.onmessage = function (event) {
      console.log(JSON.parse(event.data));
      self.gameStateUpdate(event.data);
    };
  } else {
    console.log("Your browser does not support WebSockets.");
  }

  this.cursors = this.input.keyboard.createCursorKeys();
  self.upKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.W);
  self.leftKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.A);
  self.downKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.S);
  self.rightKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.D);
}

function update(time, delta) {
  if (this.ship) {
    let direction;
    if (this.cursors.up.isDown || this.upKey.isDown) {
      direction = 'up';
    } else if (this.cursors.left.isDown || this.leftKey.isDown) {
      direction = 'left';
    } else if(this.cursors.down.isDown || this.downKey.isDown) {
      direction = 'down';
    } else if(this.cursors.right.isDown || this.rightKey.isDown) {
      direction = 'right';
    }

    if (direction) {
      this.world.move(this.ship, time, direction, conn);
      direction = null;
    }
  }
}
