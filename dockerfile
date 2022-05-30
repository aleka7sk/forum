FROM golang

LABEL project="web"

WORKDIR ./cmd/app 

COPY . .

RUN go build cmd/app/main.go

CMD ["./main"]