Block RHEL push to docker.io
=
_In order to use this plugin you must be running at least Docker 1.10 which
has support for authorization plugins._

TODO

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
Prerequisites (replace `runcom` with your own Docker Hub username):
```
$ sudo dnf install rhel-push-plugin
$ sudo systemctl start rhel-push-plugin
# edit /etc/sysconfig/docker and append --authorization-plugin=rhel-push-plugin to OPTIONS
$ sudo systemctl restart docker

$ docker build -t runcom/testrhelbased - <<EOF
FROM registry.access.redhat.com/rhel7
EOF
$ docker tag runcom/testrhelbased docker.io/runcom/testrhelbased
$ docker pull busybox
$ docker tag busybox runcom/busybox
$ docker tag runcom/busybox docker.io/runcom/busybox
$ docker login -u runcom
```
Without any `--add-regsitry` configured:
```
$ docker push runcom/testrhelbased # blocked
$ docker push docker.io/runcom/testrhelbased # blocked
$ docker push runcom/busybox # works
$ docker push docker.io/runcom/busybox # works
```
With `--add-registry=yourownadditionalregistry:5000`:
```
$ docker push runcom/testrhelbased # works
$ docker push docker.io/runcom/testrhelbased # blocked
$ docker push runcom/busybox # works
$ docker push docker.io/runcom/busybox # works
```
License
-
MIT
