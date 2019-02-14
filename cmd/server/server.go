package main

import (
	"flag"
	"fmt"

	"github.com/ilyakaznacheev/sdk-server/internal/server"
)

func main() {
	var path string
	flag.StringVar(&path, "path", "module.json", "Path to module list json file")
	flag.Parse()

	err := server.RunServer(path)
	if err != nil {
		fmt.Println(err)
	}
}
