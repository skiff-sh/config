package config

import (
	"fmt"
	"github.com/skiff-sh/appconfig/addrnet"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

// DefaultConfigDir the default directory to house config.
var DefaultConfigDir = "."

// Log represents logging config.
type Log struct {
	Level string `koanf:"level" json:"level" yaml:"level"`
	// Comma-separated output paths.
	// Valid values are:
	// * stdout
	// * stderr
	// * fullfile path
	Outputs List `koanf:"outputs" json:"outputs" yaml:"outputs"`
}

// Server represents communication to some external server.
type Server struct {
	Addr addrnet.Addr `koanf:"addr" json:"addr" yaml:"addr"`
}

// InitKoanf boilerplate for initializing koanf.
func InitKoanf(appName string, def any) *koanf.Koanf {
	k := koanf.NewWithConf(koanf.Conf{
		Delim: ".",
	})

	upperAppName := strings.ToUpper(appName)
	lowerAppName := strings.ToLower(appName)

	// Load defaults
	if def != nil {
		_ = k.Load(structs.Provider(def, "koanf"), nil)
	}

	// Load JSON config.
	_ = k.Load(file.Provider(filepath.Join(DefaultConfigDir, lowerAppName+".json")), json.Parser())

	// Load YAML config.
	_ = k.Load(file.Provider(filepath.Join(DefaultConfigDir, lowerAppName+".yml")), yaml.Parser())
	_ = k.Load(file.Provider(filepath.Join(DefaultConfigDir, lowerAppName+".yaml")), yaml.Parser())

	// Load environment variables
	_ = k.Load(env.Provider(upperAppName+"_", ".",
		func(s string) string {
			out := strings.ReplaceAll(strings.ToLower(strings.TrimPrefix(s, upperAppName+"_")), "_", ".")
			return out
		}), nil)

	return k
}

// ToEnvVars generates a key val map for env vars from a config struct.
func ToEnvVars(appName string, conf any) map[string]string {
	upperAppName := strings.ToUpper(appName)
	out := map[string]string{}
	k := InitKoanf(appName, conf)
	for key, val := range k.All() {
		envKey := upperAppName + "_" + strings.ReplaceAll(strings.ToUpper(key), ".", "_")
		out[envKey] = fmt.Sprintf("%+v", val)
	}

	return out
}

func NewLogger(log Log) (*slog.Logger, error) {
	outputs := log.Outputs.ToSlice()
	w := make([]io.Writer, 0, len(outputs))
	for _, v := range outputs {
		switch v {
		case "stdout":
			w = append(w, os.Stdout)
		case "stderr":
			w = append(w, os.Stderr)
		default:
			f, err := os.OpenFile(v, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
			if err != nil {
				return nil, err
			}
			w = append(w, f)
		}
	}

	logger := slog.New(slog.NewJSONHandler(io.MultiWriter(w...), &slog.HandlerOptions{
		AddSource: true,
		Level:     ParseLevel(log.Level),
	}))

	return logger, nil
}

func ParseLevel(lvl string) slog.Level {
	switch strings.ToLower(lvl) {
	case "info":
		return slog.LevelInfo
	case "debug":
		return slog.LevelDebug
	case "error":
		return slog.LevelError
	case "warn", "warning":
		return slog.LevelWarn
	default:
		return slog.LevelInfo
	}
}
