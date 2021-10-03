# crackmes.one
## About
This is the code of the website [crackmes.one](https://crackmes.one)

## Important informations
This website was created using the [GoWebApp](https://github.com/josephspurrier/gowebapp) from [Joseph Spurier](https://github.com/josephspurrier)

The really famous guy [Bonclay](https://twitter.com/mpgn_x64) from my team [Quokkalight](https://quokkalight.ninja) made this cool design.

## Setup for local developement.
1. Install `mongodb` for your choice of distribution.
2. Download the source code with go.

```sh
go get github.com/sar5430/crackmes.one
```

If you forked it, download it from yourself.

3.  go tho repository working dir (i.e. based of `go get`command)
```sh
cd go/src/github.com/sar5430/crackmes.one/
```

4. Modify the file app/route/route.go, uncomment

```golang
//return middleware(routes())
```

then comment

```golang
return http.HandlerFunc(redirectToHTTPS)
```

in `LoadHTTP()` function.

5. Make a `config` directory, and download this config file into it, and edit it to your liking. (You might want to edit the "Domain" value under "Session")

```sh
mkdir config
```

```sh
curl 'https://gist.githubusercontent.com/moex3/cb5225653a82dd1729525556e9175e92/raw/5fa39c308f09c1a1b44402305486bdc87fe1a61e/config.json' > config/config.json
```

6. Make a `tmp/crackme` and a `tmp/solution` directory.

```sh
mkdir -p tmp/{crackme,solution}`
````

7. Make a `static/crackme` and a `static/solution` directory.

```sh
mkdir -p static/{crackme,solution}
````

8. Build the binary.

```sh
go build
```

9. Install `python`, `zip` and `pymongo` if you want to run `validate.py`. (Also change the paths)

10. Run it.

```sh
./crackmes.one
```
