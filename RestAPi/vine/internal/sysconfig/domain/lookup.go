package domain

// SystemLookup represents system lookup table (systab1 table)
type SystemLookup struct {
	ID         int    // systab1id
	CodeAgency string // codeAgcy
	CodeKey    string // code_key
	Message    string // sys_msg
}

// LookupKey creates a composite key
func (s *SystemLookup) LookupKey() string {
	return s.CodeAgency + ":" + s.CodeKey
}
