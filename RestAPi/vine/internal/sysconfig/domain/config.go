package domain

// SystemConfig represents system configuration (sys_cfg table)
type SystemConfig struct {
	ID    int64  // sys_cfgid
	Name  string // sysname
	Value string // defavalue
}

// ConfigType represents different types of system configuration
type ConfigType string

const (
	ConfigTypeFTP     ConfigType = "ftp"
	ConfigTypeFile    ConfigType = "file"
	ConfigTypeVINE    ConfigType = "vine"
	ConfigTypeGeneral ConfigType = "general"
)

// FTPConfig groups FTP-related configurations
type FTPConfig struct {
	Username      string // gcvineftpusername
	Password      string // gcvineftppassword (encrypted)
	ServerName    string // gcvineprimaryftpservername
	DatFolder     string // gcvineftpdatfoldername
	FirewallPort  string // gnvineftpfirewalloutport
	MugshotFolder string // gcvineftpmugshotfoldername
	UseSFTP       bool   // glvineusesftp
}

// FileConfig groups file-related configurations
type FileConfig struct {
	ChargesFileHeader  string // gcvinechargesfileheader
	PrisonerFileHeader string // gcvineprisonerfileheader
	JailIDNumber       string // gcvinejailidnumber
	NewMugshotDir      string // gcvinenewmugshotdirectory
	MugshotDir         string // gcvinemugshotdirectory
	NewVineFilePath    string // gcvinenewvinefilepath
	InterFile          string // gcvineinterfile
}
