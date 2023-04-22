FROM golang:1.17 as DEPS

ENV APP_USER app

ENV STAGE "local"
ENV DOCDB_USER "myUserAdmin"
ENV DOCDB_PASS "abc123"
ENV DOCDB_ENDPOINT "172.18.0.3"
ENV DOCDB_DB "quillpen"
ENV DOCDB_ACCOUNTS "accounts"
ENV APP_HOME /go/src/quillpen


WORKDIR $APP_HOME
COPY go.mod .
COPY go.sum .


RUN go mod download

COPY . .
RUN go build
ENTRYPOINT ["./quillpen"]