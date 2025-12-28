package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/k0kubun/pp/v3"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Server   ServerConfig
	Session  SessionConfig
	Database DatabaseConfig
	Casbin   CasbinConfig
}

// ------------------------------------------------------------
// Server / Gin
// ------------------------------------------------------------

type ServerConfig struct {
	Host     string `envconfig:"HOST" default:"127.0.0.1"`
	Port     int    `envconfig:"PORT" default:"8080"`
	BasePath string `envconfig:"BASE_PATH" default:"/"`

	// Wenn gesetzt → Unix-Socket ($XDG_RUNTIME_DIR/<Socket>)
	Socket string `envconfig:"SERVER_SOCKET"`

	// true = Development Mode, false = Release Mode (Default)
	Dev bool `envconfig:"DEV" default:"false"`
}

// ------------------------------------------------------------
// Session
// ------------------------------------------------------------

type SessionConfig struct {
	Secret     string        `envconfig:"SECRET" required:"true"`
	MaxAge     time.Duration `envconfig:"MAX_AGE" default:"720h"` // 30 Tage
	CookieName string        `envconfig:"COOKIE_NAME" default:"my_hagg_app"`
}

// ------------------------------------------------------------
// Database
// ------------------------------------------------------------

type DatabaseConfig struct {
	SQLite   SQLiteConfig
	External ExternalDatabasesConfig
}

type SQLiteConfig struct {
	Path string `envconfig:"SQLITE_PATH" default:"./db.sqlite3"`
}

type ExternalDatabasesConfig struct {
	// später
}

// ------------------------------------------------------------
// Casbin
// ------------------------------------------------------------

type CasbinConfig struct {
	ModelPath  string `envconfig:"MODEL"  default:"model.conf"`
	PolicyPath string `envconfig:"POLICY" default:"policy.csv"`
}

// ------------------------------------------------------------
// Load
// ------------------------------------------------------------

func Load() (*Config, error) {
	// .env optional
	_ = godotenv.Load()

	var server ServerConfig
	if err := envconfig.Process("GIN", &server); err != nil {
		return nil, fmt.Errorf("load server config: %w", err)
	}

	var session SessionConfig
	if err := envconfig.Process("SESSION", &session); err != nil {
		return nil, fmt.Errorf("load session config: %w", err)
	}

	var database DatabaseConfig
	if err := envconfig.Process("DB", &database); err != nil {
		return nil, fmt.Errorf("load database config: %w", err)
	}

	// ----------------------------
	// Casbin
	// ----------------------------

	var casbinCfg CasbinConfig
	if err := envconfig.Process("CASBIN", &casbinCfg); err != nil {
		return nil, fmt.Errorf("load casbin config: %w", err)
	}

	cfg := &Config{
		Server:   server,
		Session:  session,
		Database: database,
		Casbin:   casbinCfg,
	}

	if err := cfg.validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(err)
	}
	return cfg
}

// ------------------------------------------------------------
// Validation
// ------------------------------------------------------------

func (c *Config) validate() error {
	if c.Server.Socket == "" {
		if c.Server.Port <= 0 || c.Server.Port > 65535 {
			return fmt.Errorf("invalid GIN_PORT: %d", c.Server.Port)
		}
	}

	if c.Server.BasePath == "" {
		return fmt.Errorf("GIN_BASE_PATH must not be empty")
	}

	if c.Database.SQLite.Path == "" {
		return fmt.Errorf("DB_SQLITE_PATH must not be empty")
	}

	if c.Casbin.ModelPath == "" {
		return fmt.Errorf("CASBIN_MODEL must not be empty")
	}

	if c.Casbin.PolicyPath == "" {
		return fmt.Errorf("CASBIN_POLICY must not be empty")
	}

	return nil
}

// ------------------------------------------------------------
// Helper
// ------------------------------------------------------------

func (c *Config) Addr() string {
	return c.Server.Host + ":" + strconv.Itoa(c.Server.Port)
}

func (c *Config) BaseURL() string {
	return "http://" + c.Addr() + c.Server.BasePath
}

func (c *Config) Pretty() {
	pp.Println(c)
}

func (c *Config) Sprint() string {
	return pp.Sprint(c)
}

func (c Config) Print() {
	fmt.Println("Config")

	printServer(c.Server)
	printDatabase(c.Database)
	printSession(c.Session)
	printCasbin(c.Casbin)
}

func printServer(s ServerConfig) {
	fmt.Println("├─ Server (GIN)")

	mode := "release"
	if s.Dev {
		mode = "debug"
	}

	fmt.Printf("│  ├─ Mode     : %s\n", mode)

	if s.Socket != "" {
		fmt.Printf("│  ├─ Socket   : %s\n", s.Socket)
	} else {
		fmt.Printf("│  ├─ Host     : %s\n", s.Host)
		fmt.Printf("│  ├─ Port     : %d\n", s.Port)
	}

	fmt.Printf("│  └─ BasePath : %s\n", s.BasePath)
}

func printDatabase(d DatabaseConfig) {
	fmt.Println("├─ Database")
	fmt.Println("│  └─ SQLite")
	fmt.Printf("│     └─ Path : %s\n", d.SQLite.Path)
}

func printSession(s SessionConfig) {
	fmt.Println("├─ Session")
	fmt.Printf("│  ├─ CookieName : %s\n", s.CookieName)
	fmt.Printf("│  ├─ MaxAge     : %s\n", s.MaxAge)
	fmt.Printf("│  └─ Secret     : %s\n", s.Secret)
}

func printCasbin(c CasbinConfig) {
	fmt.Println("└─ Casbin")
	fmt.Printf("   ├─ Model  : %s\n", c.ModelPath)
	fmt.Printf("   └─ Policy : %s\n", c.PolicyPath)
}
