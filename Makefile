REGISTRY_OWNER:=fazenda
MULTIARCH:=false
ARCHS:=linux/amd64
PROJECT_TAG:=latest

ifeq (true, $(MULTIARCH))
	# for some unkown reason linux/s390x doesn't build
	ARCHS:=linux/386,linux/amd64,linux/arm/v6,linux/arm/v7,linux/arm64/v8,linux/ppc64le
endif

all: install setup

install:
	@curl -fSL https://get.docker.com | sh
	sudo usermod -aG docker $USER
	sudo systemctl enable docker
	sudo systemctl start docker

setup:
	@docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
	@docker buildx rm builder
	@docker buildx create --name builder --driver docker-container --use
	@docker buildx inspect --bootstrap

run: build
	@docker run --rm -it --env USER=${USER} typer

build:
	@docker buildx build --platform linux/amd64 --load --tag typer .

deploy:
	@docker buildx build --platform $(ARCHS) --push --tag ${REGISTRY_OWNER}/typer:${PROJECT_TAG} .
