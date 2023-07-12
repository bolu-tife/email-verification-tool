package main

func main() {
	port := GetConfig().Port
	server := NewAPIServer(":" + port)
	server.Run()
}
