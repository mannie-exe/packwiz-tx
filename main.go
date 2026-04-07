package main

import (
	// Modules of packwiz
	"github.com/mannie-exe/packwiz-tx/cmd"
	_ "github.com/mannie-exe/packwiz-tx/curseforge"
	_ "github.com/mannie-exe/packwiz-tx/github"
	_ "github.com/mannie-exe/packwiz-tx/migrate"
	_ "github.com/mannie-exe/packwiz-tx/modrinth"
	_ "github.com/mannie-exe/packwiz-tx/settings"
	_ "github.com/mannie-exe/packwiz-tx/url"
	_ "github.com/mannie-exe/packwiz-tx/utils"
)

func main() {
	cmd.Execute()
}
