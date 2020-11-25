package main

import (
	"bufio"
	"context"
	"errors"
	"log"
	"os"
	"regexp"

	"github.com/docker/docker/client"
	"github.com/ndeloof/lancelot"
	"github.com/urfave/cli"
)

func main() {
	var cgroup string

	app := &cli.App{
		Name:  "lancelot",
		Usage: "Docker API proxy for your safety",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "cgroup",
				Usage:       "Define parent cgroup to assign all resources (defaults to own cgroup)",
				Destination: &cgroup,
			},
		},
		Action: func(c *cli.Context) error {
			return run(cgroup)
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cgroup string) error {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}
	cli.NegotiateAPIVersion(context.Background())

	if cgroup == "" {
		id, err := selfCGroup()
		if err != nil {
			return err
		}
		cgroup = id
	}

	proxy := lancelot.NewProxy(cli, cgroup)
	return proxy.Serve(":2375")
}

// selfCgroup retrieve lancelot's own cgroup and use it to restrict docker resources created through proxy
func selfCGroup() (string, error) {
	inFile, _ := os.Open("/proc/self/cgroup")
	defer inFile.Close()

	pids := regexp.MustCompile(`^[0-9]+:pids:/docker/([0-9a-z]+)`)
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		me := pids.FindStringSubmatch(scanner.Text())
		if len(me) > 0 {
			return me[1], nil
		}
	}
	return "", errors.New("can't detect cgroup")
}
