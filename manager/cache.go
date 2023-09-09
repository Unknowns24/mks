package manager

import (
	"github.com/unknowns24/mks/config"
	"github.com/unknowns24/mks/global"
	"github.com/unknowns24/mks/utils"
)

// ClearCacheAll is a function that clears the mks cache directory and temporals
func ClearCacheAll() {
	ClearCacheFiles()
	ClearCacheZip()
	ClearCacheTemporals()
}

// ClearCacheFiles is a function that clears the mks template cache directory
func ClearCacheFiles() {
	utils.DeleteFileOrDirectory(global.TemplateCachePath)
	utils.MakeDirectory(global.TemplateCachePath, config.FOLDER_PERMISSION)
}

// ClearCacheZip is a function that clears the mks zip cache directory
func ClearCacheZip() {
	utils.DeleteFileOrDirectory(global.ZipCachePath)
	utils.MakeDirectory(global.ZipCachePath, config.FOLDER_PERMISSION)
}

// ClearCacheTemporals is a function that clears the mks temporals directory
func ClearCacheTemporals() {
	utils.DeleteFileOrDirectory(global.TemporalsPath)
	utils.MakeDirectory(global.TemporalsPath, config.FOLDER_PERMISSION)
}
