package manager

import (
	"os"

	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
)

// ClearCacheAll is a function that clears the mks cache directory and temporals
func ClearCacheAll() {
	ClearCacheFiles()
	ClearCacheZip()
	ClearCacheTemporals()
}

// ClearCacheFiles is a function that clears the mks template cache directory
func ClearCacheFiles() {
	os.RemoveAll(global.TemplateCachePath)
	os.MkdirAll(global.TemplateCachePath, config.FOLDER_PERMISSION)
}

// ClearCacheZip is a function that clears the mks zip cache directory
func ClearCacheZip() {
	os.RemoveAll(global.ZipCachePath)
	os.MkdirAll(global.ZipCachePath, config.FOLDER_PERMISSION)
}

// ClearCacheTemporals is a function that clears the mks temporals directory
func ClearCacheTemporals() {
	os.RemoveAll(global.TemporalsPath)
	os.MkdirAll(global.TemporalsPath, config.FOLDER_PERMISSION)
}
