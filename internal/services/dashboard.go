package services

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var errDashboardIsStopped = errors.New("dashboard is stopped")

type DashboardService struct {
	dockerClient *client.Client
}

func NewDashboardService() (*DashboardService, error) {
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &DashboardService{dockerClient}, nil
}

func (d *DashboardService) StartDashboard(ctx context.Context, userID uint64) (int, error) {
	response, err := d.dockerClient.ContainerCreate(
		ctx,
		&container.Config{
			Image:        "tit-dashboard:latest",
			ExposedPorts: nat.PortSet{"5900": {}},
		},
		&container.HostConfig{
			AutoRemove: true,
			PortBindings: nat.PortMap{
				"5900": {{HostIP: "0.0.0.0", HostPort: "0"}},
			},
		},
		&network.NetworkingConfig{},
		nil,
		fmt.Sprintf("tit_dashboard_%d", userID),
	)
	if err != nil {
		return 0, err
	}

	err = d.dockerClient.ContainerStart(ctx, response.ID, types.ContainerStartOptions{})
	if err != nil {
		return 0, err
	}

	containerData, err := d.dockerClient.ContainerInspect(ctx, response.ID)
	if err != nil {
		return 0, err
	}

	ports := make([]string, 0, len(containerData.HostConfig.PortBindings))
	for _, v := range containerData.HostConfig.PortBindings {
		ports = append(ports, v[0].HostPort)
	}

	port, _ := strconv.Atoi(ports[0])

	return port, nil
}

func (d *DashboardService) GetDashboardPort(ctx context.Context, userID uint64) (int, error) {
	containers, err := d.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: fmt.Sprintf("tit_dashboard_%d", userID),
		}),
		All: true,
	})
	if err != nil {
		return 0, err
	}

	if len(containers) >= 1 {
		startedContainer := containers[0]

		if !strings.Contains(startedContainer.Status, "Up") {
			return 0, errDashboardIsStopped
		}

		port := startedContainer.Ports[0].PublicPort

		return int(port), nil
	}

	return 0, nil
}
