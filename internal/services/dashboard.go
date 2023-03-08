package services

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
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

type DashboardData struct {
	Port     uint16 `json:"port"`
	Password string `json:"password"`
}

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

func (d *DashboardService) getPasswordFromContainer(name string) (string, error) {
	command := exec.Command("docker", "exec", name, "cat", "/password.txt")

	output, err := command.CombinedOutput()
	if err != nil {
		return "", err
	}

	password := strings.ReplaceAll(string(output), "\n", "")

	return password, nil
}

func (d *DashboardService) StartDashboard(ctx context.Context, userID uint64) (*DashboardData, error) {
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
		return nil, err
	}

	err = d.dockerClient.ContainerStart(ctx, response.ID, types.ContainerStartOptions{})
	if err != nil {
		return nil, err
	}

	containerData, err := d.dockerClient.ContainerInspect(ctx, response.ID)
	if err != nil {
		return nil, err
	}

	ports := make([]string, 0, len(containerData.HostConfig.PortBindings))
	for _, v := range containerData.HostConfig.PortBindings {
		ports = append(ports, v[0].HostPort)
	}

	port, _ := strconv.Atoi(ports[0])

	password, err := d.getPasswordFromContainer(containerData.Name)
	if err != nil {
		return nil, err
	}

	dashboardData := &DashboardData{
		Port:     uint16(port),
		Password: password,
	}

	return dashboardData, nil
}

func (d *DashboardService) GetDashboard(ctx context.Context, userID uint64) (*DashboardData, error) {
	containers, err := d.dockerClient.ContainerList(ctx, types.ContainerListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "name",
			Value: fmt.Sprintf("tit_dashboard_%d", userID),
		}),
		All: true,
	})
	if err != nil {
		return nil, err
	}

	if len(containers) >= 1 {
		startedContainer := containers[0]

		if !strings.Contains(startedContainer.Status, "Up") {
			return nil, errDashboardIsStopped
		}

		containerData, err := d.dockerClient.ContainerInspect(ctx, startedContainer.ID)
		if err != nil {
			return nil, err
		}

		password, err := d.getPasswordFromContainer(containerData.Name)
		if err != nil {
			return nil, err
		}

		dashboardData := &DashboardData{
			Port:     startedContainer.Ports[0].PublicPort,
			Password: password,
		}

		return dashboardData, nil
	}

	return nil, errDashboardIsStopped
}
