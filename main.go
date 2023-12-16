package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"log"
	"net"
	"github.com/skratchdot/open-golang/open"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
	"github.com/docker/go-connections/nat"
)

func portBinding(ip, port string) nat.PortMap {
  binding := nat.PortBinding {
    HostIP: ip,
    HostPort: port,
  }
  bindingMap := map[nat.Port][]nat.PortBinding{nat.Port(fmt.Sprintf("%s/tcp", port)): {binding}}
  return nat.PortMap(bindingMap)
}
func main() {
  ctx := context.Background()

  cli, err := client.NewClientWithOpts(client.FromEnv)
  if err != nil {
    panic(err)
  }
  defer cli.Close()

  images, err := cli.ImageList(ctx, types.ImageListOptions{})
  if err != nil {
    panic(err)
  }

  reader, err := cli.ImagePull(ctx, "consol/rocky-xfce-vnc", types.ImagePullOptions{})
  if err != nil {
    panic(err)
  }

  defer reader.Close()

  for _, image := range images {
    if len(image.RepoTags) > 0 {
      fmt.Printf("%s %s %d Mb\n", image.RepoTags[0], image.ID, image.Size / 1000000)
    } else {
      fmt.Printf("%s %d Mb\n", image.ID, image.Size / 1000000)

    }
  }

  io.Copy(os.Stdout, reader)
  
  resp, err := cli.ContainerCreate(ctx, &container.Config{
    Image: "consol/rocky-xfce-vnc",
    ExposedPorts: nat.PortSet{"5901/tcp":struct{}{},
    "6901/tcp":struct{}{}},
  },
  &container.HostConfig{PortBindings: nat.PortMap{
    "6901/tcp": []nat.PortBinding{
      {
        HostIP:"0.0.0.0",
        HostPort: "6901",
      },
    },
    "5901/tcp": []nat.PortBinding{
      {
        HostIP:"0.0.0.0",
        HostPort: "5901",
      },
    },
  }}, nil, nil, "")

  if err != nil {
    panic(err)
  }
  if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
    panic(err)
  }

  statusCh, errCh := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
  select {
  case err := <-errCh:
    if err != nil {
      panic(err)
    }
  case <-statusCh:
  }
  
  out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
  if err != nil {
    panic(err)
  }

  stdcopy.StdCopy(os.Stdout, os.Stderr, out)



	_, err = net.Listen("tcp", "localhost:5901")
	if err != nil {
		    log.Fatal(err)
		  }

		  // The browser can connect now because the listening socket is open.

		  err = open.Start("http://localhost:6901/test")
		  if err != nil {
		   	 log.Println(err)
		 	 }

		 	 // Start the blocking server loop.

		 	 //log.Fatal(http.Serve(l, r)) 

}
