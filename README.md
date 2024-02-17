# Video-hosting PROJECT

## Launch with Docker
```shell
$ docker build -f build/Dockerfile --tag "video-hosting" .
$ docker run -d -p 8080:8080 "video-hosting"
```