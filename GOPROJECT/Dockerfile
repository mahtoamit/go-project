FROM golang:1.20




WORKDIR /app




# Download Go modules

#COPY go.mod go.sum /usr/src/

#RUN go mod download

COPY . .





# Copy the source code. Note the slash at the end, as explained in

# https://docs.docker.com/engine/reference/builder/#copy

#COPY *.go /var/src/

RUN go build -o main .




# Build

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /docker-gs-ping




# Optional:

# To bind to a TCP port, runtime parameters must be supplied to the docker command.

# But we can document in the Dockerfile what ports

# the application is going to listen on by default.

# https://docs.docker.com/engine/reference/builder/#expose

EXPOSE 8080




# Run

#CMD ["go", "run", "main.go"]

CMD ["./main"]
