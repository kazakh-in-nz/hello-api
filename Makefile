GO_VERSION :=1.22

.PHONY: install-go init-go
 
setup: install-go init-go
 
#TODO add MacOS support
install-go:
  wget "https://golang.org/dl/go$(GO_VERSION).linux-amd64.tar.gz"
  sudo tar -C /usr/local -xzf go$(GO_VERSION).linux-amd64.tar.gz
  rm go$(GO_VERSION).linux-amd64.tar.gz
 
init-go:
  echo 'export PATH=$$PATH:/usr/local/go/bin' >> $${HOME}/.zshrc
  echo 'export PATH=$$PATH:$${HOME}/go/bin' >> $${HOME}/.zshrc

build:
  go build -o api cmd/main.go
