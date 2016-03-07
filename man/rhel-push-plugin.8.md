% RHEL-PUSH-PLUGIN(8)
% Antonio Murdaca
% MARCH 2016
# NAME
rhel-push-plugin - Blocks RHEL content push to docker.io

# SYNOPSIS
**rhel-push-plugin**
[**--host**=[=*unix:///var/run/docker.sock*]]

# DESCRIPTION

This plugin avoids any RHEL based image to be pushed to the default `docker.io` registry preventing
users to violate the RH subscription agreement.

# OPTIONS

**--host**="unix:///var/run/docker.sock"
  Specifies the host where to contact the docker daemon.

# AUTHORS
Antonio Murdaca <runcom@redhat.com>
