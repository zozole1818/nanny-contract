# nanny-contract

App for calculating ZUS contributions for nanny. RCA and DRA calculations.

There are two variants of the application:
- [server](https://github.com/zozole1818/nanny-contract/tree/master/cmd/server) - HTTP server
- [cli](https://github.com/zozole1818/nanny-contract/tree/master/cmd/cli) - command line interface

If you run the HTTP server, you can access the UI at [http://localhost:8080/views/reports-new](http://localhost:8080/views/reports-new).

If you run the command line interface, you can use the `rca` and `dra` commands.
```bash
$ zus -h
CLI to quickly generate ZUS RCA/DRA data

Usage:
  zus [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  dra         Help to calculate ZUS DRA
  help        Help about any command
  rca         Help to calculate ZUS RCA

Flags:
  -h, --help      help for zus
  -v, --verbose   verbose output

Use "zus [command] --help" for more information about a command.

```

