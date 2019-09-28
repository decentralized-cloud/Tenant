// Package util implements different utilities required by the tenant service
package util

// Set by linker
var (
	version  = "undefined"
	platform = "undefined"
	commit   = "undefined"
	date     = "undefined"
)

// Version defines the structure containns all information to be printed when 'version' command is requested.
type Version struct {
	Version  string `yaml:"version"`
	Platform string
	Commit   string
	Date     string
}

// GetVersion returns the version information
// Returns the version inforamtion
func GetVersion() Version {
	return Version{
		Version:  version,
		Platform: platform,
		Commit:   commit,
		Date:     date,
	}
}
