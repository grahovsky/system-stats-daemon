# Собираем в гошке
FROM golang:1.21 as build

ENV BIN_FILE_SERVICE /opt/smdaemon/service
ENV BIN_FILE_CLIENT /opt/smdaemon/client
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_SERVICE} cmd/service/* && \
    CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE_CLIENT} cmd/client/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="smdaemon"
LABEL MAINTAINERS="grahovsky@gmail.com"

ENV BIN_FILE_SERVICE "/opt/smdaemon/service"
ENV BIN_FILE_CLIENT "/opt/smdaemon/client"

COPY --from=build ${BIN_FILE_SERVICE} ${BIN_FILE_SERVICE}
COPY --from=build ${BIN_FILE_CLIENT} ${BIN_FILE_CLIENT}

ENV CONFIG_FILE /etc/smdaemon/config.yaml
COPY ./configs/config.yaml ${CONFIG_FILE}

CMD ${BIN_FILE_SERVICE} --config ${CONFIG_FILE}
