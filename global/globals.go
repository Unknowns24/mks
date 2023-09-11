package global

// Command global variables
var (
	Verbose bool // verbose flag for comands (Print what script is doing)
)

// Mks global variables
var (
	InstalledTemplates      = []string{}
	ExportPath              = "" // Path to exports folder
	ZipCachePath            = "" // Path to zip cache folder
	TemporalsPath           = "" // Path to temporals folder
	AutoBackupsPath         = "" // path to store auto backups of feature adds
	TemplateCachePath       = "" // Path to template cache folder
	MksDataFolderPath       = "" // Path to mks data folder inside user config
	MksTemplatesFolderPath  = "" // Path to templates folder inside mks
	UserTemplatesFolderPath = "" // Path to templates folder inside user config
)

// Application global variables
var (
	BasePath        = "" // Current application inside path
	ApplicationName = "" // Current application module name
)
