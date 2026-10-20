package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// client is a thin Supabase REST + Auth-admin client using the service_role
// key. The service_role key bypasses RLS, so this must only ever run server /
// CLI side — never in a browser.
type client struct {
	cfg  config
	http *http.Client
}

func newClient(cfg config) *client {
	return &client{cfg: cfg, http: &http.Client{Timeout: 30 * time.Second}}
}

// do performs a request against the Supabase REST (PostgREST) or Auth API and
// returns the response body, or an error carrying the status + body on != 2xx.
func (c *client) do(method, path string, query url.Values, body any) ([]byte, error) {
	var reader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reader = bytes.NewReader(b)
	}

	u := c.cfg.url + path
	if len(query) > 0 {
		u += "?" + query.Encode()
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", c.cfg.serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.cfg.serviceRoleKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("%s %s → %d: %s", method, path, resp.StatusCode, string(data))
	}
	return data, nil
}

// userIDByEmail resolves a Lumen auth user's UUID from their email via the
// GoTrue admin API, paging until found.
func (c *client) userIDByEmail(email string) (string, error) {
	for page := 1; page <= 50; page++ {
		q := url.Values{"page": {fmt.Sprint(page)}, "per_page": {"200"}}
		data, err := c.do(http.MethodGet, "/auth/v1/admin/users", q, nil)
		if err != nil {
			return "", err
		}
		var out struct {
			Users []struct {
				ID    string `json:"id"`
				Email string `json:"email"`
			} `json:"users"`
		}
		if err := json.Unmarshal(data, &out); err != nil {
			return "", err
		}
		if len(out.Users) == 0 {
			break
		}
		for _, u := range out.Users {
			if u.Email == email {
				return u.ID, nil
			}
		}
	}
	return "", fmt.Errorf("no user with email %q — register in the app first", email)
}

// upsert POSTs rows to a table with on-conflict merge, returning the created /
// updated rows (via Prefer: return=representation).
func (c *client) upsert(table string, onConflict string, rows any) ([]byte, error) {
	q := url.Values{}
	if onConflict != "" {
		q.Set("on_conflict", onConflict)
	}
	// Build the request with the merge-duplicates Prefer header.
	b, err := json.Marshal(rows)
	if err != nil {
		return nil, err
	}
	u := c.cfg.url + "/rest/v1/" + table
	if len(q) > 0 {
		u += "?" + q.Encode()
	}
	req, err := http.NewRequest(http.MethodPost, u, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("apikey", c.cfg.serviceRoleKey)
	req.Header.Set("Authorization", "Bearer "+c.cfg.serviceRoleKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Prefer", "resolution=merge-duplicates,return=representation")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("upsert %s → %d: %s", table, resp.StatusCode, string(data))
	}
	return data, nil
}

// deleteWhere deletes rows matching an exact-match filter, e.g. {"entry_id":"eq.<id>"}.
func (c *client) deleteWhere(table string, filters url.Values) error {
	_, err := c.do(http.MethodDelete, "/rest/v1/"+table, filters, nil)
	return err
}

// selectIDBy returns the id column of the first row matching the filters.
func (c *client) selectIDBy(table string, filters url.Values) (string, error) {
	filters.Set("select", "id")
	filters.Set("limit", "1")
	data, err := c.do(http.MethodGet, "/rest/v1/"+table, filters, nil)
	if err != nil {
		return "", err
	}
	var rows []struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(data, &rows); err != nil {
		return "", err
	}
	if len(rows) == 0 {
		return "", nil
	}
	return rows[0].ID, nil
}
