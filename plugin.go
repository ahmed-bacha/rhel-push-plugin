package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"path/filepath"
	"regexp"
	"strings"

	dockerapi "github.com/docker/docker/api"
	"github.com/docker/docker/reference"
	dockerclient "github.com/docker/engine-api/client"
	"github.com/docker/engine-api/types"
	"github.com/docker/go-plugins-helpers/authorization"
)

func newPlugin(dockerHost, certPath string, tlsVerify bool) (*rhelpush, error) {
	var transport *http.Transport
	if certPath != "" {
		tlsc := &tls.Config{}

		cert, err := tls.LoadX509KeyPair(filepath.Join(certPath, "cert.pem"), filepath.Join(certPath, "key.pem"))
		if err != nil {
			return nil, fmt.Errorf("Error loading x509 key pair: %s", err)
		}

		tlsc.Certificates = append(tlsc.Certificates, cert)
		tlsc.InsecureSkipVerify = !tlsVerify
		transport = &http.Transport{
			TLSClientConfig: tlsc,
		}
	}

	client, err := dockerclient.NewClient(dockerHost, dockerapi.DefaultVersion.String(), transport, nil)
	if err != nil {
		return nil, err
	}
	return &rhelpush{client: client}, nil
}

var (
	pushRegExp = regexp.MustCompile(`/images/(.*)/push(\?tag=(.*))?$`)
)

const (
	RHELVendorLabel     = "Red Hat, Inc."
	RHELNameLabelPrefix = "rhel"
)

type rhelpush struct {
	client *dockerclient.Client
}

func (p *rhelpush) AuthZReq(req authorization.Request) authorization.Response {
	if req.RequestMethod == "POST" && pushRegExp.MatchString(req.RequestURI) {
		res := pushRegExp.FindStringSubmatch(req.RequestURI)
		if len(res) < 3 {
			return authorization.Response{Err: "unable to find repository name and reference"}
		}

		repoName := res[1]
		tag := res[3]

		if tag == "" {
			imgListOptions := types.ImageListOptions{}
			imgListOptions.MatchName = repoName
			images, err := p.client.ImageList(imgListOptions)
			if err != nil {
				return authorization.Response{Err: err.Error()}
			}
			for _, img := range images {
				rb, err := p.isRHELBased(img.ID)
				if err != nil {
					return authorization.Response{Err: err.Error()}
				}
				if rb {
					return authorization.Response{Err: fmt.Sprintf("%s is RHEL based (refers to: %s), please push your image using a tag directly to avoid pushing the RHEL based image", img.ID, strings.Join(img.RepoTags, ", "))}
				}
			}
			goto allow
		}

		repoName = fmt.Sprintf("%s:%s", repoName, tag)
		RHELBased, err := p.isRHELBased(repoName)
		if err != nil {
			return authorization.Response{Err: err.Error()}
		}
		if !RHELBased {
			goto allow
		}

		// any direct push to docker.io/ with a qualified image is rejected
		if strings.HasPrefix(repoName, "docker.io/") {
			goto noallow
		}

		ref, err := reference.ParseNamed(repoName)
		if err != nil {
			return authorization.Response{Err: err.Error()}
		}
		// ref.Hostname() uses the docker/docker/reference implementation, which automatically
		// maps unspecified hostname to reference.DefaultHostname.
		// Due to the strings.HasPrefix() check above, we now know that the repository name does
		// not contain a host name so it's not a direct push to docker.io.
		//
		// This `if` will match pushing *unqualified* images to the default registry
		// with the projectatomic/docker codebase and the docker official binary.
		if ref.Hostname() == "docker.io" {
			registries, err := p.getAdditionalDockerRegistries()
			if err != nil {
				return authorization.Response{Err: err.Error()}
			}
			if len(registries) != 0 {
				// We have a projectatomic/docker implementation: pushing without specifying a host name
				// automatically uses the first just discovered registry configured with --add-registry
				// If the first registry configured in the daemon is docker.io blocks.
				if registries[0] == "docker.io" {
					goto noallow
				}
				// otherwise let the user push rhel content to his first configured registry.
				goto allow
			}
			// this is the official docker binary case
			goto noallow
		}
	}
allow:
	return authorization.Response{Allow: true}

noallow:
	return authorization.Response{Msg: "RHEL based images are not allowed to be pushed to docker.io"}
}

func (p *rhelpush) AuthZRes(req authorization.Request) authorization.Response {
	return authorization.Response{Allow: true}
}

// TODO(runcom): official engine-api client doesn't have Registries
// hacked into Godeps/_workspace/src/github.com/docker/engine-api/types/types.go
func (p *rhelpush) getAdditionalDockerRegistries() ([]string, error) {
	i, err := p.client.Info()
	if err != nil {
		return nil, err
	}
	regs := []string{}
	for _, r := range i.Registries {
		regs = append(regs, r.Name)
	}
	return regs, nil
}

func (p *rhelpush) isRHELBased(repoName string) (bool, error) {
	for {
		if repoName == "" {
			return false, nil
		}
		image, _, err := p.client.ImageInspectWithRaw(repoName, false)
		if err != nil {
			return false, err
		}
		if image.Config.Labels["Vendor"] == RHELVendorLabel && strings.HasPrefix(image.Config.Labels["Name"], RHELNameLabelPrefix) {
			return true, nil
		}
		repoName = image.Parent
	}
}
