```
                                                 ,-.
                                                  ) \
                                              .--'   |
                                             /       /
                                             |_______|
                                            (  O   O  )
                                             {'-(_)-'}
                                           .-{   ^   }-.
                                          /   '.___.'   \
                                         /  |    o    |  \
                                         |__|    o    |__|
                                         (((\_________/)))
                                             \___|___/
                                        jgs.--' | | '--.
                                           \__._| |_.__/
```

Warden in Go, because why not.

* [![Build Status](https://travis-ci.org/pivotal-cf-experimental/garden.png?branch=master)](https://travis-ci.org/pivotal-cf-experimental/garden)
* [![Coverage Status](https://coveralls.io/repos/pivotal-cf-experimental/garden/badge.png?branch=HEAD)](https://coveralls.io/r/pivotal-cf-experimental/garden?branch=HEAD)
* [Tracker](https://www.pivotaltracker.com/s/projects/962374)
* [Warden](https://github.com/cloudfoundry/warden)

# Running

For development, you can just spin up the Vagrant VM and run the server
locally, pointing at its host:

```bash
# if you need it:
vagrant plugin install vagrant-omnibus

# then:
librarian-chef install
vagrant up
ssh-copy-id vagrant@192.168.50.5
ssh vagrant@192.168.50.5 sudo cp -r .ssh/ /root/.ssh/
./bin/add-route
./bin/run-garden-remote-linux

# or run from inside the vm:
vagrant ssh
sudo su -
goto garden
./bin/run-garden-linux
```

This runs the server locally and configures the Linux backend to do everything
over SSH to the Vagrant box.

# Testing

## Pre-requisites

* Go 1.2 or later
* git
* mercurial
* godep

```
mkdir ~/go
```

Assuming Go was installed using gvm:
```
export GOPATH=/home/<user>/go:$GOPATH
export PATH=$PATH:/home/<user>/go/bin
```

Download a root filesystem, extract it as root, and point to it:
```
curl -O http://cfstacks.s3.amazonaws.com/lucid64.dev.tgz
sudo mkdir -p /var/warden/rootfs
sudo tar xzf lucid64.dev.tgz -C /var/warden/rootfs
export GARDEN_TEST_ROOTFS=/var/warden/rootfs
```

Get garden and its dependencies:
```
go get github.com/pivotal-cf-experimental/garden
cd ~/go/github.com/pivotal-cf-experimental/garden
godep restore
```

Make the C code:
```
make
```

Install ginkgo:
```
cd ~/go/github.com/onsi/ginkgo/ginkgo
go install
```

Run the tests:
```
cd ~/go/github.com/pivotal-cf-experimental/garden
ginkgo -r
```
