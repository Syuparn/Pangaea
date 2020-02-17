parser/y.go: parser/parser.go.y
	goyacc -o ./parser/y.go -v ./parser/y.output ./parser/parser.go.y