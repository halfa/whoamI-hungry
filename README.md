# whoamI

Tiny Go webserver that prints os information and HTTP request to output

```sh
$ docker run -d -P --name iamfoo emilevauge/whoami
$ docker inspect --format '{{ .NetworkSettings.Ports }}'  iamfoo
map[80/tcp:[{0.0.0.0 32769}]]
$ curl "http://0.0.0.0:32769"
Hostname :  6e0030e67d6a
IP :  127.0.0.1
IP :  ::1
IP :  172.17.0.27
IP :  fe80::42:acff:fe11:1b
GET / HTTP/1.1
Host: 0.0.0.0:32769
User-Agent: curl/7.35.0
Accept: */*
```

Burn
----

Use `/burn` to produce a heavy CPU load. Current iteration number (`100000`) is the reference for a `i7-4800MQ` CPU.

### Curl-it
```
while true    
do
  curl -H Host:whoami2.docker.localhost http://localhost/burn
  sleep 0.05
done
```
