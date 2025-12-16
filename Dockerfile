FROM golang:1.24-alpine

WORKDIR /workspace
ENV GO111MODULE=on

COPY go.mod go.sum entrypoint.sh ./
RUN go mod download && go mod verify

COPY app app
COPY migrations migrations
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest 
COPY generate-subscription-api.sh generate-subscription-api.sh
RUN sh generate-subscription-api.sh
RUN go build -o orgnote app/main.go
RUN go install -tags 'mongodb' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.1

ENTRYPOINT ["./entrypoint.sh"]
