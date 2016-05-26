msplapi
-----------

Data access layer for marsupi


Setup
=====

Install and run *Postgres*: [Postgres App](http://postgresapp.com/) on OSX is usually the easiest way

Install and run *RabbitMQ*: `brew update` + `brew install rabbitmq` + `rabbitmq-server`
[https://www.rabbitmq.com/install-homebrew.html](https://www.rabbitmq.com/install-homebrew.html) for more help.

Build msplapi and setup the datbase: `make setup`

If you're running on localhost, you need to install and run [ngrok](https://ngrok.com) for webhooks to work. Upgrading to ngrok Pro plan is required.

`make run-ngrok` to run it


Run
===

`msplapi run`


Testing
=======

The api_test package needs postgres to be running. Connection limit needs to be increased (set `max_connections = 1000` in `postgresql.conf`)

`make test`


Commands
========

`msplapi` shows available commands


Environment Config
==================

see `config/config.go`

