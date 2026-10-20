package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// config holds the Supabase connection details the importer needs.
type config struct {
	url            string // https://<ref>.supabase.co
	serviceRoleKey string // secret; bypasses RLS
}

// configFromEnv reads the required Supabase settings from the environment.
func configFromEnv() (config, error) {
	c := config{
		url:            strings.TrimRight(os.Getenv("SUPABASE_URL"), "/"),
		serviceRoleKey: os.Getenv("SUPABASE_SERVICE_ROLE_KEY"),
	}
	if c.url == "" {
		return c, fmt.Errorf("SUPABASE_URL not set")
	}
	if c.serviceRoleKey == "" {
		return c, fmt.Errorf("SUPABASE_SERVICE_ROLE_KEY not set (this is the SECRET key — CLI only)")
	}
	return c, nil
}

// loadDotEnv loads KEY=VALUE lines from a .env file into the process
// environment if they are not already set. It walks up from the current dir
// to find the repo root .env. Missing file is not an error.
func loadDotEnv() {
	for _, path := range dotEnvCandidates() {
		f, err := os.Open(path)
		if err != nil {
			continue
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, val, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			val = strings.Trim(strings.TrimSpace(val), `"'`)
			if _, exists := os.LookupEnv(key); !exists {
				_ = os.Setenv(key, val)
			}
		}
		return // first .env found wins
	}
}

func dotEnvCandidates() []string {
	var out []string
	if cwd, err := os.Getwd(); err == nil {
		dir := cwd
		for i := 0; i < 4; i++ { // cwd and up to 3 parents
			out = append(out, filepath.Join(dir, ".env"))
			parent := filepath.Dir(dir)
			if parent == dir {
				break
			}
			dir = parent
		}
	}
	return out
}
