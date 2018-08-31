import 'phaser'
import ioClient from 'socket.io-client'
import { WIDTH, HEIGHT } from './constants'

import { World } from './gs/world'
import {Player} from "./el/player";

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
  this.socket = ioClient('http://localhost:8081');
  this.socket.on('currentPlayers', function (players) {
    Object.keys(players).forEach(function (id) {
      if (players[id].playerId === self.socket.id) {
        addPlayer(self, players[id]);
      } else {
        addOtherPlayers(self, players[id]);
      }
    });
  });
  this.socket.on('newPlayer', function (playerInfo) {
    addOtherPlayers(self, playerInfo);
  });
  this.socket.on('disconnect', function (playerId) {
    // self.otherPlayers.getChildren().forEach(function (otherPlayer) {
    //   if (playerId === otherPlayer.playerId) {
    //     otherPlayer.destroy();
    //   }
    // });
  });
  this.socket.on('playerMoved', function (playerInfo) {
    // self.otherPlayers.getChildren().forEach(function (otherPlayer) {
    //   if (playerInfo.playerId === otherPlayer.playerId) {
    //     otherPlayer.setPosition(playerInfo.x, playerInfo.y);
    //   }
    // });
  });

  this.cursors = this.input.keyboard.createCursorKeys();
  self.upKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.W);
  self.leftKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.A);
  self.downKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.S);
  self.rightKey = this.input.keyboard.addKey(Phaser.Input.Keyboard.KeyCodes.D);

  // this.blueScoreText = this.add.text(16, 16, '', { fontSize: '32px', fill: '#0000FF' });

  this.socket.on('scoreUpdate', function (scores) {
    // self.blueScoreText.setText('CoinCount: ' + scores.blue);
  });

  this.socket.on('starLocation', function (starLocation) {
    // if (self.star) self.star.destroy();
    // self.star = self.arcade.add.image(starLocation.x, starLocation.y, 'star');
    // self.arcade.add.overlap(self.ship, self.star, function () {
    //   this.socket.emit('starCollected');
    // }, null, self);
  });

  self.gameStateUpdate = (rawGameState) => {
    let jsonGameState = JSON.parse(rawGameState);
    this.world.setState(jsonGameState);
  };

  this.socket.on('stateUpdate', self.gameStateUpdate);

  let stubbedJsonGameState = {
    globalPlayerLocation: {
      x: '2',
      y: '2',
    },
    coordinates: {
      "0,1": { type: 'coin', id: '33' },
      "3,4": { type: 'rock', id: '-1' },
      "1,1": { type: 'rock', id: '-1' },
    },
    objects: {
      // player: {
      //   "1": {
      //     hp: "10",
      //     alive: "true",
      //     coinCount: "22",
      //   },
      //   "2": {
      //     hp: "7",
      //     alive: "true"
      //   }
      // },
      coin: {
        "33": {
          amount: "11",
        },
      },
      rock: {
        "-1": {}
      }
    },
  };

  this.world.setState(stubbedJsonGameState);

  stubbedJsonGameState = {
    coordinates: {
      "0,1": { type: 'coin', id: '33' },
      "4,4": { type: 'rock', id: '-1' },
      "3,4": { type: 'rock', id: '-1' },
    },
    objects: {
      coin: {
        "33": {
          amount: "11",
        },
      },
      rock: {
        "-1": {}
      }
    },
  };

  this.world.setState(stubbedJsonGameState);

  stubbedJsonGameState = {
    globalPlayerLocation: {
      x: '3',
      y: '3',
    },
    coordinates: {
      "0,1": { type: 'coin', id: '33' },
      "4,4": { type: 'rock', id: '-1' },
      "3,4": { type: 'rock', id: '-1' },
    },
    objects: {
      coin: {
        "33": {
          amount: "11",
        },
      },
      rock: {
        "-1": {}
      }
    },
  };
  this.world.setState(stubbedJsonGameState);
}

function addPlayer(self) {
  // Refactor later when we are ready to handle first response from server,
  // should set gamestate and player location
  self.ship = new Player(self);
}

function addOtherPlayers(self, playerInfo) {
  // const otherPlayer = self.add.sprite(playerInfo.x, playerInfo.y, 'otherPlayer').setOrigin(0.5, 0.5).setDisplaySize(53, 40);
  // if (playerInfo.team === 'blue') {
  //   otherPlayer.setTint(0x0000ff);
  // } else {
  //   otherPlayer.setTint(0xff0000);
  // }
  // otherPlayer.playerId = playerInfo.playerId;
  //self.otherPlayers.add(otherPlayer);
}

function update(time, delta) {
  if (this.ship) {
    let direction = null;
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
      this.world.move(this.ship, time, direction, this.socket);
    }
  }
}
