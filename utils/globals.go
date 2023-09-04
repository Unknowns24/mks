package utils

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

func SetTemplatesFolderPathGlobal() error {
	// Get the directory path of the current file (generator.go)
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to get current file path")
	}

	mksDir := filepath.Dir(filename)

	// Save in a global variable the path to templates folder inside MKS
	global.TemplatesFolderPath = filepath.Join(mksDir, "..", config.FOLDER_LIBS, config.FOLDER_TEMPLATES)

	return nil
}
