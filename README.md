# timescaleproxy
A proxy for writing Telegraf JSON HTTP outputs into a TimescaleDB. Written in Go. Feedback is welcome. Bear in mind this is my first Go-project.

As the Telegraf plugin for TimescaleDB is ~[pending](https://github.com/influxdata/telegraf/pull/3428)~ [pending](https://github.com/influxdata/telegraf/pull/8651) for being included in the Telegraf codebase. 
I saw the need for a workaround using a HTTP-proxy that recieves metrics in JSON and writes them to the TimescaleDB.

It's inspired by the [plugin](https://github.com/svenklemm/telegraf/tree/postgres/plugins/outputs/postgresql) but instead of a general PostgreSQL output it outputs only to TimescaleDB.

Security is non-existant at the moment so this should only be used far away from the internetz. I'd recommend either behind a secured nginx reverse proxy or over a [wireguard](https://www.wireguard.com) interface.

## Install
  ```sh
  git clone https://github.com/dafvid/timescaleproxy.git
  cd timescaleproxy
  go install
  go build
  ```

## Config
  ```sh
  ./timescaleproxy -printconf > config.json
  ```
  
  ### Db ###
  - **Db.Host**: PostgreSQL server host
  - **Db.Port**: PostgreSQL port
  - **Db.Schema**: Timescale database schema
  - **Db.Database**: Timescale database name
  - **Db.User**: Timescale database user
  - **Db.Password**: Timescale database user password
  - **Db.MaxConns**: MaxConns setting for pgxpool
  - **Db.MinConns**: MinConns setting for pgxpool
  ### Listen ###
  - **Listen.adress**: Timescale proxy listening adress
  - **Listen.port**: Timescale proxy listening port
  ### Misc ###
  - **TimestampUnit**: matches the setting `json_timestamp_units` in section `[[outputs.http]]` in `telegraf.conf`
  - **DefaultDropPolicy**: Interval setting for [add_retention_policy()](https://docs.timescale.com/api/latest/data-retention/add_retention_policy/). Clear to disable retention.
  - **LogLevel**: Log level. Clear to disable logs.
    
## Run
  ```sh
  ./timescaleproxy -c config.json
  ```
  

## Telegraf config
```yaml
[[outputs.http]]
  url = "http://url.to.server:8432/"
  data_format = "json"
  json_timestamp_units = "1ms"
  [outputs.http.headers]
    Content-Type = "application/json; charset=utf-8"
```

## Updates
### 2022-10-07
- Added Config documentation
### 2022-08-10
- Added FreeBSD RC-script to run the proxy as a daemon
### 2022-06-22
- Changed the timestamp column name to ts
### 2020-07-31:
- Export and read JSON-config
- Create tables and tags in DB from first Metric
- Writes JSON-metrics to TimescaleDB

# TODO
- Per measurement config for retention
- ~~config~~
- ~~create tables from first measurement~~
- ~~write measurement to db~~
- ~~tags as FK~~
- ~~schema config~~
- handle influx line protocol (less portable maybe)
- measurement config (column type for field)
- sanitize strings
- ~~default retention policy~~
- ~~error handling~~
- update table if measurement changes (unlikely)
