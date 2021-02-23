# A Prometheus exporter for FreeSWITCH

This service exports the following information from a running FreeSWITCH instance
for consumption by Prometheus:

* Calls

## Environment Variables

* `FS_HOST`: host and port of FreeSWITCH socket (default: `127.0.0.1:8021`)
* `FS_PASSWORD`: password for FreeSWITCH socket (default: `ClueCon`)
* `LISTEN`: Address to listen at(default: `:2112`)

## Authors

* Original Author: TomP <tomp@tomp.uk>
* This fork: Daniel Molkentin <danimo@infra.run>
