# http_request_checker

This is used to check requests' detail.

It executes a HTTP server that shows a request as a response on a browser and a log on STDOUT.

# Usage

Default port is `8080`.

```console
go get
go build -o main main.go
chmod 700 ./main

./main
```

## With specific port

```console
PORT=3939 ./main
```

## With a responser for Slack URL verification

```console
FOR_SLACK=true ./main
```

Slack requires a URL that returns `challenge` value for an APP. [url\_verification event \| Slack](https://api.slack.com/events/url_verification)

This sets up the endpoint to meet the requirement.
