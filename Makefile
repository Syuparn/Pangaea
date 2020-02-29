BINNAME = pangaea

ifeq ($(OS),Windows_NT)
	BINNAME = pangaea.exe
endif

all: bin

bin: parser/y.go
	go build -o $(BINNAME)

parser/y.go: parser/parser.go.y
	goyacc -o ./parser/y.go -v ./parser/y.output ./parser/parser.go.y