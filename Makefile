build-windows:
	GOOS=windows go build -o ./bin/win.easy-ege.exe ./main.go
build-linux:
	GOOS=linux go build -o ./bin/x64linux.easy-ege ./main.go
