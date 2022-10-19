# kutego-api

KuteGo is an API to play with cute [Aurelie Vache's Cats](https://github.com/Pandahoro/cats)

For now, the API provide:

* the Cats list
* a route to display one Cat of your choice 😍
* and a random Cat!

You can resize them with query param 
```
http://localhost:8080/cat/random?size=[x-small,small,medium]
```


![Cat McFly](https://raw.githubusercontent.com/Pandahoro/cats/main/back-to-the-future-v2.gif)

## Deploy to Cloud Run

[![Run on Google Cloud](https://deploy.cloud.run/button.svg)](https://deploy.cloud.run)

## How to install 

### prerequisites
Install Go in 1.16 version minimum.  
Install [Taskfile](https://taskfile.dev/#/installation) (optional)

### Build 

``` 
$ go build -o bin/kutego-api internal/main.go

// or 

$ task build
```

### Run app 

``` 
$ go run internal/main.go

// or 

$ task run
```

### Serve Swagger UI 

This will open you browser on Swagger UI
``` 
$ task swagger:serve
```
### View API up & running

Get list of available Cats:

```
$ curl http://localhost:8080/cats

[
  {
    "name": "5eme-element",
    "path": "5eme-element.gif",
    "url": "https://raw.githubusercontent.com/Pandahoro/cats/main/5eme-element.gif"
  },
  {
    "name": "arrow-cat",
    "path": "arrow-cat.gif",
    "url": "https://raw.githubusercontent.com/Pandahoro/cats/main/arrow-cat.gif"
  },

  [...]

  {
    "name": "back-to-the-future-v2",
    "path": "back-to-the-future-v2.gif",
    "url": "https://raw.githubusercontent.com/Pandahoro/cats/main/back-to-the-future-v2.gif"
  }
]
```
Filter by name: 
```
$ curl http://localhost:8080/cats?name=5eme-element
[
  {
    "name": "5eme-element",
    "path": "5eme-element.gif",
    "url": "https://raw.githubusercontent.com/Pandahoro/cats/main/5eme-element.gif"
  }
]
```

Get Cat by name:

```
$ curl -O localhost:8080/cat/back-to-the-future-v2

$ file back-to-the-future-v2
back-to-the-future-v2: PNG image data, 552 x 616, 8-bit/color RGBA, non-interlaced
```

Get Cat by name in medium size:  
(available values : `x-small`, `small`, `medium`)
```
$ curl -O localhost:8080/cat/back-to-the-future-v2?size=medium

$ file back-to-the-future-v2
back-to-the-future-v2: PNG image data, 269 x 300, 8-bit/color RGBA, non-interlaced
```

Get a Random cat:  
(available values : `x-small`, `small`, `medium`)
```
$ curl -O localhost:8080/cat/random?size=medium

$ file back-to-the-future-v2
back-to-the-future-v2: PNG image data, 269 x 300, 8-bit/color RGBA, non-interlaced
```

Enjoy to see a so cute Cat! 

## Notes

This API use [go-swagger](https://goswagger.io/install.html)

## Thanks 

Special thanks to [Scraly](https://github.com/Pandahoro) who wrote an awesome article about [How to write an API by example](https://dev.to/aurelievache/learning-go-by-examples-part-2-create-an-http-rest-api-server-in-go-1cdm)
