VERSION := $(shell git describe --tags --dirty 2>/dev/null)
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null)
TIME := $(shell date +"%Y-%m-%d.%H:%M:%S")

GOOS ?= $(shell uname -s | tr [:upper:] [:lower:])
GOARCH ?= amd64

IGNORE = ./ ./vendor/
PACKAGES := $(filter-out $(IGNORE), $(sort $(dir $(wildcard ./*/))))
LINTS := $(filter-out ,$(PACKAGES))

FLAGS = -ldflags "-X stem-cell/cmd.Version=$(VERSION) -X stem-cell/cmd.GitDegest=$(COMMIT) "

app:
	@export GOOS=$(GOOS) GOARCH=$(GOARCH)
	go build $(FLAGS)


.PHONE: app

lint: $(LINTS:=.lint)

$(LINTS:=.lint):
	golint -set_exit_status $(subst .lint,,$@)...

unittest: $(PACKAGES:=.test) $(EXAMPLES:=.build)

$(PACKAGES:=.test):
	$(eval d := $(patsubst %.test,%,$@))
	@cd $d && go test ./...
