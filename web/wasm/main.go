package main

func main() {
	c := make(chan struct{}, 0)
	RegisterPangaea()

	// block return
	<-c
}
