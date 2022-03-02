FROM golang:latest

# create a working directory
WORKDIR /delta

# Fetch dependencies on separate layer as they are less likely to
# change on every build and will therefore be cached for speeding
# up the next build
COPY ./go.mod ./go.sum ./
RUN go mod download

# copy source from the host to the working directory inside
# the container
COPY . .

RUN go build -o server ./main.go

# This container exposes port 1337 to the outside world
EXPOSE 8000
