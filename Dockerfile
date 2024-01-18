FROM golang:1.21.6-alpine

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY **/*.go ./

RUN go build -o /thesaurum

CMD [ "/thesaurum" ]
