Devastator
==========

[![Build Status](https://travis-ci.org/nbusy/devastator.svg?branch=master)](https://travis-ci.org/nbusy/devastator) [![GoDoc](https://godoc.org/github.com/nbusy/devastator?status.svg)](https://godoc.org/github.com/nbusy/devastator)

Devastator is a messaging server for delivering chat messages to mobile and Web clients. For each delivery target, the server uses different protocol. i.e. GCM for Android apps, WebSockets for browsers, etc. The server is completely written in Go and makes huge use of goroutines and channels.

Tech Stack
----------

GCM CCS (for message delivery and retrieval from Android clients), GAE Sockets API (CCS XMPP delivery protocol, used in place of plain TCP on AppEngine)

Architecture
------------

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

Client-Server Protocol
----------------------

Client server communication protocol is based on [JSON RPC](http://www.jsonrpc.org/specification) 2.0 specs. Mobile devices connect with the TLS endpoint and the Web browsers utilizes the WebSocket endpoint. The message framing on the TLS endpoint is quite simple:

```
[uint32] 4 Bytes Content Length Header + [JSON] Message Body
```

Following is a valid TLS message frame for connecting mobile devices (except the header, which should be a 4 byte binary encoded uint32):

```
xxxx{method: "ping"}
```

Client Authentication
---------------------

First-time registration is done through Google+ OAuth 2.0 flow. After a successful registration, the connecting device receives a client-side TLS certificate (for mobile devices) or a JSON Web Token (for browsers), to be used for successive connections.

Typical Client-Server Sequence
------------------------------

Client-server communication sequence is pretty similar to that of XMPP, except we are using JSON RPC packaging.

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

Any message that was not acknowledged will be delivered to the client again (hence at-least-once delivery princinple), unless TTL of the message was reached.

Testing
-------

All the tests can be executed by `GORACE="halt_on_error=1" go test -v -race -cover ./...` command. Optionally you can add `-v` flag to observe all connection logs. Integration tests require environment variables defined in the next section. If they are missing, integration tests are skipped.

Environment Variables
---------------------

Following environment variables needs to be present on any dev or production environment:

```bash
export GOOGLE_API_KEY=
export GOOGLE_PREPROD_API_KEY=
```

Logging and Metrics
-------------------

Only actionable events are logged. You can use logs as event sources. Anything else is considered telemetry and exposed with `expvar`. Queue lengths, active connection/request counts, performance metrics, etc. Metrics are exposed via HTTP at /debug/vars in JSON format.

Performance Notes
-----------------

The messaging server is designed to make max usage of available CPU resources. However exceeding 100% CPU usage will cause a memory usage spike as marshalled/unmarshalled messages and other allocated byte buffers will have to reside in memory much longer. Ideally, server process' CPU usage should never exceed 95% of overall system CPU resources. Currently there is no clustering support but it is a top priority.

License
-------

[MIT](LICENSE)
