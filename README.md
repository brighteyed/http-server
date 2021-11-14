# http-server
Simple http server with idle tracking. It can serve local files and contents of zip archives also.

## Basic command line usage

```
$ Usage of http-server:
  -d string
        the directory of files to host
  -p string
        port to serve on (default "8100")
  -t uint
        duration before shutdown while inactive (0 – disable)
```

## Configuration
Application supports multiple served locations. Configuration is loaded from all yaml files placed in `http-server` subdirectory of `XDG_DATA_DIRS` or `XDG_CONFIG_HOME` directories. 

> **NOTE**: If root directory is specified with `-d` command line argument then server doesn't load its configuration from files.

### Example
```
locations:
    - path: /tasks/
      root: /var/www/tasks
    - path: /notes/
      root: /var/notes.zip
```

