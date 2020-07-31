# timescaleproxy
A proxy for writing Telegraf JSON HTTP outputs into a TimescaleDB. Written in Go. Feedback is welcome. Bear in mind this is my first Go-project.

As the Telegraf plugin for TimescaleDB is [pending](https://github.com/influxdata/telegraf/pull/3428) for being included in the Telegraf codebase. 
I saw the need for a workaround using a HTTP-proxy that recieves metrics in JSON and writes them to the TimescaleDB.

I'm looking at the [plugin](https://github.com/svenklemm/telegraf/tree/postgres/plugins/outputs/postgresql) but instead of a general PostgreSQL output it outputs only to TimescaleDB.

Security is non-existant at the moment so this should only be used far away from the internetz. 

Telegraf config
```yaml
[[outputs.http]]
  url = "http://url.to.server:aport/"
  data_format = "json"
  [outputs.http.headers]
    Content-Type = "application/json; charset=utf-8"
```

Functionality as of 2020-07-31:
- Export and read JSON-config
- Create tables and tags in DB from first Metric
- Writes JSON-metrics to TimescaleDB




# TODO
- ~~config~~
- ~~create tables from first measurement~~
- update table if measurement changes (unlikely)
- ~~write measurement to db~~
- ~~tags as FK~~
- ~~schema config~~
- handle influx line protocol (less portable maybe)
- measurement config (column type for field)
- sanitize strings
- default retention policy
- error handling
