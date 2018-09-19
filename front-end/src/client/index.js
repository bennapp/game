import 'phaser'
import { WIDTH, HEIGHT } from './constants'

import { World } from './gs/world'
import { Player } from "./el/player";
import msgpack from "msgpack-lite";

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
  this.load.image('rocks', 'assets/sprites/rock1.png');

  this.load.image('coins1', 'assets/sprites/coins/goldCoin1.png');
  this.load.image('coins2', 'assets/sprites/coins/goldCoin2.png');
  this.load.image('coins3', 'assets/sprites/coins/goldCoin3.png');
  this.load.image('coins4', 'assets/sprites/coins/goldCoin4.png');
  this.load.image('coins5', 'assets/sprites/coins/goldCoin5.png');
  this.load.image('coins6', 'assets/sprites/coins/goldCoin6.png');
  this.load.image('coins7', 'assets/sprites/coins/goldCoin7.png');
  this.load.image('coins8', 'assets/sprites/coins/goldCoin8.png');
  this.load.image('coins9', 'assets/sprites/coins/goldCoin9.png');
}

function create() {
  this.anims.create({
      key: 'coins',
      frames: [
          { key: 'coins1' },
          { key: 'coins2' },
          { key: 'coins3' },
          { key: 'coins4' },
          { key: 'coins5' },
          { key: 'coins6' },
          { key: 'coins7' },
          { key: 'coins8' },
          { key: 'coins9' },
      ],
      frameRate: 15,
      repeat: -1
  });

  var self = this;
  self.world = new World(self);

  self.gameStateUpdate = (gameState) => {
      this.world.setState(gameState);
  };

  if (window["WebSocket"]) {
    console.log('websockts!');
    conn = new WebSocket("ws://" + "localhost:8081" + "/ws");
    conn.binaryType = 'arraybuffer';

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
      //need to cast the raw buffer to a sequence of typed elements
      var typedData = new Uint8Array(event.data)
      var decodedEvent = (msgpack.decode(typedData));

      self.gameStateUpdate(decodedEvent);
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
