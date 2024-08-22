package main

import "parkingCam/internal/server"

func main() {
	err := server.Run()
	if err != nil {
		panic(err)
	}
}
