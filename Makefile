BINNAME = pangaea

ifeq ($(OS),Windows_NT)
	BINNAME = pangaea.exe
endif

all: bin

bin: parser/y.go statik/
	go build -o $(BINNAME)

statik/: native/*
	statik -src=native/ -include=*.pangaea

parser/y.go: parser/parser.go.y
	goyacc -o ./parser/y.go -v ./parser/y.output ./parser/parser.go.y
