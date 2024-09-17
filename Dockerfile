# syntax=docker/dockerfile:1

FROM golang:1.19

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the source code. Note the slash at the end, as explained in
# https://docs.docker.com/engine/reference/builder/#copy
#COPY *.go ./
#COPY *.go ./
COPY . ./

RUN ls

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /prewave

# To bind to a TCP port, runtime parameters must be supplied to the docker command.
EXPOSE 8080

# Run
CMD [ "/prewave" ]