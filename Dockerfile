FROM golang:1.24 AS build

ENV BIN_FILE /opt/app/anti_bruteforce
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

ARG LDFLAGS

RUN go install github.com/swaggo/swag/cmd/swag@latest && \
    swag i -d internal/server/http/,internal/app/,internal/common/ -o internal/server/http/docs/ -g server.go

RUN CGO_ENABLED=0 go build -v \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} cmd/calendar/*.go


FROM alpine:latest

ENV BIN_FILE /opt/calendar/calendar

COPY --from=build ${BIN_FILE} ${BIN_FILE}

ENV CONFIG_FILE /etc/calendar/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE} -config ${CONFIG_FILE}
