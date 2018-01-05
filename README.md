## shrinkUrl

A simple URL shortening API project in Go.


### Requirements

- [Go](https://golang.org)
- [Redis](https://redis.io)

### Installation

`go get github.com/khaight/shrinkUrl`


### Running

`go build -o main && ./main`


### API Usage

##### Create New URL
```
curl http://localhost:3000/api/url -X POST -d '{ "url": "http://www.github.com/khaight" }'
```

##### Get URL
```
curl http://localhost:3000/api/url/6i
```

```json
{
  "shortURL":"http://localhost:3000/6i",
  "longURL":"http://www.github.com/khaight",
  "visits":0,
  "created":1515083989107066000
}
```

### Config EnvVars

#####  SHRINK_UR_APP_HOST
The app host name

```sh
export SHRINK_URL_APP_HOST=localhost:3000
```

#####  SHRINK_URL_APP_PORT  (default 3000)

```sh
export SHRINK_URL_APP_PORT=localhost
```

#####  SHRINK_URL_REDIS_HOST  (default localhost:6379)

```sh
export SHRINK_URL_REDIS_HOST=localhost:6379
```
