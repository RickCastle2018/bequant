# bequant test project

## How to use:

Simply type `docker-compose up -d --build db warden distributor` in terminal while in project's directory.
MySQL, Warden (updates data) and Distributor (provides API) will be started.
Access API at `localhost:2311`.

## API Description:

One handle: `/price`

Two params, both required: `?fsyms=CURRENCY&tsyms=CURRENCY`, where *CURRENCY* is a currency symbol (USDT, for example)

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

## Requirements:

- REST API (mb, no need to use WebSockets?)
- If Cryptocompare is not accessible service must return data from database via own API
(It's the purpose, or service should act just as gateway if API is accessible?)
- Data in response must be fresh (realtime). 2-3 minutes discrepancy is ok.
(Again, realtime means being a gateway? I'll just do continuous caching)
- Currency pairs should be configurable.
- MySQL parameters should be configurable.
- Service must store data to MySQL by sheduler (rawjson is ok).
- Service must work in background.

## Not required, but appreciated:

- WebSockets (Why? If the service can return 2-3 min old data. To discard HTTP overhead? It won't be a problem, I think)
- Clean Architecture (Uncle Bob's?)
- Scalability (Will do!)

---

Cryptocompare API: `https://min-api.cryptocompare.com/data/pricemultifull?fsyms=BTC&tsyms=USD,EUR`
