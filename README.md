# mongoOplogStats
Tool to pull metrics out of a mongod oplog


A common issue in mongo (especially with mmap) is that there is a spike of oplog
activity which causes secondary read latency. After the fact all you have from
the oplog metrics is whether it was an insert/update/delete -- which isn't
enough to look into it. This CLI tool connects to the oplog and generates metrics
per namespace to aid in investigation


## How to use
```
./mongoOplogStats --start="2018-03-15 12:22:00" --end="2018-03-15 12:27:00" | sort -k 3 -rn
```
