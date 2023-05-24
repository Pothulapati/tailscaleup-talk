package tailscale

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"golang.org/x/oauth2/clientcredentials"
)

var (
	tsKey struct {
		Key string `json:"key"`
	}
)

func GetTodoAuthKeyFromEnv() (string, error) {
	// use tailscale oauth client
	var oauthConfig = &clientcredentials.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		TokenURL:     "https://api.tailscale.com/api/v2/oauth/token",
	}

	tailnet, ok := os.LookupEnv("TAILNET")
	if !ok {
		return "", fmt.Errorf("TAILNET env var not set")
	}

	// todo: add tags first
	client := oauthConfig.Client(context.Background())
	reqBody := `{
		"capabilities": {
		  "devices": {
			"create": {
			  "reusable": false,
			  "ephemeral": true,
			  "preauthorized": false,
			  "tags": [ "tag:tailtodo" ]
			}
		  }
		},
		"expirySeconds": 86400
	  }`

	resp, err := client.Post(fmt.Sprintf("https://api.tailscale.com/api/v2/tailnet/%s/keys", tailnet), "application/json", strings.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("error getting keys: %v", err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	// convert body into tsKey struct
	err = json.Unmarshal(body, &tsKey)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling response body: %v", err)
	}

	return tsKey.Key, nil
}
