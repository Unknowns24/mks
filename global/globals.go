package global

// Command global variables
var (
	Verbose bool // verbose flag for comands (Print what script is doing)
)

// Mks global variables
var (
	ConfigFolderPath        = "" // Path to mks config folder inside user config
	MksTemplatesFolderPath  = "" // Path to templates folder inside mks
	UserTemplatesFolderPath = "" // Path to templates folder inside user config
	InstalledTemplates      = []string{}
	ZipCachePath            = "" // Path to zip cache folder
	TemplateCachePath       = "" // Path to template cache folder
	TemporalsPath           = "" // Path to temporals folder
	ExportPath              = "" // Path to exports folder
)

// Application global variables
var (
	BasePath        = "" // Current application inside path
	ApplicationName = "" // Current application module name
)
