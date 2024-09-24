package main

import (
	"fmt"
	"syscall/js"
	"time"

	"github.com/davidhorak/space-wars/kernel/game"
	"github.com/davidhorak/space-wars/kernel/physics"
)

func main() {
	fmt.Println("Space Wars")
	fmt.Println("Copyright (C) 2024 David Horak")

	done := make(chan struct{}, 0)
	var instance *game.Game
	jsGlobal := js.Global()

	initializeGameCb := JsFuncIn(func(args []js.Value) {
		method := Method("init", args)

		width, err := method.FloatArg(0, "width")
		if err != nil {
			fmt.Println(err)
		}
		height, err := method.FloatArg(1, "height")
		if err != nil {
			fmt.Println(err)
		}

		seed := time.Now().UnixNano()
		if len(args) == 3 {
			seed, err = method.IntArg(2, "seed")
			if err != nil {
				fmt.Println(err)
			}
		}

		instance = game.NewGame(physics.Size{Width: width, Height: height}, seed)
	})
	tickCb := JsFuncIn(func(args []js.Value) {
		method := Method("tick", args)
		deltaTimeMs, err := method.FloatArg(0, "deltaTimeMs")
		if err != nil {
			fmt.Println(err)
		}
		instance.Update(deltaTimeMs)
	})
	startGameCb := JsFunc(func() { instance.Start() })
	pauseGameCb := JsFunc(func() { instance.Pause() })
	resetGameCb := JsFunc(func() { instance.Reset() })
	gameStateCb := JsFuncOut(func() any { return instance.Serialize() })
	addSpaceshipCb := JsFuncIn(func(args []js.Value) {
		method := Method("addSpaceship", args)

		shipName, err := method.StringArg(0, "shipName")
		if err != nil {
			fmt.Println(err)
		}
		x, err := method.FloatArg(1, "x")
		if err != nil {
			fmt.Println(err)
		}
		y, err := method.FloatArg(2, "y")
		if err != nil {
			fmt.Println(err)
		}
		rotation, err := method.FloatArg(3, "rotation")
		if err != nil {
			fmt.Println(err)
		}

		instance.AddSpaceship(shipName, physics.Vector2{X: x, Y: y}, rotation)
	})
	spaceShipActionCb := JsFuncIn(func(args []js.Value) {
		method := Method("action", args)

		action, err := method.StringArg(0, "action")
		if err != nil {
			fmt.Println(err)
		}
		shipName, err := method.StringArg(1, "shipName")
		if err != nil {
			fmt.Println(err)
		}

		instance.SpaceshipAction(shipName, func(spaceShip *game.Spaceship, gameManager *game.GameManager) {
			switch action {
			case "setEngineThrust":
				mainEngineThrust, err := method.FloatArg(2, "mainEngineThrust")
				if err != nil {
					fmt.Println(err)
					return
				}
				leftEngineThrust, err := method.FloatArg(3, "leftEngineThrust")
				if err != nil {
					fmt.Println(err)
					return
				}
				rightEngineThrust, err := method.FloatArg(4, "rightEngineThrust")
				if err != nil {
					fmt.Println(err)
					return
				}

				spaceShip.SetEngineThrust(mainEngineThrust, leftEngineThrust, rightEngineThrust)
			case "setStartPosition":
				x, err := method.FloatArg(2, "x")
				if err != nil {
					fmt.Println(err)
					return
				}
				y, err := method.FloatArg(3, "y")
				if err != nil {
					fmt.Println(err)
					return
				}
				rotation, err := method.FloatArg(4, "rotation")
				if err != nil {
					fmt.Println(err)
					return
				}

				spaceShip.SetStartPosition(physics.Vector2{X: x, Y: y})
				spaceShip.SetStartRotation(rotation)
			case "fireLaser":
				spaceShip.FireLaser(gameManager)
			case "fireRocket":
				spaceShip.FireRocket(gameManager)
			default:
				fmt.Errorf("invalid action: %s", action)
			}
		})
	})

	jsGlobal.Set("spaceWars", map[string]interface{}{
		"init":         initializeGameCb,
		"tick":         tickCb,
		"start":        startGameCb,
		"pause":        pauseGameCb,
		"reset":        resetGameCb,
		"state":        gameStateCb,
		"addSpaceship": addSpaceshipCb,
		"action":       spaceShipActionCb,
	})

	<-done

	initializeGameCb.Release()
	startGameCb.Release()
	pauseGameCb.Release()
	resetGameCb.Release()
	gameStateCb.Release()
	addSpaceshipCb.Release()
	spaceShipActionCb.Release()
}
