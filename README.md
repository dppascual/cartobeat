# Cartobeat

Cartobeat is a beat based on metricbeat which was generated with metricbeat/metricset generator.

## Requirements

- [Golang](https://golang.org/doc/install)
- [Docker](https://www.docker.com/products/docker-engine)

## Getting started

To get started run the following command. 

```
mkdir -p $GOPATH/src/github.com/dppascual
git clone git clone git@github.com:dppascual/cartobeat.git
go get -u github.com/kardianos/govendor
cd $GOPATH/src/github.com/dppascual/cartobeat
govendor sync
```

This will fech all dependencies and put them into a directory `vendor` inside your repository. [govendor](https://github.com/kardianos/govendor) is used for vendoring.

Copy `Makefiles` that are needed to build the package:

```
go get -u github.com/elastic/beats
make copy-vendor
```

## Packaging

The beat frameworks provides tools to crosscompile and package your beat for different platforms. This requires [docker](https://www.docker.com/) and vendoring as described above. To build packages of your beat, run the following command:

```
go get -u github.com/magefile/mage
make release
```

This will fetch and create all images required for the build process. The whole process to finish can take several minutes.
