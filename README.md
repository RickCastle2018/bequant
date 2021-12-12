# bequant test project

## How to use:

Simply type `docker-compose up -d --build db warden distributor` in terminal while in project's directory.
MySQL, Warden (updates data) and Distributor (provides API) will be started.
Access API at `localhost:8080`.

Before running, create .env file based on example.env. This is the way to configure the service.

## API Description:

One handle: `/price`

Two params, both required: `?fsyms=CURRENCY,...&tsyms=CURRENCY,...`, where *CURRENCY* is a currency symbol (USDT, for example)

Example:

```js
Curl Request:
curl 127.0.0.1:8080/price?fsyms=BTC&tsyms=USD

JSON Response:
{
  "RAW": {
    "BTC": {
      "USD": {
        "CHANGE24HOUR": -13.25,
        "CHANGEPCT24HOUR": -0.18152873223073468,
        "OPEN24HOUR": 7299.12,
        "VOLUME24HOUR": 47600.120073200706,
        "VOLUME24HOURTO": 348033250.4911315,
        "LOW24HOUR": 7197.22,
        "HIGH24HOUR": 7426.64,
        "PRICE": 7285.87,
        "LASTUPDATE": 1586433196,
        "SUPPLY": 18313937,
        "MKTCAP": 133432964170.19
      }
    }
  },
  "DISPLAY": {
    "BTC": {
      "USD": {
        "CHANGE24HOUR": "$ -13.25",
        "CHANGEPCT24HOUR": "-0.18",
        "OPEN24HOUR": "$ 7,299.12",
        "VOLUME24HOUR": "Ƀ 47,600.1",
        "VOLUME24HOURTO": "$ 348,033,250.5",
        "HIGH24HOUR": "$ 7,426.64",
        "PRICE": "$ 7,285.87",
        "FROMSYMBOL": "Ƀ",
        "TOSYMBOL": "$",
        "LASTUPDATE": "Just now",
        "SUPPLY": "Ƀ 18,313,937.0",
        "MKTCAP": "$ 133.43 B"
      }
    }
  }
}
```

## Stack:

- Golang
- MySQL
- Docker

## Task:

Make an API cacher, I would say. We have Cryptocompare API and the task is to make
microservice which stores API answers and provides it's own API to get this data.

## Load testing

```
> ab -n 100000 -c 10 'http://localhost:8080/price?fsyms=BTC&tsyms=USD'

This is ApacheBench, Version 2.3 <$Revision: 1879490 $>
Copyright 1996 Adam Twiss, Zeus Technology Ltd, http://www.zeustech.net/
Licensed to The Apache Software Foundation, http://www.apache.org/

Benchmarking localhost (be patient)
Completed 10000 requests
Completed 20000 requests
Completed 30000 requests
Completed 40000 requests
Completed 50000 requests
Completed 60000 requests
Completed 70000 requests
Completed 80000 requests
Completed 90000 requests
Completed 100000 requests
Finished 100000 requests


Server Software:        
Server Hostname:        localhost
Server Port:            8080

Document Path:          /price?fsyms=BTC&tsyms=USD
Document Length:        584 bytes

Concurrency Level:      10
Time taken for tests:   9.498 seconds
Complete requests:      100000
Failed requests:        0
Total transferred:      69300000 bytes
HTML transferred:       58400000 bytes
Requests per second:    10528.45 [#/sec] (mean)
Time per request:       0.950 [ms] (mean)
Time per request:       0.095 [ms] (mean, across all concurrent requests)
Transfer rate:          7125.21 [Kbytes/sec] received

Connection Times (ms)
              min  mean[+/-sd] median   max
Connect:        0    0   0.1      0       1
Processing:     0    1   0.3      1       9
Waiting:        0    1   0.3      1       8
Total:          0    1   0.3      1       9

Percentage of the requests served within a certain time (ms)
  50%      1
  66%      1
  75%      1
  80%      1
  90%      1
  95%      1
  98%      2
  99%      2
 100%      9 (longest request)
```

## Requirements:

- [X] REST API
- [X] If Cryptocompare is not accessible service must return data from database via own API
(It's the purpose, or service should act just as gateway if API is accessible?)
- [X] Data in response must be fresh (realtime). 2-3 minutes discrepancy is ok.
(Again, realtime means being a gateway? I'll just do continuous caching)
- [X] Currency pairs should be configurable.
- [X] MySQL parameters should be configurable.
- [X] Service must store data to MySQL by sheduler (rawjson is ok).
- [X] Service must work in background.

## Not required, but appreciated:

- [ ] WebSockets (Why? If the service can return 2-3 min old data. To discard HTTP overhead? It won't be a problem, I think)
- [X] Clean Architecture (Uncle Bob's?)
- [X] ? Scalability (Will do!) *Add Caddy server proxy with Round-Robin algorithm (https://nknv.ru/caddy-load-balancing)*

---

Cryptocompare API: `https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BTC&tsyms=USD,EUR`
