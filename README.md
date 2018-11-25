[Задание](https://github.com/oleggator/query-exporter/blob/master/intern_task.txt)

# PostgreSQL Query Exporter

## Installation
```
export GO111MODULE=on
go install
```

## Usage
```
query-exporter -t threads_count -c config_path
Usage of query-exporter:
  -c string
    	config file (default "./config.yml")
  -t int
    	threads count (default 1)
```

## Usage example

### Start and init DB
```
docker run -p5432:5432 -e POSTGRES_PASSWORD=postgres --name postgres --rm -d postgres:11.1-alpine
docker exec -i postgres psql -U postgres < db_up.sql
```

### Install random data generator
```
export GO111MODULE=on
go install github.com/oleggator/query-exporter/...
```

### Fill DB
```
generator -c cmd/generator/config.yml
```

### Export
```
query-exporter -t 4 -c config.yml
```
