msplapi
-----------

Data access layer for marsupi


Setup
=====

Run postgres SQL locally

`make`

If you're running on localhost, you need to run [ngrok](https://ngrok.com) `ngrok http -subdomain=msplapi 8081` for webhooks to work. Upgrading to ngrok Pro plan is required.


Run API Server
==============

`msplapi run`


Testing
=======

The api_test package needs postgres to be running. Connection limit needs to be increased (set `max_connections = 1000` in `postgresql.conf`

`make test`


Database
========

`msplapi migrate up`
`msplapi migrate down`
`msplapi create-db`
`msplapi clean-db`
`msplapi seed-db`

Environment Config
==================

```
# URL for the api service
export MSPL_API_URL=http://msplapi.ngrok.io
```
