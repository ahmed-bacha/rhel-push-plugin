% RHEL-PUSH-PLUGIN(8)
% Antonio Murdaca
% MARCH 2016
# NAME
rhel-push-plugin - Blocks RHEL content push to docker.io

# SYNOPSIS
**rhel-push-plugin**
[**--cert-path**=[=*""*]]
[**--host**=[=*unix:///var/run/docker.sock*]]
[**--tls-verify**=[=*false*]]

# DESCRIPTION

This plugin avoids any RHEL based image to be pushed to the default `docker.io` registry preventing
users to violate the RH subscription agreement.

# OPTIONS

**--cert-path**=""
  Certificates path to connect to Docker (cert.pem, key.pem)
**--host**="unix:///var/run/docker.sock"
  Specifies the host where to contact the docker daemon.
**--tls-verify**="false"
  Whether to verify certificates or not

# AUTHORS
Antonio Murdaca <runcom@redhat.com>
