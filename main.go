package main

import (
	"fmt"
	"strings"
	"time"

	"pidanker-client/logger"
)

type DockerStats struct {
	CpuStats    CpuStats    `json:"cpu_stats"`
	PreCpuStats CpuStats    `json:"precpu_stats"`
	MemoryStats MemoryStats `json:"memory_stats"`
	Name        string      `json:"name"`
	ID          string      `json:"id"`
}

type CpuStats struct {
	CpuUsage struct {
		TotalUsage uint64 `json:"total_usage"`
	} `json:"cpu_usage"`
	SystemCpuUsage uint64 `json:"system_cpu_usage"`
	OnlineCpus     int    `json:"online_cpus"`
}

type MemoryStats struct {
	Usage uint64 `json:"usage"`
	Limit uint64 `json:"limit"`
}

func main() {

	logger.InitLogger()

	logger.Logger.Info("Starting PiDanKer!")

	containers, err := GetContainers()
	if err != nil {
		logger.Logger.Error("Failed to get containers:", err)
		return
	}

	if len(containers) <= 0 {
		logger.Logger.Error("No containers found")
		return
	}

	for _, container := range containers {
		cleanName := strings.TrimPrefix(container.Names[0], "/")

		createdTime := time.Unix(container.Created, 0).Format(time.RFC3339)

		fmt.Println("--------------------------------------------------")
		fmt.Printf("容器名称: %s\n", cleanName)
		fmt.Printf("容器 ID:  %s\n", container.ID[:12])
		fmt.Printf("镜像:     %s\n", container.Image)
		fmt.Printf("状态:     %s\n", container.State)
		fmt.Printf("详细状态: %s\n", container.Status)
		fmt.Printf("创建时间: %s\n", createdTime)
		fmt.Println("--------------------------------------------------")
	}

}
