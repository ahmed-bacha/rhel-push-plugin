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
	flDockerHost   = flag.String("host", defaultDockerHost, "Specifies the host where to contact the docker daemon")
	flCertPath     = flag.String("cert-path", "", "Certificates path to connect to Docker (cert.pem, key.pem)")
	flTLSVerify    = flag.Bool("tls-verify", false, "Whether to verify certificates or not")
	flPluginSocket = flag.String("plugin-sock", pluginSocket, "Specifies the plugin socket path")
)

func main() {
	flag.Parse()

	rhelpush, err := newPlugin(*flDockerHost, *flCertPath, *flTLSVerify)
	if err != nil {
		logrus.Fatal(err)
	}

	h := authorization.NewHandler(rhelpush)

	if err := h.ServeUnix("root", *flPluginSocket); err != nil {
		logrus.Fatal(err)
	}
}
