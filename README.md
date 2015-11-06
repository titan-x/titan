# Devastator

[![Build Status](https://travis-ci.org/nbusy/devastator.svg?branch=master)](https://travis-ci.org/nbusy/devastator)
[![GoDoc](https://godoc.org/github.com/nbusy/devastator?status.svg)](https://godoc.org/github.com/nbusy/devastator)

Devastator is a messaging server for delivering chat messages to mobile and Web clients. For each delivery target, the server uses different protocol. i.e. GCM for Android apps, WebSockets for browsers, etc. The server is completely written in Go and makes huge use of goroutines and channels. Client server communication is full-duplex bidirectional.

## Example

See [example_test.go](example_test.go) file.

## Architecture

Messaging server utilizes device specific delivery infrastructure for notifications and messages; GCM + TLS for Android, APNS + TLS for iOS, and WebSockets for the Web browsers.

```
+-----------+------------+-------------+
| GCM + TLS | APNS + TLS | Web Sockets |
+-----------+------------+-------------+
|           Messaging Server           |
+--------------------------------------+
```

Following is the overview of the server application's components:

```
+---------------------------------------+
|              main + config            |
+---------------------------------------+
|                 server                |
+---------------------------------------+
|   router  |   listener  |  msg queue  |
+---------------------------------------+
```

## Client-Server Protocol

(Devastator server is entirely built on top of [Neptulon](https://github.com/nbusy/neptulon) framework. You can browse Neptulon repository to get more in-depth info.)

Client server communication protocol is based on [JSON RPC](http://www.jsonrpc.org/specification) 2.0 specs. Mobile devices connect with the TLS endpoint and the Web browsers utilizes the WebSocket endpoint. The message framing on the TLS endpoint is quite simple:

```
[uint32] 4 Bytes Content Length Header + [JSON] Message Body
```

Following is a valid TLS message frame for connecting mobile devices (except the header, which should be a 4 byte binary encoded uint32):

```
xxxx{method: "ping"}
```

## Client Authentication

First-time registration is done through Google+ OAuth 2.0 flow. After a successful registration, the connecting device receives a client-side TLS certificate (for mobile devices) or a JSON Web Token (for browsers), to be used for successive connections.

## Typical Client-Server Communication

Client-server communication sequence is pretty similar to that of XMPP, except we are using JSON RPC packaging for messages.

```
[Server]                    [Client]
+                                  +
|---------GCM Notification-------->| [offline]
|                                  |
|                                  |
|<-----------auth.cert-------------| [online]
|                                  |
|----------ACK/batch[msg]--------->|
|                                  |
|<-----------batch[ACK]------------|
|                                  |
|                                  |
|<------------msg.echo-------------|
|                                  |
|---------------ACK--------------->|
|                                  |
|-------------msg.echo------------>|
|                                  |
|<--------------ACK----------------|
+                                  +
```

Any message that was not acknowledged by the client will be delivered again (hence at-least-once delivery princinple), unless TTL of the message was reached. Client implementations should be ready to handle occasional duplicate deliveries of messages by the server. Message IDs will remain the same for duplicates.

## Testing

All the tests can be executed with `GORACE="halt_on_error=1" go test -race -cover ./...` command. Optionally you can add `-v` flag to observe all connection logs. Integration tests require environment variables defined in the next section. If they are missing, integration tests are skipped.

## Environment Variables

Following environment variables needs to be present on any dev or production environment:

```bash
export GOOGLE_API_KEY=
export GOOGLE_PREPROD_API_KEY=
```

## Logging and Metrics

Only actionable events are logged (i.e. server started, client connected on IP ..., client disconnected, etc.). You can use logs as event sources. Anything else is considered telemetry and exposed with `expvar`. Queue lengths, active connection/request counts, performance metrics, etc. Metrics are exposed via HTTP at /debug/vars in JSON format.

## Performance Notes

The messaging server is designed to make max usage of available CPU resources. However exceeding 100% CPU usage will cause a memory usage spike as marshalled/unmarshalled messages and other allocated byte buffers will have to reside in memory much longer. Ideally, 95% CPU usage should trigger the clustering mechanism which should spawn more server instance. Currently there is no clustering support built-in, but it is a top priority.

## License

[MIT](LICENSE)
