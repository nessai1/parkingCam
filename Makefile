# TODO: make common build of server and camera in one archive
build-camera:
	# TODO: make version postfix naming (by commit tag)
	GOOS=linux GOARCH=amd64 go build -o bin/camera cmd/cam/main.go