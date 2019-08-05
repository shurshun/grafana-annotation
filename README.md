grafana-annotation: Post graphite annotation to grafana 5+.

# Install

go get -d "github.com/contentsquare/grafana-annotation"

## Create a Bearer Token

[Read the Docs](http://docs.grafana.org/http_api/auth/)

# Build

```
go build
```

# Call

## Options

```
NAME:
   annotation-poster - Tool to post graphite annotations to grafana

USAGE:
   main [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --data value   Additional data. [$GRAFANA_DATA]
   --what value   The What item to post. [$GRAFANA_WHAT]
   --tags value   Tags. [$GRAFANA_TAGS]
   --token value  Bearer Token. [$GRAFANA_TOKEN]
   --uri value    Example: https://some-grafana-host.tld [$GRAFANA_URI]
   --help, -h     show help
   --version, -v  print the version
```

## Example call

```
~$ grafana-annotation --data "Details on this event" --tags foo,bar \
   -what "Something happened on system foo with bar event"
```

