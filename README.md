Go-ShortURL
========

Shorturl like twitter as Daemon , tinier, faster redirect 
This e.g. I run with "proxy_pass" which in nignx 


* Require redis
* Require "github.com/hoisie/redis"

###### e.g. [http://w.adango.cn/?s=http://www.1010g.com](http://w.adango.cn/?s=http://www.1010g.com)
###### e.g. [http://w.adango.cn/fSzKfq](http://w.adango.cn/fSzKfq)

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

##### how to proxy in nginx

```nginx
	server {
       listen       80;
       server_name  w.adango.cn;

       location / {
           proxy_pass http://localhost:8887;
       }
    } 
```
