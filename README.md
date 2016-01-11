# Simple web server in [Go](http://golang.org) with [PostgresSQL 9.5](http://www.postgresql.org/)

## Build
```
go get github.com/lib/pq
go install github.com/michalkowol/web-pg/server
./server
```

## Dependencies:
* [pg](https://github.com/lib/pq) - A pure Go postgres driver for Go's database/sql package