package config

// Mks command arguments
const (
	ARG_CLEAR_CACHE_ALL     = "all"
	ARG_CLEAR_CACHE_FILES   = "cachefiles"
	ARG_CLEAR_CACHE_ZIP     = "cachezip"
	ARG_CLEAR_CACHE_TEMP    = "temporals"
	ARG_CLEAR_CACHE_DEFAULT = ARG_CLEAR_CACHE_ALL
)

// Mks command flags
const (
	FLAG_VERBOSE_LONG  = "verbose"
	FLAG_VERBOSE_SHORT = "v"

	FLAG_FEATURE_LONG  = "feature"
	FLAG_FEATURE_SHORT = "f"

	FLAG_USE_LONG  = "use"
	FLAG_USE_SHORT = "u"
)
