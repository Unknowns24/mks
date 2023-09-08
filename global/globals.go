package global

// Command global variables
var (
	Verbose bool // verbose flag for comands (Print what script is doing)
)

// Mks global variables
var (
	ConfigFolderPath        = ""
	MksTemplatesFolderPath  = "" // Path to templates folder inside mks
	UserTemplatesFolderPath = "" // Path to templates folder inside mks
	InstalledTemplates      = []string{}
)

// Application global variables
var (
	BasePath    = "" // Current application inside path
	ServiceName = "" // Current application module name
)
