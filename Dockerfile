FROM golang AS builder

WORKDIR /build

COPY . .

RUN go mod download
RUN go build -o todo todo.go

FROM ubuntu

WORKDIR /todo-app

COPY --from=builder /build/todo .

EXPOSE 8080

CMD [ "./todo" ]
