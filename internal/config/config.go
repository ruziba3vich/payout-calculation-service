package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	App      AppConfig
	HTTP     HTTPConfig
	Postgres PostgresConfig
	Redis    RedisConfig
	JWT      JWTConfig
	Admin    AdminConfig
}

// AdminConfig seeds a first admin on startup when both fields are set.
type AdminConfig struct {
	Username string
	Password string
}

type AppConfig struct {
	Name     string
	Env      string
	LogLevel string
}

type HTTPConfig struct {
	Port            string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration
}

type PostgresConfig struct {
	URL             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	MigrationsPath  string
}

type RedisConfig struct {
	Addr          string
	Password      string
	DB            int
	PayoutTTL     time.Duration
	PayoutListTTL time.Duration
	CourierTTL    time.Duration
}

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
	Issuer     string
}

func Load() (Config, error) {
	var errs []error

	cfg := Config{
		App: AppConfig{
			Name:     getString("APP_NAME", "payout-calculation-service"),
			Env:      getString("APP_ENV", "local"),
			LogLevel: getString("LOG_LEVEL", "info"),
		},
		HTTP: HTTPConfig{
			Port:            getString("HTTP_PORT", "8080"),
			ReadTimeout:     getDuration("HTTP_READ_TIMEOUT", 10*time.Second, &errs),
			WriteTimeout:    getDuration("HTTP_WRITE_TIMEOUT", 15*time.Second, &errs),
			IdleTimeout:     getDuration("HTTP_IDLE_TIMEOUT", 60*time.Second, &errs),
			ShutdownTimeout: getDuration("HTTP_SHUTDOWN_TIMEOUT", 10*time.Second, &errs),
		},
		Postgres: PostgresConfig{
			URL:             mustString("DB_URL", &errs),
			MaxConns:        int32(getInt("DB_MAX_CONNS", 20, &errs)),
			MinConns:        int32(getInt("DB_MIN_CONNS", 2, &errs)),
			MaxConnLifetime: getDuration("DB_MAX_CONN_LIFETIME", time.Hour, &errs),
			MaxConnIdleTime: getDuration("DB_MAX_CONN_IDLE_TIME", 15*time.Minute, &errs),
			MigrationsPath:  getString("MIGRATIONS_PATH", "./migrations"),
		},
		Redis: RedisConfig{
			Addr:          mustString("REDIS_ADDR", &errs),
			Password:      getString("REDIS_PASSWORD", ""),
			DB:            getInt("REDIS_DB", 0, &errs),
			PayoutTTL:     getDuration("REDIS_PAYOUT_TTL", 10*time.Minute, &errs),
			PayoutListTTL: getDuration("REDIS_PAYOUT_LIST_TTL", 2*time.Minute, &errs),
			CourierTTL:    getDuration("REDIS_COURIER_TTL", 10*time.Minute, &errs),
		},
		JWT: JWTConfig{
			Secret:     mustString("JWT_SECRET", &errs),
			AccessTTL:  getDuration("JWT_ACCESS_TTL", 15*time.Minute, &errs),
			RefreshTTL: getDuration("JWT_REFRESH_TTL", 7*24*time.Hour, &errs),
			Issuer:     getString("JWT_ISSUER", "payout-calculation-service"),
		},
		Admin: AdminConfig{
			Username: getString("ADMIN_USERNAME", ""),
			Password: getString("ADMIN_PASSWORD", ""),
		},
	}

	if len(cfg.JWT.Secret) > 0 && len(cfg.JWT.Secret) < 32 {
		errs = append(errs, errors.New("JWT_SECRET must be at least 32 characters"))
	}
	if len(errs) > 0 {
		return Config{}, fmt.Errorf("config: %w", errors.Join(errs...))
	}
	return cfg, nil
}

func (c Config) IsProd() bool {
	return strings.EqualFold(c.App.Env, "prod")
}

func getString(key, def string) string {
	if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
		return v
	}
	return def
}

func mustString(key string, errs *[]error) string {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		*errs = append(*errs, fmt.Errorf("%s is required", key))
		return ""
	}
	return v
}

func getInt(key string, def int, errs *[]error) int {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: invalid int %q", key, v))
		return def
	}
	return n
}

func getDuration(key string, def time.Duration, errs *[]error) time.Duration {
	v, ok := os.LookupEnv(key)
	if !ok || strings.TrimSpace(v) == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		*errs = append(*errs, fmt.Errorf("%s: invalid duration %q", key, v))
		return def
	}
	return d
}
