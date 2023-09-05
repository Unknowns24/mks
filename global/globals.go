package global

// Command global variables
var (
	Verbose bool // verbose flag for comands (Print what script is doing)
)

// Mks global variables
var (
	TemplatesFolderPath = "" // Path to templates folder inside mks
	ExecutableBasePath  = ""
	InstalledTemplates  = []string{}
)

// Microservice global variables
var (
	BasePath    = "" // Current microservice inside path
	ServiceName = "" // Current microservice module name
)
