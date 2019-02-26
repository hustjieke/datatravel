datatravel is designed for data migration between data sources such as relational databases, NoSQL, and big data (OLAP).

Shift

$make build

Usage:
```
./bin/shift --from=[host:port] --from-database=[database] --from-table=[table] --from-user=[user] --from-password=[password] --to=[host:port] --to-database=[database] --to-table=[table]  --to-user=[user] --to-password=[password]
```

For example:
```
./bin/shift --from=192.168.0.2:3306 --from-database=sbtest --from-table=benchyou0_0031 --from-user=mock --from-password=mock --to=192.168.0.9:3306 --to-database=sbtest --to-table=benchyou0_0031 --to-user=mock --to-password=mock
```
