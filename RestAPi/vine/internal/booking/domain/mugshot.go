package domain

import (
	"fmt"
	"time"
)

// Mugshot represents mugshot image metadata (sys_img table)
type Mugshot struct {
	ID        int64     // sys_imgid
	SystemID  int64     // sysid (name_id)
	SystemKey string    // syskey (always "N" for names)
	Ext1      int       // ext1 (image sequence)
	Ext2      int       // ext2 (version)
	AddTime   time.Time // addtime
}

// FileName generates the mugshot filename
func (m *Mugshot) FileName() string {
	return fmt.Sprintf("%d.%02d%d", m.SystemID, m.Ext1, m.Ext2)
}

// IsForName checks if mugshot is for a name record
func (m *Mugshot) IsForName() bool {
	return m.SystemKey == "N"
}
