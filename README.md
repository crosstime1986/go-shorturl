Go-ShortURL
========

shorturl like twitter as  Daemon 

* request redis
* request "github.com/hoisie/redis"

###### e.g. [http://w.1010g.cn/?s=http://www.1010g.com](http://w.1010g.cn/?s=http://www.1010g.com)
###### e.g. [http://w.1010g.cn/fSzKfq](http://w.1010g.cn/fSzKfq)

##### how to user

```shell
go run shorturl.go [:port]
go run shorturl.go 8887 
```

as you'd like build

```shell
go build -o shorturl shorturl.go 
./shorturl 8887
```