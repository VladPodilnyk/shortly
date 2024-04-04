### shortly

A URL shortener service at home xD
🚧 This repo is under rework 🚧

#

## How to run an app (OLD)

There are two ways to start an application. For the fist one you will need
both `make` and docker installed on your machine. For the second one you will
need to have golang installed on your local machine.

### Using make and docker

To run an app simply execute the follwing command in a project root.

```bash
> make start
```

This should start an app in a docker container. An app will run on 4000 port.

### If you have Go installed already

In case you have Golang on your machine you can start an app
with either `go ./cmd/api` or `make run-dev`.

In order to test the application use `make-test`
User `make help` to list all commands.

### Interacting with an app

There are 3 endpoints: `/v1/status`, `/v1/encode`, `/v1/decode`

A status endpoint is just a simple healthcheck:

```bash
curl -i -X GET localhost:4000/v1/status
```

An example response:

```bash
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 19 Sep 2022 19:06:45 GMT
Content-Length: 79

{
	"status": "available",
	"environment": "development",
	"version": "1.0.0"
}
```

For an encode endpoint it is required to pass an original url and an optional alias.
Here are expamples of a payload and test request:

```
{"url": "https://www.google.com/"}
{"url": "https://www.google.com/", "alias": "abc"}
curl -i -X POST localhost:4000/v1/encode -H 'Content-Type: application/json' -d '{"url": "https://www.google.com/", "alias": "abc"}
```

Response example:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 19 Sep 2022 19:11:41 GMT
Content-Length: 42

{
	"short_url": "https://short.est/abc"
}
```

Similary to encode endpoint a decode enpoint expects from user the following request body:

```
{"short_url": "https://short.est/0"}
curl -i -X POST localhost:4000/v1/decode -H 'Content-Type: application/json' -d '{"short_url": "https://short.est/0"}'
```

Response example:

```
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 19 Sep 2022 19:18:10 GMT
Content-Length: 47

{
	"original_url": "https://www.google.com/"
}
```
