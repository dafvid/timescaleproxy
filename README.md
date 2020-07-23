# timescaleproxy
A proxy for writing Telegraf JSON HTTP outputs into a TimescaleDB. Written in Go

As the Telegraf plugin for TimescaleDB is [pending](https://github.com/influxdata/telegraf/pull/3428) for being included in the Telegraf codebase. 
I saw the need to try and create a simple workaround by using a simple HTTP-proxy that recieves JSON and writes that to the TimescaleDB.

The plan is to borrow a lot of code from the [plugin](https://github.com/svenklemm/telegraf/tree/postgres/plugins/outputs/postgresql) but instead of making it a general PostgreSQL output it just writes to TimescaleDB.

As of now (2020-07-21) it's just the initial commit of something that recieves JSON and creates a Go map.
- 22/7  Created a module, read/write config
- 23/7  Parse JSON into []Metric

# TODO
- ~~config~~
- create tables from first measurement
- update table if measurement changes
- write measurement to db
- tags as FK
- schema config
