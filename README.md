Devastator
==========

[![Build Status](https://travis-ci.org/nbusy/devastator.svg?branch=master)](https://travis-ci.org/nbusy/devastator)

Devastator is a messaging server for delivering chat messages to mobile and Web clients. For each delivery target, the server uses different protocol. i.e. GCM for Android apps, WebSockets for browsers, etc. The server is completely written in Go and makes huge use of goroutines and channels.

Tech Stack
----------

GCM CCS (for message delivery and retrieval from Android clients), GAE Sockets API (CCS XMPP delivery protocol, used in place of plain TCP on AppEngine)

Architecture
------------

Messaging server utilizes device specific delivery options; GCM for Android, APNS+TCP for iOS, WebSockets for Web browsers.

```
+-------+------------+---------------+
|  GCM  |  APNS+TCP  |  Web Sockets  |
+-------+------------+---------------+
|          Messaging Server          |
+------------------------------------+
```

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

[Apache License 2.0](LICENSE)
