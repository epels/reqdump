# reqdump

reqdump is an extremely basic server dumping any incoming HTTP requests.

## Usage

```bash
$ reqdump -help
Usage of reqdump:
  -addr string
    	addr optionally specifies the TCP address for the server to listen on, in the form "host:port" (default ":8080")
  -body
    	If true, also dump the request body (default true)
```

## Example

```bash
$ reqdump
[INFO]: 2020/12/27 16:32:40 Accepting connections on ":8080"

GET /path HTTP/1.1
Host: localhost:8080
Accept: */*
User-Agent: curl/7.64.1


PUT /some-other-path HTTP/1.1
Host: localhost:8080
Accept: */*
Content-Length: 14
Content-Type: application/json
User-Agent: curl/7.64.1

{"foo": "bar"}
``` 
