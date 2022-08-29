FROM golang AS builder

WORKDIR /app/backend

COPY . .

RUN go mod download
RUN go build -o todo todo.go

FROM ubuntu

WORKDIR /todo-app

COPY --from=builder /app/backend/todo .

EXPOSE 8010

CMD [ "./todo" ]
