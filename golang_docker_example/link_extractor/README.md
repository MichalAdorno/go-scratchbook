# Dockerization of a simple Golang application

Compilation:
```
    go clean
    go build
    ./link_extractor [domain-host] [domain-port]
```
Starting an HTTP server (this app):
```
    ./link_extractor localhost 8080
```
Building and running a dockerized app:
```
  docker build -t link_extractor_container --file=Dockerfile.golang .
  docker run -d -p 8080:8030 link_extractor_container
```
Calling a webservice:
```
  curl http://[domain-host]:[domain-port]/extract?from=[an-address-from-which-the-links-have-to-be-extracted]
```
For example:
```
  curl http://localhost:8080/extract?from=zeit.de
```
