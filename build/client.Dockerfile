# Собираем в гошке
FROM golang:1.21 as build

ENV BIN_FILE /opt/smdaemon/client
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
        -o ${BIN_FILE} cmd/client/*

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="OTUS Online Education"
LABEL SERVICE="smdaemon"
LABEL MAINTAINERS="grahovsky@gmail.com"

ENV BIN_FILE "/opt/smdaemon/client"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

CMD ${BIN_FILE}
