# Debugger functions for OpenFaaS

This repo contains several functions that can be helpful to testing or debugging an OpenFaaS installation.

## Functions

1. `echo` - a very simple echo function that returns in the input as output

    Usage:

    ```sh
    $ faas-cli deploy \
        --image=ghcr.io/lucasroesler/echo:latest \
        --name=echo
    $ faas-cli invoke echo <<< "payload data here"
    payload data here
    ```

2. `status-echo` - a slightly more complicated echo function that returns a JSON payload that contains the original request data: path, headers, payload

    You can request a specific response status code by adding `/<desired code>` to the function path.

    Usage:

    ```sh
    $ faas-cli deploy \
        --image=ghcr.io/lucasroesler/status-echo:latest \
        --name=status-echo
    $ faas-cli invoke status-echo <<< "payload data here" | jq
    {
        "status": 200,
        "method": "POST",
        "path": "/",
        "message": "payload data here",
        "headers": {
            "Accept": [
                "*/*"
            ],
            "Accept-Encoding": [
                "gzip"
            ],
            "Content-Type": [
                "application/x-www-form-urlencoded"
            ],
            "User-Agent": [
                "curl/7.68.0"
            ]
        }
    }
    $ curl -vs localhost:8080/function/status-echo/code/400 | jq
    *   Trying 127.0.0.1:8080...
    * TCP_NODELAY set
    * Connected to localhost (127.0.0.1) port 8080 (#0)
    > GET /function/status-echo/code/400 HTTP/1.1
    > Host: localhost:8080
    > User-Agent: curl/7.68.0
    > Accept: */*
    >
    * Mark bundle as not supporting multiuse
    < HTTP/1.1 400 Bad Request
    < Content-Length: 310
    < Content-Type: application/json
    < Date: Sat, 20 Feb 2021 14:28:37 GMT
    < X-Call-Id: a2b4f5a3-f5ac-40b5-8c59-de902cb0bb08
    < X-Duration-Seconds: 0.000311
    < X-Start-Time: 1613831317854614082
    <
    { [310 bytes data]
    * Connection #0 to host localhost left intact
    {
    "status": 400,
    "method": "GET",
    "path": "/code",
    "message": "",
    "headers": {
        "Accept": [
        "*/*"
        ],
        "Accept-Encoding": [
        "gzip"
        ],
        "User-Agent": [
        "curl/7.68.0"
        ],
        "X-Call-Id": [
        "a2b4f5a3-f5ac-40b5-8c59-de902cb0bb08"
        ],
        "X-Forwarded-For": [
        "127.0.0.1:47968"
        ],
        "X-Forwarded-Host": [
        "localhost:8080"
        ],
        "X-Start-Time": [
        "1613831317854614082"
        ]
    }
    }
    ```

3. `nats-publisher` - is a test function to use with the [`nats-connector`](https://github.com/openfaas/nats-connector). It can be configured to send messages to a specified NATS topic. The default topic is `nats-test`

    ```sh
    $ faas-cli deploy \
        --name publish-message \
        --image ghcr.io/lucasroesler/nats-publisher:latest \
        --fprocess='./handler' \
        --env nats_url=nats://nats.openfaas:4222
    $ faas-cli invoke publish-message <<< "test message"
    $ faas-cli invoke publish-message -H "nats-subject=alerts" <<< "test message"
    ```

4. `nats-reciever` - is a test function to use with the [`nats-connector`](https://github.com/openfaas/nats-connector). It can be configured at deploy time via the annotations to receive messages from a NATS topic. The received messages are then echoed to the logs for debugging

    ```sh
    $ faas-cli deploy \
        --name receive-message \
        --image ghcr.io/lucasroesler/nats-receiver:latest \
        --fprocess='./handler' \
        --annotation topic="nats-test"
    $ faas-cli logs receive-message --tail
    ```

    **Note:**
    a) this requires `nats-connector` to be installed and configured correct, see the README [instructions](https://github.com/openfaas/nats-connector#deploy-on-kubernetes)
