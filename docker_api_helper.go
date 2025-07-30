package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"

	"pidanker-client/configs"
)

type Container struct {
	ID      string   `json:"Id"`
	Names   []string `json:"Names"`
	Image   string   `json:"Image"`
	ImageID string   `json:"ImageID"`
	Command string   `json:"Command"`
	Created int64    `json:"Created"`
	State   string   `json:"State"`  // 例如 "running", "exited", "created"
	Status  string   `json:"Status"` // 例如 "Up 2 hours", "Exited (0) 2 days ago"
}

func GetContainers() ([]Container, error) {
	dockerClient := http.Client{
		Transport: &http.Transport{
			DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
				return net.Dial("unix", configs.DockerSocketPath)
			},
		},
	}

	statsURL := fmt.Sprintf("http://%s/containers/json?all=true", configs.DockerAPIVersion)

	req, err := http.NewRequest("GET", statsURL, nil)
	if err != nil {
		log.Fatalf("Could not create request to Docker API: %v", err)
	}

	resp, err := dockerClient.Do(req)
	if err != nil {
		log.Fatalf("Could not connect to Docker API. Is Docker running? Error: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Docker API Resp Error: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Read resp body Error: %w", err)
	}

	var containers []Container

	if err := json.Unmarshal(body, &containers); err != nil {
		return nil, fmt.Errorf("Resp body json unmarshal Error: %w", err)
	}

	return containers, nil
}
