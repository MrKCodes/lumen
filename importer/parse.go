package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// ── Goals ────────────────────────────────────────────────────────────

// Goal mirrors one item in goals.yaml.
type Goal struct {
	Slug        string  `yaml:"slug"`
	Title       string  `yaml:"title"`
	Description string  `yaml:"description"`
	Color       string  `yaml:"color"`
	TargetHours float64 `yaml:"target_hours"`
	Status      string  `yaml:"status"`
}

type goalsFile struct {
	Goals []Goal `yaml:"goals"`
}

func parseGoals(path string) ([]Goal, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var gf goalsFile
	if err := yaml.Unmarshal(raw, &gf); err != nil {
		return nil, err
	}
	for i := range gf.Goals {
		if gf.Goals[i].Status == "" {
			gf.Goals[i].Status = "active"
		}
		if gf.Goals[i].Slug == "" || gf.Goals[i].Title == "" {
			return nil, fmt.Errorf("goal %d: slug and title are required", i)
		}
	}
	return gf.Goals, nil
}

// ── Entries ──────────────────────────────────────────────────────────

// Link is a URL attached to an entry. In YAML it is either a bare string or
// a mapping {url, label}; UnmarshalYAML handles both shapes.
type Link struct {
	URL   string
	Label string
}

func (l *Link) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode { // bare URL string
		l.URL = value.Value
		return nil
	}
	var m struct {
		URL   string `yaml:"url"`
		Label string `yaml:"label"`
	}
	if err := value.Decode(&m); err != nil {
		return err
	}
	l.URL, l.Label = m.URL, m.Label
	return nil
}

// Entry is one learning session, ready to upsert.
type Entry struct {
	GoalSlug  string
	Title     string
	Date      string // YYYY-MM-DD
	Minutes   int
	Notes     string
	Tags      []string
	Links     []Link
	SourceKey string // stable key: relpath#index
}

// frontMatter is the YAML block at the top of a markdown file. It supports two
// shapes: a single entry (goal/title/date/minutes/...), or a multi-entry file
// via the entries: list.
type frontMatter struct {
	Goal    string         `yaml:"goal"`
	Title   string         `yaml:"title"`
	Date    string         `yaml:"date"`
	Minutes *int           `yaml:"minutes"`
	Time    string         `yaml:"time"` // alt to minutes, e.g. "1h30m"
	Tags    []string       `yaml:"tags"`
	Links   []Link         `yaml:"links"`
	Entries []subEntry     `yaml:"entries"`
}

type subEntry struct {
	Date    string   `yaml:"date"`
	Title   string   `yaml:"title"`
	Minutes *int     `yaml:"minutes"`
	Time    string   `yaml:"time"`
	Note    string   `yaml:"note"`
	Tags    []string `yaml:"tags"`
	Links   []Link   `yaml:"links"`
}

var (
	frontMatterRe = regexp.MustCompile(`(?s)\A---\r?\n(.*?)\r?\n---\r?\n?(.*)\z`)
	dateFileRe    = regexp.MustCompile(`(\d{4}-\d{2}-\d{2})`)
)

// parseEntries walks dir for *.md files (excluding README.md) and returns all
// entries found, sorted by date then source key for stable output.
func parseEntries(dir string) ([]Entry, error) {
	var entries []Entry

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasSuffix(strings.ToLower(d.Name()), ".md") {
			return nil
		}
		if strings.EqualFold(d.Name(), "README.md") {
			return nil
		}
		rel, _ := filepath.Rel(dir, path)
		es, perr := parseFile(path, rel)
		if perr != nil {
			return fmt.Errorf("%s: %w", rel, perr)
		}
		entries = append(entries, es...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].Date != entries[j].Date {
			return entries[i].Date < entries[j].Date
		}
		return entries[i].SourceKey < entries[j].SourceKey
	})
	return entries, nil
}

func parseFile(path, rel string) ([]Entry, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	m := frontMatterRe.FindSubmatch(raw)
	if m == nil {
		// No frontmatter → skip quietly; it is a prose doc, not a log.
		return nil, nil
	}
	var fm frontMatter
	if err := yaml.Unmarshal(m[1], &fm); err != nil {
		return nil, err
	}
	body := strings.TrimSpace(string(m[2]))

	dateFromName := ""
	if dm := dateFileRe.FindString(filepath.Base(path)); dm != "" {
		dateFromName = dm
	}

	// Multi-entry file.
	if len(fm.Entries) > 0 {
		out := make([]Entry, 0, len(fm.Entries))
		for i, se := range fm.Entries {
			date := firstNonEmpty(se.Date, dateFromName, fm.Date)
			if date == "" {
				return nil, fmt.Errorf("entry %d: no date (set date: in the item or name the file YYYY-MM-DD.md)", i)
			}
			out = append(out, Entry{
				GoalSlug:  fm.Goal,
				Title:     firstNonEmpty(se.Title, fm.Title),
				Date:      date,
				Minutes:   resolveMinutes(se.Minutes, se.Time),
				Notes:     se.Note,
				Tags:      mergeTags(fm.Tags, se.Tags),
				Links:     append(append([]Link{}, fm.Links...), se.Links...),
				SourceKey: fmt.Sprintf("%s#%d", rel, i),
			})
		}
		return out, nil
	}

	// Single-entry file.
	date := firstNonEmpty(fm.Date, dateFromName)
	if date == "" {
		return nil, fmt.Errorf("no date (set date: in frontmatter or name the file YYYY-MM-DD.md)")
	}
	if fm.Goal == "" || fm.Title == "" {
		return nil, fmt.Errorf("goal and title are required")
	}
	return []Entry{{
		GoalSlug:  fm.Goal,
		Title:     fm.Title,
		Date:      date,
		Minutes:   resolveMinutes(fm.Minutes, fm.Time),
		Notes:     body,
		Tags:      fm.Tags,
		Links:     fm.Links,
		SourceKey: fmt.Sprintf("%s#0", rel),
	}}, nil
}

// ── small helpers ────────────────────────────────────────────────────

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

func mergeTags(a, b []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, t := range append(append([]string{}, a...), b...) {
		if t != "" && !seen[t] {
			seen[t] = true
			out = append(out, t)
		}
	}
	return out
}

// resolveMinutes prefers an explicit minutes value; otherwise parses a Go-style
// duration like "1h30m" or "45m". Returns 0 if neither is usable.
func resolveMinutes(minutes *int, timeStr string) int {
	if minutes != nil {
		return *minutes
	}
	if timeStr == "" {
		return 0
	}
	total := 0
	num := ""
	for _, r := range timeStr {
		switch {
		case r >= '0' && r <= '9':
			num += string(r)
		case r == 'h':
			if n, err := strconv.Atoi(num); err == nil {
				total += n * 60
			}
			num = ""
		case r == 'm':
			if n, err := strconv.Atoi(num); err == nil {
				total += n
			}
			num = ""
		}
	}
	return total
}
