package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/bSkracic/delta-cli/dockercli"
	"github.com/bSkracic/delta-cli/wrap"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
)

// Prerequesites:
// 		docker
//		docker images: gcc:latest, soon: java:latest, python:latest
// TODO: 0. parse os.Args (-l lang, -t time, OPT: -c compiler_options -a program_args)
// TODO: 1. create docker client
// TODO: 2. check if corresponding execution environment container is running (docker ps):
// 			if not -> check if that container (docker ps -a) is stopped or exited:
//                if yes -> remove it (docker rm container_name);
//          	 finally start new container with infinite sleep (detached) (docker run -d container_name bash -c "sleep infinity")
// TODO: 3. copy main file to that container -> check if file exists! (docker cp file_name container_name:/ )
// TODO: 4. execute source code (docker exec container_name bash -c "g++ -o main ./main.c && ./main") as a goroutine
// TODO: 5. if arg:time was passed, start timer and kill docker exec process if timeout is reached before exec result
// TODO: Log performance of program and executon in ms
// OPT: create swarm of multiple containers with same program execution environment for each language (c, java, python)
// OPT: pass program cmdline arguments and/or compiler/interpreter arguments

var t = flag.Int("t", 0, "Execution time limit (ms)")
var lang = flag.String("l", "c", "Language (c, java, python)")

func init() {
	// 0
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: [options] [file]\n")
		flag.PrintDefaults()
	}
}

func main() {

	start := time.Now()

	flag.Parse()

	if len(os.Args) < 2 {
		log.Fatalln("You need to provide main file")
	} else if _, err := os.Stat(os.Args[len(os.Args)-1]); errors.Is(err, os.ErrNotExist) {
		log.Fatalf("File %s does not exist\n", os.Args[1])
	}

	var lc *dockercli.LanguageConf

	switch *lang {
	case "c":
		lc = &dockercli.LanguageConf{Name: "c", Compiler: "gcc", Extension: "c", Cmd: "gcc -o main ./main.c && ./main", Image: "gcc:4.9", File: "main.c"}
	case "java":
		lc = &dockercli.LanguageConf{Name: "java", Compiler: "javac", Extension: "java", Cmd: "javac Main.java && java Main", Image: "openjdk:latest", File: "Main.java"}
	case "python":
		lc = &dockercli.LanguageConf{Name: "python", Compiler: "python3", Extension: "py", Cmd: "python3 main.py", Image: "python:latest", File: "main.py"}
	default:
		log.Fatalf("Unsupported lang: %s\n", *lang)
	}

	if *t < 0 {
		log.Fatalf("Invalid value for time limit: %d\n", *t)
	}

	filePath := os.Args[len(os.Args)-1]
	conName := fmt.Sprintf("delta-%s", *lang)

	dcli := dockercli.CreateClient()

	id := ""

	// List containers and find the suitable one for execution environment
	filters := filters.NewArgs()
	filters.Add("name", conName)
	containers, err := dcli.Cli.ContainerList(dcli.Ctx, types.ContainerListOptions{Filters: filters})
	if err != nil {
		panic(err)
	}

	if len(containers) == 0 {
		id = dcli.Run(lc)
	} else {
		found := false
		for _, c := range containers {
			if strings.Contains(c.Status, "Up") {
				id = c.ID
				found = true
				break
			} else {
				// Remove all containers that are exited or not running, may be slow, depending on number of containers!
				dcli.Remove(c.ID)
			}
		}

		if !found {
			id = dcli.Run(lc)
		}
	}

	// Prepare file in tar archive
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalln(err)
	}

	var buf bytes.Buffer
	buf.Write(content)
	file, err := wrap.Generate(lc.File, buf.String())
	if err != nil {
		log.Fatalln(err)
	}

	dcli.Copy(file, id)

	// Start another goroutine to execute commands in container
	ch := make(chan string, 1)

	startExec := time.Now()
	go dcli.Exec(id, lc, ch)

	if *t != 0 {
		select {
		case res := <-ch:
			fmt.Printf("Finished :)\n%s", res)
		case <-time.After(time.Duration(*t) * time.Millisecond):
			fmt.Println("Interrupted :/")
		}
	} else {
		fmt.Printf("Finished :)\n%s", <-ch)
	}

	// OPT: Meassurment
	fmt.Printf("Elapsed exec time: %dms\n", time.Now().Sub(startExec).Milliseconds())
	fmt.Printf("Elapsed total time: %dms\n", time.Now().Sub(start).Milliseconds())
}
