package dockercli

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/bSkracic/delta-rest/lib/wrap"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type Dockercli struct {
	Cli *client.Client
	Ctx context.Context
}

func CreateClient() *Dockercli {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}

	return &Dockercli{Cli: cli, Ctx: context.Background()}
}

func (dc *Dockercli) RetreiveAvailableContainer(langName string, image string) string {
	// List containers with requested language in name, find the first suitable one.
	// If error occurred while restarting stopped containers and none suitable started containers where found: start a new one.
	contName := fmt.Sprintf("delta-%s", langName)
	filters := filters.NewArgs()
	filters.Add("name", contName)
	containers, err := dc.Cli.ContainerList(dc.Ctx, types.ContainerListOptions{Filters: filters, All: true})
	if err != nil {
		panic(err)
	}

	var id string
	if len(containers) == 0 {
		id = dc.Run(fmt.Sprintf("delta-%s-%v", langName, time.Now().Unix()), image)
	} else {
		found := false
		for _, c := range containers {
			if strings.Contains(c.Status, "Up") == false {
				if err := dc.Cli.ContainerStart(dc.Ctx, c.ID, types.ContainerStartOptions{}); err != nil {
					continue
				}
			}
			id = c.ID
			found = true
			break
		}
		if !found {
			id = dc.Run(fmt.Sprintf("delta-%s-%v", langName, time.Now().Unix()), image)
		}
	}
	return id
}

func (dc *Dockercli) Kill(id string) {
	dc.Cli.ContainerKill(dc.Ctx, id, "")
	log.Default().Printf("Killed: %s\n", id)
}

func (dc *Dockercli) Remove(id string) {
	dc.Cli.ContainerRemove(dc.Ctx, id, types.ContainerRemoveOptions{})
	log.Default().Printf("Removed: %s\n", id)
}

func (dc *Dockercli) Run(name string, image string) string {
	resp, err := dc.Cli.ContainerCreate(dc.Ctx, &container.Config{
		Image: image,
		Cmd:   []string{"sleep", "infinity"},
		Tty:   false,
	}, nil, nil, nil, name)
	if err != nil {
		panic(err)
	}

	if err := dc.Cli.ContainerStart(dc.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func (dc *Dockercli) CreateDir(id string, dirName string) error {
	cmd := fmt.Sprintf("mkdir %s", dirName)
	config := types.ExecConfig{Cmd: []string{"bash", "-c", cmd}}
	respCreate, err := dc.Cli.ContainerExecCreate(dc.Ctx, id, config)
	if err != nil {
		return err
	}

	// Listen for an event and return only after exec finishes
	msgs, errs := dc.Cli.Events(dc.Ctx, types.EventsOptions{})
	dc.Cli.ContainerExecStart(dc.Ctx, respCreate.ID, types.ExecStartCheck{})

	for {
		select {
		case err := <-errs:
			return err
		case msg := <-msgs:
			log.Printf("%v\n", msg)
			if msg.Action == "exec_die" && msg.Actor.ID == id {
				return nil
			}
		}
	}

}

func (dc *Dockercli) Copy(id string, fileName string, file []byte, path string) error {
	var buf bytes.Buffer
	buf.Write(file)
	f, err := wrap.Generate(fileName, buf.String())
	if err != nil {
		return err
	}

	return dc.Cli.CopyToContainer(dc.Ctx, id, path, f, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true})
}

type ExecOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (dc *Dockercli) Exec(id string, cmd string, ch chan<- *ExecOutput) {
	config := types.ExecConfig{AttachStdout: true, Cmd: []string{"bash", "-c", cmd}}
	respCreate, err := dc.Cli.ContainerExecCreate(dc.Ctx, id, config)
	if err != nil {
		log.Fatalln(err)
	}

	respExec, err := dc.Cli.ContainerExecAttach(dc.Ctx, respCreate.ID, types.ExecStartCheck{})
	defer respExec.Close()

	// Read the output
	var outBuf, errBuf bytes.Buffer
	outputDone := make(chan error)

	go func() {
		// StdCopy demultiplexes the stream into two buffers
		_, err = stdcopy.StdCopy(&outBuf, &errBuf, respExec.Reader)
		outputDone <- err
	}()

	select {
	case err := <-outputDone:
		if err != nil {
			log.Fatalln(err)
		}
		break

	case <-dc.Ctx.Done():
	}

	stdout, err := ioutil.ReadAll(&outBuf)
	if err != nil {
		log.Fatalln(err)
	}
	stderr, err := ioutil.ReadAll(&errBuf)
	if err != nil {
		log.Fatalln(err)
	}

	res, err := dc.Cli.ContainerExecInspect(dc.Ctx, respCreate.ID)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("\ncommand: %s\nstdout: %s\nstderr: %s\nexit_code: %v\n\n", cmd, stdout, stderr, res.ExitCode)

	ch <- &ExecOutput{Stdout: string(stdout), Stderr: string(stderr), ExitCode: res.ExitCode}
}
