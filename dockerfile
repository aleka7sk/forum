FROM golang

LABEL project="g-t"

WORKDIR ./cmd/app 

COPY . .

RUN go build cmd/app/main.go

CMD ["./main"]