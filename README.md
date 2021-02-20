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
    ```
