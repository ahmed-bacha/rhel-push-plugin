.PHONY: all binary man install clean
export GOPATH:=$(CURDIR)/Godeps/_workspace:$(GOPATH)

LIBDIR=${DESTDIR}/lib/systemd/system
BINDIR=${DESTDIR}/usr/libexec/docker/
PREFIX ?= ${DESTDIR}/usr
MANINSTALLDIR=${PREFIX}/share/man

all: man binary

binary:
	go build  -o rhel-push-plugin .

man:
	go-md2man -in man/rhel-push-plugin.8.md -out rhel-push-plugin.8

install:
	install -d -m 0755 ${LIBDIR}
	install -m 644 systemd/rhel-push-plugin.service ${LIBDIR}
	install -d -m 0755 ${LIBDIR}
	install -m 644 systemd/rhel-push-plugin.socket ${LIBDIR}
	install -d -m 0755 ${BINDIR}
	install -m 755 rhel-push-plugin ${BINDIR}
	install -m 644 rhel-push-plugin.8 ${MANINSTALLDIR}/man8/

clean:
	rm -f rhel-push-plugin
