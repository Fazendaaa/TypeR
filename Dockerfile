FROM golang:1.11-alpine AS build

ARG CGO_ENABLED=0
ARG GOOS=linux

WORKDIR /go/src/github.com/Fazendaaa/TypeR

COPY src .

RUN [ "go", "test", "./..." ]
RUN [ "go", "build", "-o", "/usr/bin/typer", "main.go" ]

FROM scratch AS bin

COPY --from=build /usr/bin/typer /

ENTRYPOINT [ "/typer" ]
