FROM golang:1.22-alpine

WORKDIR /app

COPY . .

RUN go build -o server ./main/main.go

EXPOSE 5050

CMD [ "./server" ]