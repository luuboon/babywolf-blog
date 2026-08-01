package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// AdminCreateAuthUser creates a user directly in Supabase Auth using the
// service role key (Practica 11-12: alta de usuarios por el administrador).
// The public.users profile row is populated automatically by the
// on_auth_user_created trigger (see supabase_auth_trigger.sql).
func AdminCreateAuthUser(email, password, username string) (string, error) {
	body, _ := json.Marshal(map[string]any{
		"email":         email,
		"password":      password,
		"email_confirm": true,
		"user_metadata": map[string]string{"username": username},
	})

	resp, err := adminRequest(http.MethodPost, "/auth/v1/admin/users", body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return "", err
	}
	return out.ID, nil
}

// AdminSetUserPassword resets a user's password via the Supabase Auth admin API.
func AdminSetUserPassword(supabaseUID, newPassword string) error {
	body, _ := json.Marshal(map[string]any{"password": newPassword})
	resp, err := adminRequest(http.MethodPut, "/auth/v1/admin/users/"+supabaseUID, body)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return nil
}

func adminRequest(method, path string, body []byte) (*http.Response, error) {
	baseURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	req, err := http.NewRequest(method, baseURL+path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", serviceKey)
	req.Header.Set("Authorization", "Bearer "+serviceKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= 300 {
		defer resp.Body.Close()
		msg, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("supabase admin API error (%d): %s", resp.StatusCode, string(msg))
	}
	return resp, nil
}
