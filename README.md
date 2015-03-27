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

Client-Server Protocol
----------------------

Client server communication protocol is based on [JSON RPC](http://www.jsonrpc.org/specification) 2.0 specs. Mobile devices connect with the TLS endpoint and the Web browsers utilizes the WebSocket endpoint.

Client Authentication
---------------------

First-time registration is done through Google+ OAuth 2.0 flow. After a successful registration, the connecting device receives a client-side TLS certificate (for mobile devices) or a JSON Web Token (for browsers), to be used for successive connections.

Testing
-------

All the tests can be executed by `go test -race -cover ./...` command. Optionally you can add `-v` flag to observe all connection logs. Integration tests require environment variables defined in the next section. If they are missing, integration tests are skipped.

Environment Variables
---------------------

Following environment variables needs to be present on any dev or production environment:

```bash
export GOOGLE_API_KEY=
export GOOGLE_PREPROD_API_KEY=
```

License
-------

[MIT](LICENSE)
