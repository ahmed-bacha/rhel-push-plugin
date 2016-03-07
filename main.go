package main

import (
	"flag"

	"github.com/Sirupsen/logrus"
	"github.com/docker/go-plugins-helpers/authorization"
)

const (
	defaultDockerHost = "unix:///var/run/docker.sock"
	pluginSocket      = "/run/docker/plugins/rhel-push-plugin.sock"
)

var (
	flDockerHost = flag.String("host", defaultDockerHost, "Docker host the plugin connects to when inspecting")
	// TODO(runcom): add tls option to connect to docker?
	// TODO(runcom): add plugin tls option (need to learn more...)
	// TODO(runcom): add config tls option based on Dan's suggestion to block based on AuthN
)

func main() {
	flag.Parse()

	rhelpush, err := newPlugin(*flDockerHost)
	if err != nil {
		logrus.Fatal(err)
	}

	// TODO(runcom): parametrize this when the bin starts
	h := authorization.NewHandler(rhelpush)

	if err = h.ServeUnix("root", pluginSocket); err != nil {
		logrus.Fatal(err)
	}
}
