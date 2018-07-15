# About

Gets information about SSL certificates, so you know when it's time to renew them.

# Usage

Print table with colors and formatted header

```
go-cert --colors --formatting example.com github.com
```

Print JSON

```
go-cert --output json example.com github.com
```

# Help
```
NAME:
   go-cert - check days left on SSL certificates

USAGE:
   go-cert [global options] command [command options] domains

GLOBAL OPTIONS:
   --days value, -d value      days left on certificate warning (default: 20)
   --location value, -l value  used for time zone, e.g. Europe/Stockholm. Defaults to local
   --output value, -o value    table (default), json
   --colors, -c                add colors in table output
   --formatting, -f            uses bold in table header
   --help, -h                  show help
   --version, -v               print the version
```
