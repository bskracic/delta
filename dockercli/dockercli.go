package dockercli

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/docker/pkg/stdcopy"
)

type LanguageConf struct {
	Name      string
	Compiler  string
	Extension string
	File      string
	Cmd       string
	Image     string
}

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

func (dc *Dockercli) Kill(id string) {
	dc.Cli.ContainerKill(dc.Ctx, id, "")
	log.Default().Printf("Killed: %s\n", id)
}

func (dc *Dockercli) Remove(id string) {
	dc.Cli.ContainerRemove(dc.Ctx, id, types.ContainerRemoveOptions{})
	log.Default().Printf("Removed: %s\n", id)
}

func (dc *Dockercli) Run(lc *LanguageConf) string {
	conName := "delta-" + lc.Name
	resp, err := dc.Cli.ContainerCreate(dc.Ctx, &container.Config{
		Image: lc.Image,
		Cmd:   []string{"sleep", "infinity"},
		Tty:   false,
	}, nil, nil, nil, conName)
	if err != nil {
		panic(err)
	}

	if err := dc.Cli.ContainerStart(dc.Ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	return resp.ID
}

func (dc *Dockercli) Copy(file io.Reader, id string) {
	filePath := "/"
	if err := dc.Cli.CopyToContainer(dc.Ctx, id, filePath, file, types.CopyToContainerOptions{AllowOverwriteDirWithFile: true}); err != nil {
		log.Fatalf("copy failed: %s", err)
	}
}

func (dc *Dockercli) Exec(id string, lc *LanguageConf, ch chan<- string) {
	// fmt.Printf("cmd: %s\n", lf.cmd)
	config := types.ExecConfig{AttachStdin: true, AttachStdout: true, Cmd: []string{"bash", "-c", lc.Cmd}}
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

	ch <- fmt.Sprintf("STDOUT:\n%s\nSTDERR:\n%s\nEXIT CODE: %d\n", string(stdout), string(stderr), res.ExitCode)
}
