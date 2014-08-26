NBusy Server
============

[![Build Status](https://travis-ci.org/nbusy/nbusy-server.svg?branch=master)](https://travis-ci.org/nbusy/nbusy-server)

NBusy messaging server for delivering all chat messages to all devices (mobile apps + the browser). For each delivery target, the server uses different protocol. i.e. GCM for the NBusy Android app, WebSockets for nbusy.com, etc. The server is completely written in Go and makes huge use of goroutines and channels.

## Tech Stack
GCM CCS (for message delivery and retrieval from Android clients), GAE Sockets API (CCS XMPP delivery protocol, used in place of plain TCP on AppEngine)

## Environment Variables
Following environment variables needs to be present on any dev or production environment:

```bash
export GCM_SENDER_ID=
export GOOGLE_API_KEY=
```