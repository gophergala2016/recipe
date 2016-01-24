package main

import (
	"fmt"
	"github.com/gophergala2016/recipe/recipe/cmd"
	"os"
)

func main() {
	if err := cmd.RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
