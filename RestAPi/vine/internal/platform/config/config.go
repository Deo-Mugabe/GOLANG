package config

type Config struct {
	Environment string
	Server      ServerConfig
	Database    DatabaseConfig
	Redis       RedisConfig
	Queue       QueueConfig
	FTP         FTPConfig
	Files       FilesConfig
	VINE        VINEConfig
	Scheduler   SchedulerConfig
	Auth        AuthConfig
	Logging     LoggingConfig
	Metrics     MetricsConfig
}

type VINEConfig struct {
	ChargesFileHeader  string `mapstructure:"charges_file_header"`
	PrisonerFileHeader string `mapstructure:"prisoner_file_header"`
	JailIDNumber       string `mapstructure:"jail_id_number"`
	RecordSeparator    string `mapstructure:"record_separator"`
	FieldSeparator     string `mapstructure:"field_separator"`
}

type FilesConfig struct {
	NewMugshotDir   string `mapstructure:"new_mugshot_dir"`
	MugshotDir      string `mapstructure:"mugshot_dir"`
	NewVineFilePath string `mapstructure:"new_vine_file_path"`
	InterFile       string `mapstructure:"inter_file"`
	TempDir         string `mapstructure:"temp_dir"`
	MaxFileSize     int64  `mapstructure:"max_file_size"`
}
