# About

Gets information about SSL certificates, so you know when it's time to renew them.

# Usage

Print a table to the terminal with colors and formatted header. Header is bold, and days left is either red or green depending on time left.

```bash
$ go-cert --colors --formatting example.com github.com

   Domain     Days left       End date       Status
example.com         135   2018-11-28 13:00     ok  
github.com          688   2020-06-03 14:00     ok
```

Print JSON

```bash
$ go-cert --output json example.com github.com
```
```json
{"domains":[{"name":"example.com","daysLeft":135,"endTime":"2018-11-28T13:00:00+01:00","status":"ok"},{"name":"github.com","daysLeft":688,"endTime":"2020-06-03T14:00:00+02:00","status":"ok"}]}
```

Print text

```bash
$ go-cert --output text example.com github.com

example.com|135|2018-11-28 13:00:00 +0100 CET|ok
github.com|688|2020-06-03 14:00:00 +0200 CEST|ok
```


# Help
```
NAME:
   go-cert - check days left on SSL certificates

USAGE:
   go-cert [OPTION]... DOMAIN [DOMAIN ...]

VERSION:
   0.0.1

AUTHOR:
   Rasmus Lindroth <rasmus@lindroth.xyz>

COMMANDS:
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --days INT, -d INT      days INT left on certificate warning (default: 20)
   --location LOC, -l LOC  LOC used for time zone, e.g. Europe/Stockholm. Defaults to local
   --output TYPE, -o TYPE  output TYPE: table, json, text (| seperator) (default: "table")
   --colors, -c            add colors in table output
   --formatting, -f        add bold in table header
   --help, -h              show help
   --version, -v           print the version
```
