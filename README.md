NBusy Server
============

NBusy messaging server for delivering all chat messages to all devices (mobile apps + the browser). For each delivery target, the server uses different protocol. i.e. GCM for the NBusy Android app, WebSockets for nbusy.com, etc. The server is completely written in Go and makes huge use of goroutines and channels.
