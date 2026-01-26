package domain

import sysconfigdomain "github.com/Deo-Mugabe/GOLANG/internal/sysconfig/domain"

// TransferConfig encapsulates transfer configuration
type TransferConfig struct {
	FTP  sysconfigdomain.FTPConfig
	File sysconfigdomain.FileConfig
}

// FTPConfig (imported from sysconfig domain)
// FileConfig (imported from sysconfig domain)
