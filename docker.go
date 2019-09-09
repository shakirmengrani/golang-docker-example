package main

import (
	"context"
	"os"
	"fmt"
    "strings"

    "github.com/docker/docker/client"
    "github.com/docker/docker/api/types"
    "github.com/docker/docker/api/types/container"
    "github.com/docker/docker/pkg/stdcopy"
)

func startContainer(ctx context.Context,cli *client.Client, id string){
	if id == "" {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			panic(err)
		}
		for _, container := range containers {
			if strings.HasPrefix(container.Names[0], "/test") {
				if err := cli.ContainerStart(ctx, container.ID[:12], types.ContainerStartOptions{}); err != nil {
					panic(err)
				}
			}
			
		}
	}else {
		if err := cli.ContainerStart(ctx, id, types.ContainerStartOptions{}); err != nil {
			panic(err)
		}
	}
}

func createContainer(ctx context.Context,cli *client.Client){
    // create an container
    resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: "golang",
		Tty: true,
		WorkingDir: "/code",
        Cmd:   []string{"bash"},
    }, nil, nil, "test-go-docker")
    if err != nil {
        panic(err)
	}
	
 	// start created container
	 startContainer(ctx, cli, resp.ID)
	
    // waiting for running up
	if false {
		statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
		select {
		case err := <-errCh:
			if err != nil {
				panic(err)
			}
		case <-statusCh:
		}
	}	
    // getting output for making decisions
	if false {
		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			panic(err)
		}
		stdcopy.StdCopy(os.Stdout, os.Stderr, out)
	}
}

func stopContainer(ctx context.Context, cli *client.Client, id string){
	if id == "" {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			panic(err)
		}
		for _, container := range containers {
			if strings.HasPrefix(container.Names[0], "/test") {
				err = cli.ContainerStop(ctx, container.ID[:12], nil)
				if err != nil {
					panic(err)
				}
			}
			
		}
	}else {
		err := cli.ContainerStop(ctx, id, nil)
		if err != nil {
			panic(err)
		}
	}
}

func removeContainer(ctx context.Context, cli *client.Client, opts types.ContainerRemoveOptions, id string){
	if id == "" {
		containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
		if err != nil {
			panic(err)
		}
		for _, container := range containers {
			if strings.HasPrefix(container.Names[0], "/test") {
				err := cli.ContainerRemove(ctx, container.ID[:12], opts)
				if err != nil {
					panic(err)
				}
			}
			
		}
	}else {
		err := cli.ContainerRemove(ctx, id, opts)
		if err != nil {
			panic(err)
		}
	}
}

func listContainer(ctx context.Context, cli *client.Client, opts types.ContainerListOptions){
    containers, err := cli.ContainerList(ctx, opts)
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
        fmt.Printf("ID: {%s} Image: {%s} Name: {%s} Status: {%s}\n", container.ID[:12], container.Image, container.Names[0], container.Status)
	}
}

func Run(cmd string, id string){
	ctx := context.Background()
    cli, err := client.NewClientWithOpts(client.FromEnv)
    if err != nil {
        panic(err)
    }
    cli.NegotiateAPIVersion(ctx)
    if cmd == "create" {
		createContainer(ctx, cli)
	}else if cmd == "start" {
		startContainer(ctx, cli, id)
	}else if cmd == "stop" {
		stopContainer(ctx, cli, id)
	}else if cmd == "list" {
		listContainer(ctx, cli, types.ContainerListOptions{All: true})
	}else if cmd == "remove" {
		removeContainer(ctx, cli, types.ContainerRemoveOptions{Force: true, RemoveLinks: false, RemoveVolumes: false}, id)
	} else {
		fmt.Println(fmt.Sprintf("%s not valid", cmd))
	} 
}