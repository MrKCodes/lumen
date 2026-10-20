package main

import (
	"encoding/json"
	"fmt"
	"net/url"
)

// importAll upserts goals then entries (with links and tags) for one user.
// It returns the number of entries written.
//
// Everything keys off stable identifiers so re-running is idempotent:
//   - goals  : unique (user_id, slug)
//   - entries: unique (user_id, source_key)
func importAll(c *client, userID string, goals []Goal, entries []Entry) (int, error) {
	// 1. Upsert goals, keep slug→id so entries can reference them.
	goalID := map[string]string{}
	if len(goals) > 0 {
		rows := make([]map[string]any, 0, len(goals))
		for _, g := range goals {
			row := map[string]any{
				"user_id":     userID,
				"slug":        g.Slug,
				"title":       g.Title,
				"description": g.Description,
				"status":      g.Status,
			}
			if g.Color != "" {
				row["color"] = g.Color
			}
			if g.TargetHours > 0 {
				row["target_hours"] = g.TargetHours
			}
			rows = append(rows, row)
		}
		data, err := c.upsert("goals", "user_id,slug", rows)
		if err != nil {
			return 0, fmt.Errorf("upserting goals: %w", err)
		}
		var got []struct {
			ID   string `json:"id"`
			Slug string `json:"slug"`
		}
		if err := json.Unmarshal(data, &got); err != nil {
			return 0, err
		}
		for _, g := range got {
			goalID[g.Slug] = g.ID
		}
	}

	// 2. Upsert entries one by one so we can wire links/tags per entry.
	written := 0
	for _, e := range entries {
		gid, ok := goalID[e.GoalSlug]
		if !ok && e.GoalSlug != "" {
			fmt.Printf("  ! skip %s: unknown goal slug %q (not in goals.yaml)\n", e.SourceKey, e.GoalSlug)
			continue
		}

		row := map[string]any{
			"user_id":    userID,
			"goal_id":    nullable(gid),
			"entry_date": e.Date,
			"title":      e.Title,
			"notes":      e.Notes,
			"minutes":    e.Minutes,
			"source_key": e.SourceKey,
		}
		data, err := c.upsert("entries", "user_id,source_key", []map[string]any{row})
		if err != nil {
			return written, fmt.Errorf("upserting entry %s: %w", e.SourceKey, err)
		}
		var got []struct {
			ID string `json:"id"`
		}
		if err := json.Unmarshal(data, &got); err != nil || len(got) == 0 {
			return written, fmt.Errorf("entry %s: no id returned", e.SourceKey)
		}
		entryID := got[0].ID

		if err := syncLinks(c, userID, entryID, e.Links); err != nil {
			return written, fmt.Errorf("links for %s: %w", e.SourceKey, err)
		}
		if err := syncTags(c, userID, entryID, e.Tags); err != nil {
			return written, fmt.Errorf("tags for %s: %w", e.SourceKey, err)
		}
		written++
	}
	return written, nil
}

// syncLinks replaces an entry's links with the given set (delete-then-insert —
// simplest correct behaviour for a re-import).
func syncLinks(c *client, userID, entryID string, links []Link) error {
	if err := c.deleteWhere("entry_links", url.Values{"entry_id": {"eq." + entryID}}); err != nil {
		return err
	}
	if len(links) == 0 {
		return nil
	}
	rows := make([]map[string]any, 0, len(links))
	for _, l := range links {
		if l.URL == "" {
			continue
		}
		rows = append(rows, map[string]any{
			"entry_id": entryID,
			"user_id":  userID,
			"url":      l.URL,
			"label":    nullable(l.Label),
		})
	}
	if len(rows) == 0 {
		return nil
	}
	_, err := c.upsert("entry_links", "", rows)
	return err
}

// syncTags ensures each tag exists for the user, then links it to the entry.
func syncTags(c *client, userID, entryID string, tags []string) error {
	if err := c.deleteWhere("entry_tags", url.Values{"entry_id": {"eq." + entryID}}); err != nil {
		return err
	}
	for _, name := range tags {
		if name == "" {
			continue
		}
		// upsert the tag (unique user_id,name) and read its id back
		if _, err := c.upsert("tags", "user_id,name",
			[]map[string]any{{"user_id": userID, "name": name}}); err != nil {
			return err
		}
		tagID, err := c.selectIDBy("tags", url.Values{
			"user_id": {"eq." + userID},
			"name":    {"eq." + name},
		})
		if err != nil {
			return err
		}
		if tagID == "" {
			continue
		}
		if _, err := c.upsert("entry_tags", "entry_id,tag_id",
			[]map[string]any{{"entry_id": entryID, "tag_id": tagID}}); err != nil {
			return err
		}
	}
	return nil
}

// nullable turns an empty string into a JSON null so PostgREST stores NULL.
func nullable(s string) any {
	if s == "" {
		return nil
	}
	return s
}
