FROM golang:alpine

COPY main.go /app/main.go

CMD ["go","run","/app/main.go"]
