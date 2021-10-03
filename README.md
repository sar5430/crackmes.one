# crackmes.one
## About
This is the code of the website [crackmes.one](https://crackmes.one)

## Important informations
This website was created using the [GoWebApp](https://github.com/josephspurrier/gowebapp) from [Joseph Spurier](https://github.com/josephspurrier)

The really famous guy [Bonclay](https://twitter.com/mpgn_x64) from my team [Quokkalight](https://quokkalight.ninja) made this cool design.

## Setup for local developement.
* Install `mongodb` for your choice of distribution.
* Download the source code with go. `go get github.com/sar5430/crackmes.one`. If you forked it, download it from yourself.
* `cd go/src/github.com/0xstan/crackmes.one/`
* Modify the file app/route/route.go, uncomment 
```
//return middleware(routes())
```  
then comment  
```
return http.HandlerFunc(redirectToHTTPS)
```
in `LoadHTTP()` function.
* Make a `config` directory, and download this config file into it, and edit it to your liking. (You might want to edit the "Domain" value under "Session")
`curl 'https://gist.githubusercontent.com/moex3/cb5225653a82dd1729525556e9175e92/raw/5fa39c308f09c1a1b44402305486bdc87fe1a61e/config.json' > config/config.json`
* Make a `tmp/crackme` and a `tmp/solution` directory. `mkdir -p tmp/{crackme,solution}`
* Make a `static/crackme` and a `static/solution` directory. `mkdir -p static/{crackme,solution}`
* Build the binary. `go build`
* Install `python`, `zip` and `pymongo` if you want to run `validate.py`. (Also change the paths)
* Run it. `./crackmes.one`
