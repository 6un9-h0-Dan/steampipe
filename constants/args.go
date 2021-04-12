package constants

// Argument name constants
const (
	ArgDynamic                 = "dynamic"
	ArgHTML                    = "html"
	ArgJSON                    = "json"
	ArgCSV                     = "csv"
	ArgTable                   = "table"
	ArgLine                    = "line"
	ArgListAllTableNames       = "L"
	ArgSelectAll               = "A"
	ArgForce                   = "force"
	ArgTimer                   = "timing"
	ArgOn                      = "on"
	ArgOff                     = "off"
	ArgPortDeprecated          = "db-port"
	ArgPort                    = "database-port"
	ArgListenAddressDeprecated = "listen"
	ArgListenAddress           = "database-listen"
	ArgSearchPath              = "search-path"
	ArgSearchPathPrefix        = "search-path-prefix"
	// search path set in the database config
	ArgServiceSearchPath = "database.search-path"
	// search path set in the terminal config
	ArgSearchPathTerminal = "terminal.search-path"
	ArgInvoker            = "invoker"
	ArgRefresh            = "refresh"
	ArgLogLevel           = "log-level"
	ArgUpdateCheck        = "update-check"
	ArgInstallDir         = "install-dir"
	ArgSqlFile            = "sql-file"
	ArgNoColor            = "no-color"
	ArgNoProgress         = "no-progress"
	ArgOutputFileDir      = "output-file-dir"
	ArgOutputFileFormat   = "output-file-format"
)

/// metaquery mode arguments
var ArgOutput = ArgFromMetaquery(CmdOutput)
var ArgSeparator = ArgFromMetaquery(CmdSeparator)
var ArgHeader = ArgFromMetaquery(CmdHeaders)
var ArgMultiLine = ArgFromMetaquery(CmdMulti)

// BoolToOnOff :: convert a boolean value onto the string "on" or "off"
func BoolToOnOff(val bool) string {
	if val {
		return ArgOn
	}
	return ArgOff
}

// BoolToEnableDisable :: convert a boolean value onto the string "enable" or "disable"
func BoolToEnableDisable(val bool) string {
	if val {
		return "enable"
	}
	return "disable"

}
