package dborm

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

func RunDockerContainer(Image string, projectName string) (string, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())

	if err != nil {
		fmt.Println("error while creating docker client", err)
		return "", err
	}
	ctx := context.Background()
	imageExist, err := isImage(cli, ctx, Image)
	if err != nil {
		return "", err
	}
	if !imageExist {
		res, err := cli.ImagePull(ctx, Image, image.PullOptions{})
		if err != nil {
			fmt.Println("error is", err)
			return "", err
		}
		defer res.Close()
		_, err = io.Copy(os.Stdout, res)
		if err != nil {
			return "", err
		}
	}
	fmt.Println("image pull/found succesffuly")
	//create container with volume with same name as projectname
	dbString, err := createContainerWithVolume(cli, ctx, Image, projectName)
	if err != nil {
		return "", err
	}
	return dbString, nil
}
func createContainerWithVolume(cli *client.Client, ctx context.Context, Image string, projectName string) (string, error) {
	var dbString string
	if Image == "mongo" {
		res, err := cli.ContainerCreate(ctx, &container.Config{
			Image: Image,
			Env: []string{
				"MONGO_INITDB_ROOT_USERNAME=mongouser",
				"MONGO_INITDB_ROOT_PASSWORD=mongopassword",
				"MONGO_INITDB_DATABASE=mydb",
			},
			ExposedPorts: nat.PortSet{
				"27017/tcp": struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				"27017/tcp": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: "27017",
					},
				},
			},
			Mounts: []mount.Mount{{
				Type:   mount.TypeVolume,
				Source: projectName,
				Target: "/data/db"},
			},
		}, nil, nil, projectName)
		if err != nil {
			return "", err
		}
		if err = cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
			return "", err
		}
		dbString = "mongodb://mongouser:mongopassword@127.0.0.1:27017/mydb"
		fmt.Println("res from creating mongo container is", res)
	} else {
		res, err := cli.ContainerCreate(ctx, &container.Config{
			Image: Image,
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=password",
				"POSTGRES_DB=mydb",
			},
			ExposedPorts: nat.PortSet{
				"5432/tcp": struct{}{},
			},
		}, &container.HostConfig{
			PortBindings: nat.PortMap{
				"5432/tcp": []nat.PortBinding{
					{
						HostIP:   "127.0.0.1",
						HostPort: "5432",
					},
				},
			},
			Mounts: []mount.Mount{{
				Type:   mount.TypeVolume,
				Source: projectName,
				Target: "/var/lib/postgresql/data",
			}},
		}, nil, nil, projectName)
		if err != nil {
			return "", err
		}
		fmt.Println("res from Postgres is ", res)
		if err = cli.ContainerStart(ctx, res.ID, container.StartOptions{}); err != nil {
			return "", err
		}
		dbString = "postgresql://postgres:password@127.0.0.1:5432/mydb"
	}

	return dbString, nil
}
func isImage(cli *client.Client, ctx context.Context, Image string) (bool, error) {
	filter := filters.NewArgs()
	filter.Add("reference", Image)
	images, err := cli.ImageList(ctx, image.ListOptions{
		Filters: filter,
	})
	if err != nil {
		return false, err
	}
	return len(images) > 1, nil
}
