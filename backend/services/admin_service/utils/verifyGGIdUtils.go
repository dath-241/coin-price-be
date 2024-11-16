package utils

import (
    "fmt"
    "net/http"
    "encoding/json"
)

func VerifyGoogleIDToken(idToken string) (map[string]interface{}, error) {
    url := fmt.Sprintf("https://oauth2.googleapis.com/tokeninfo?id_token=%s", idToken)
    resp, err := http.Get(url)
    if err != nil {
        return nil, fmt.Errorf("failed to reach tokeninfo endpoint: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("invalid token: status %d", resp.StatusCode)
    }

    var result map[string]interface{}
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode tokeninfo response: %v", err)
    }

    return result, nil
}