// Command lumen-import reads a directory of markdown learning logs plus a
// goals.yaml file and upserts them into a Lumen Supabase project.
//
// It is deliberately stdlib-first (only yaml.v3 is external) so it doubles as a
// readable Go learning artifact. See importer/README.md.
//
//	go run . --dir ../content/sre-learnings --owner you@example.com
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// ── flags ────────────────────────────────────────────────────────
	dir := flag.String("dir", "", "content directory holding goals.yaml and *.md (required)")
	owner := flag.String("owner", "", "email of the Lumen user the data belongs to (required)")
	dryRun := flag.Bool("dry-run", false, "parse and print, but do not write to Supabase")
	flag.Parse()

	if *dir == "" || *owner == "" {
		flag.Usage()
		log.Fatal("both --dir and --owner are required")
	}

	// ── parse content ────────────────────────────────────────────────
	// Parsing needs no network, so dry-run works without any Supabase config.
	goals, err := parseGoals(filepath.Join(*dir, "goals.yaml"))
	if err != nil {
		log.Fatalf("goals.yaml: %v", err)
	}
	entries, err := parseEntries(*dir)
	if err != nil {
		log.Fatalf("parsing entries: %v", err)
	}
	fmt.Printf("parsed %d goal(s), %d entr(ies) from %s\n", len(goals), len(entries), *dir)

	if *dryRun {
		for _, g := range goals {
			fmt.Printf("  goal  %-18s %q\n", g.Slug, g.Title)
		}
		for _, e := range entries {
			fmt.Printf("  entry %s  %-14s %4dm  %q\n", e.Date, e.GoalSlug, e.Minutes, e.Title)
		}
		fmt.Println("dry-run: nothing written.")
		return
	}

	// ── config from environment (only needed for a real write) ──────
	// Load .env from the repo root if present (convenience; real env wins).
	loadDotEnv()
	cfg, err := configFromEnv()
	if err != nil {
		log.Fatalf("config: %v", err)
	}

	// ── push to Supabase ─────────────────────────────────────────────
	client := newClient(cfg)

	userID, err := client.userIDByEmail(*owner)
	if err != nil {
		log.Fatalf("resolving owner %q: %v", *owner, err)
	}
	fmt.Printf("owner %s → %s\n", *owner, userID)

	n, err := importAll(client, userID, goals, entries)
	if err != nil {
		log.Fatalf("import: %v", err)
	}
	fmt.Printf("done: %d goal(s) and %d entr(ies) upserted.\n", len(goals), n)
	os.Exit(0)
}
