package main

import (
	"fmt"
	"os"

	"github.com/bordviz/datasphere/internal/config"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%+v\n", cfg)
}
