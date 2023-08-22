package mks

import (
	"github.com/unknowns24/mks/cmd"
	"github.com/unknowns24/mks/config"
)

func main() {
	config.LoadPackages()
	cmd.Execute()
}
