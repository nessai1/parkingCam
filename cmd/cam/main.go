package main

import "parkingCam/internal/cam"

func main() {
	err := cam.Run()
	if err != nil {
		panic(err)
	}
}
