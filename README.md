# http-server
Simple static file server with idle tracking

## Basic command line usage

```
Usage of http-server:
  -d string
        the directory of files to host (default ".")
  -p string
        port to serve on (default "8100")
  -t uint
        duration before shutdown while inactive (0 – disable)
```

## Configuration
Application supports multiple served locations that should be configured in the `config.yml` (`~/.config/http-server/config.yml`)

```
locations:
    - path: "/example/"
      root: "/var/www/example.com"
    - path: "/doc/"
      root: "/var/www/doc"
```

In this case root directory specified via command line arguments ignored