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

	contName := fmt.Sprintf("delta-%s", langName)

	filters := filters.NewArgs()
	filters.Add("name", contName)
	containers, err := dc.Cli.ContainerList(dc.Ctx, types.ContainerListOptions{Filters: filters})
	if err != nil {
		panic(err)
	}

	var id string

	if len(containers) == 0 {
		id = dc.Run(fmt.Sprintf("delta-%s-%v", langName, time.Now().Unix()), image)
	} else {
		found := false
		for _, c := range containers {
			if strings.Contains(c.Status, "Up") {
				id = c.ID
				found = true
				break
			} else {
				// Remove all containers that are exited or not running, may be slow, depending on number of containers!
				dc.Remove(c.ID)
			}
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

func (dc *Dockercli) Copy(mainFile []byte, mainFileName string, id string) {
	var buf bytes.Buffer
	buf.Write(mainFile)
	file, err := wrap.Generate(mainFileName, buf.String())
	if err != nil {
		log.Fatalln(err)
	}

	filePath := "/"
	if err := dc.Cli.CopyToContainer(dc.Ctx, id, filePath, file, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
		log.Fatalf("copy failed: %s", err)
	}
}

type ExecOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func (dc *Dockercli) Exec(id string, cmd string, ch chan<- *ExecOutput) {
	config := types.ExecConfig{AttachStdin: true, AttachStdout: true, Cmd: []string{"bash", "-c", cmd}}
	respCreate, err := dc.Cli.ContainerExecCreate(dc.Ctx, id, config)
	if err != nil {
		log.Fatalln(err)
	}

	respExec, err := dc.Cli.ContainerExecAttach(dc.Ctx, respCreate.ID, types.ExecStartCheck{})
	defer respExec.Close()

	// read the output
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

	log.Printf("\ncommand: %s\n\tstdout: %s\n\tstderr: %s\n\texit_code: %v\n\n", cmd, stdout, stderr, res.ExitCode)

	ch <- &ExecOutput{Stdout: string(stdout), Stderr: string(stderr), ExitCode: res.ExitCode}
}
