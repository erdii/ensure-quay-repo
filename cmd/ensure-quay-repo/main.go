package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"
)

func main() {
	authToken, ok := os.LookupEnv("AUTH_TOKEN")
	if !ok {
		fmt.Fprintln(os.Stderr, "AUTH_TOKEN must be configured")
		os.Exit(1)
	}

	// validate command arguments and print usage if needed
	if len(os.Args) != 3 || os.Args[1] == "" || os.Args[2] == "" {
		exe, err := os.Executable()
		if err != nil {
			panic(fmt.Errorf("resolving executable when printing usage because arg is missing: %w", err))
		}
		fmt.Fprintf(os.Stderr, "usage: %s <org> <repo name>\n", exe)
		os.Exit(1)
	}
	org := os.Args[1]
	repo := os.Args[2]

	// instanciate client and ensure repo
	rc := newRepoClient(authToken)
	if err := rc.ensure(org, repo); err != nil {
		panic(err)
	}
	os.Exit(0)
}

type repoClient struct {
	*http.Client

	apiUrl    string
	authToken string
	logger    *slog.Logger
}

func newRepoClient(authToken string) *repoClient {
	return &repoClient{
		Client: &http.Client{
			Timeout: 5 * time.Second,
		},
		apiUrl:    "https://quay.io/api/v1/",
		authToken: authToken,
		logger:    slog.Default(),
	}
}

func (rc *repoClient) ensure(org, name string) error {
	rc.logger.Info("Ensuring Repo", "org", org, "name", name)

	exists, err := rc.exists(fmt.Sprintf("%s/%s", org, name))
	if err != nil {
		return err
	}

	if exists {
		rc.logger.Info("Repo already exists", "org", org, "name", name)
		return nil
	}

	rc.logger.Info("Repo does not exist. Creating Repo", "org", org, "name", name)
	return rc.create(org, name)
}

func (rc *repoClient) exists(namespacedName string) (bool, error) {
	rc.logger.Info("Checking if Repo exists", "namespacedName", namespacedName)

	req, err := rc.newRequest(http.MethodGet, fmt.Sprintf(
		"/repository/%s?includeTags=false&includeStats=false",
		namespacedName,
	), nil)
	if err != nil {
		return false, err
	}

	resp, err := rc.Do(req)
	if err != nil {
		return false, fmt.Errorf("checking if repo exists: %w", err)
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200, nil
}

func (rc *repoClient) create(org, name string) error {
	values := map[string]string{
		"repo_kind":   "image",
		"namespace":   org,
		"visibility":  "public",
		"repository":  name,
		"description": "",
	}

	marshaled, err := json.Marshal(values)
	if err != nil {
		return fmt.Errorf("marshaling json data: %w", err)
	}

	req, err := rc.newRequest(http.MethodPost, "/repository", bytes.NewReader(marshaled))
	if err != nil {
		return err
	}

	resp, err := rc.Do(req)
	if err != nil {
		return fmt.Errorf("creating repo: %w", err)
	}

	if resp.StatusCode != 201 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("reading body: %w", err)
		}
		return fmt.Errorf("statuscode was not 201: %d\n%s", resp.StatusCode, string(b))
	}
	return nil
}

func (rc *repoClient) newRequest(method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(
		method,
		fmt.Sprintf("%s%s", rc.apiUrl, path),
		body,
	)
	if err != nil {
		return nil, fmt.Errorf("preparing request: %s", err)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+rc.authToken)

	return req, nil
}
