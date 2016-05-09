msplapi
-----------

Data access layer for marsupi


Setup
=====

Run postgres SQL locally

`make`


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
