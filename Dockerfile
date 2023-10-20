FROM golang:1.21.3-alpine

WORKDIR /workspace
ENV GO111MODULE=on

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY app app
RUN go install github.com/deepmap/oapi-codegen/cmd/oapi-codegen@latest 
COPY generate-subscription-api.sh generate-subscription-api.sh
RUN sh generate-subscription-api.sh
RUN go build -o orgnote app/main.go

ENTRYPOINT ["./orgnote"]
