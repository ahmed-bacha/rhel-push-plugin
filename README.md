Block RHEL push to docker.io
=
_In order to use this plugin you must be running at least Docker 1.10 which
has support for authorization plugins._

This plugin avoids any RHEL based image to be pushed to the default `docker.io` registry preventing
users to violate the RH subscription agreement.

Building
-
```sh
$ export GOPATH=~ # optional if you already have this
$ mkdir -p ~/src/github.com/projectatomic # optional, from now on I'm assuming GOPATH=~
$ cd ~/src/github.com/projectatomic && git clone https://github.com/projectatomic/rhel-push-plugin
$ cd rhel-push-plugin
$ make
```
Installing
-
```sh
$ sudo make install
$ systemctl enable rhel-push-plugin
```
Running
-
Specify `--authorization-plugin=rhel-push-plugin` in the `docker daemon` command line
flags (either in the systemd unit file or in `/etc/sysconfig/docker` in `$OPTIONS`
or when manually starting the daemon).
The plugin must be started before `docker` (done automatically via systemd unit file).
If you're not using the systemd unit file:
```sh
$ rhel-push-plugin &
```
Just restart `docker` and you're good to go!
Systemd socket activation
-
The plugin can be socket activated by systemd. You just have to basically use the file provided
under `systemd/` (or installing via `make install`). This ensures the plugin gets activated
if it goes down for any reason.
How to test
-
Given the plugin is enabled:

- case docker/docker daemon

  - if the image is not rhel based and qualified-> allow pushing
  - if the image is not rhel based and unqualified -> allow pushing
  - if the image is rhel based and qualified with docker.io -> disallow pushing
  - if the image is rhel based and qualified with myregistry.com:5000 -> allow pushing
  - if the image is rhel based and unqualified -> disallow pushing

- case projectatomic/docker daemon with additional registries REST route and 1 additional registry configured at myregistry.com:5000

if the image is not rhel based and qualified-> allow pushing
if the image is not rhel based and unqualified -> allow pushing (it goes to myregistry.com:5000)
if the image is rhel based and qualified with docker.io -> disallow pushing **
if the image is rhel based and qualified with myregistry.com:5000 -> allow pushing
if the image is rhel based and unqualified -> allow pushing (it goes to myregistry.com:5000)

- case projectatomic/docker daemon with additional registries REST route and NO additional registry configured which implies registries[0] == docker.io
case docker/docker daemon

if the image is not rhel based and qualified-> allow pushing
if the image is not rhel based and unqualified -> allow pushing
if the image is rhel based and qualified with docker.io -> disallow pushing
if the image is rhel based and qualified with myregistry.com:5000 -> allow pushing
if the image is rhel based and unqualified -> disallow pushing
License
-
MIT
