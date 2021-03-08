# TPLink EasySmart Switch Exporter

Exports port stats from TPLink's EasySmart Switches. Very alpha, basic functionality so far:
- Port Stats

Tested on:

TPLink EasySmart Gigabit 8 Port Switch (TL-SG108E):
- v3 Firmware: 1.0.0 Build 20171214 Rel.70905 (based on web GUI)
- v4 Firmware: 1.0.0 Build 20181120 Rel.40749

Should work on the switches of the same family, but untested personally:
- TPLink EasySmart Gigabit 16 Port Switch TL-SG116E

## Grafana Dashboard

Basic Grafana dashboard using the exporter - https://grafana.com/grafana/dashboards/12517

## Usage

go run main.go --host <IP/host of switch> --username <WEBGUI username> --password <WEBGUI password>

Default username and password for this switch is admin and admin, so:

go run main.go --host 10.0.0.3 --username admin --password admin

## Run with docker:

docker run -it -p 9717:9717 thelastguardian/tplinkexporter --host 10.0.0.3 --username admin --password admin

## Metrics Exported on :9717/metrics

```
tplinkexporter_portstats_state{portnum="1"-"8",host="host"}
tplinkexporter_portstats_linkstatus{portnum="1"-"8",host="host"}
tplinkexporter_portstats_rxgoodpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_rxbadpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_txgoodpkt{portnum="1"-"8",host="host"}
tplinkexporter_portstats_txbadpkt{portnum="1"-"8",host="host"}
```