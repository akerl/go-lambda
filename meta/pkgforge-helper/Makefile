DOCKER_CMD = docker run \
	--rm -t -i \
	-v $$(pwd):/opt/build \
	-v /tmp:/tmp \
	$(DOCKER_FLAGS) \
	$(CONTAINER_NAME)

.PHONY : default manual dircheck container prereqs release build

default: prereqs
	$(DOCKER_CMD) pkgforge build $(PKGFORGE_FLAGS)

build: prereqs
	pkgforge build -ts

release: prereqs
	$(DOCKER_CMD) pkgforge release $(PKGFORGE_FLAGS)

manual: prereqs
	$(DOCKER_CMD) bash || true

prereqs: dircheck container auth

ifdef GITHUB_CREDS
GITHUB_CRED_NAME ?= targit
auth:
	@echo "$(GITHUB_CRED_NAME): $(GITHUB_CREDS)" > .github || true
else
auth:
	@true
endif

PKGFORGE_FLAGS =
ifdef PKGFORGE_STATEFILE
PKGFORGE_FLAGS += --statefile $(PKGFORGE_STATEFILE)
endif
ifdef DEBUG
PKGFORGE_FLAGS += -ts
endif

ifneq ("$(wildcard .pkgforge)","")
dircheck:
	@true
else
dircheck:
	@echo ".pkgforge not found; run make from the repo root"
	@false
endif

ifneq ("$(wildcard Dockerfile)","")
CONTAINER_NAME = $$(awk '/^name / {print $$2}' .pkgforge | tr -d "'")-pkg
container:
	docker build -t $(CONTAINER_NAME) .
else
CONTAINER_NAME = dock0/pkgforge
container:
	@true
endif

ifneq ("$(wildcard docker.conf)","")
DOCKER_FLAGS = $$(cat docker.conf)
else
DOCKER_FLAGS =
endif
