package config

type DatabaseEndpointConfig struct {
	Name     string `yaml:"name"`
	Driver   string `yaml:"driver"`
	Path     string `yaml:"path"`
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Database string `yaml:"database"`
	Password string `yaml:"password"`
	Port     uint16 `yaml:"port"`
	Timeout  int    `yaml:"timeout"`
}

// type SMTPConfig struct {
// 	Host     string `yaml:"host"`
// 	User     string `yaml:"user"`
// 	Email    string `yaml:"email"`
// 	Password string `yaml:"password"`
// 	Port     uint16 `yaml:"port"`
// }

type HTTPServerConfig struct {
	Endpoints     []string `yaml:"endpoints"`
	CookieSecure  bool     `yaml:"cookie_secure"`
	TokenDuration int      `yaml:"token_duration"`
	Domain        string   `yaml:"domain"`
}

type JWTConfig struct {
	SigningKey        string `yaml:"signing_key"`
	Issuer            string `yaml:"issuer"`
	EncryptPassphrase string `yaml:"passphrase"`
	Method            string `yaml:"method"`
}

type Config struct {
	Database DatabaseEndpointConfig `yaml:"database"`
	// SMTPConfig model.SMTPConfig             `yaml:"smtp"`
	HTTPServer HTTPServerConfig `yaml:"http"`
	JWT        JWTConfig        `yaml:"jwt"`

	TemplatesLocation string `yaml:"templates_location"`
	StaticLocation    string `yaml:"static_location"`
	DataLocation      string `yaml:"data_location"`
}
