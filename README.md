# 同志社休講情報

[![Go Report Card](https://goreportcard.com/badge/github.com/g-hyoga/kyuko)](https://goreportcard.com/report/github.com/g-hyoga/kyuko)

[今出川twitter](https://twitter.com/kyuko_imadegawa)
[田辺twitter](https://twitter.com/kyuko_tanabe)


Ruby -> Golang -> Golang on AWS Lambda

[ここ](http://duet.doshisha.ac.jp/kyuko/i/)
から休講情報をスクレイピングしてきて
twitterになげています

# Develope


## Build 

If you can use Docker

```sh
make build
```


or, If you can use Golang 

```sh
make local-build
```

output/handler.zip will be created

## Test

testing using docker(not working...)

```sh
make test
```

or 

testing using local Golang

```sh
make local-test
```

# Deploy

not working...
