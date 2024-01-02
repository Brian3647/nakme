# nakme

_simple, fast & secure project sharing platform_

## what?

nakme is a platform created for sharing projects with others, giving & receiving feedback and finding collaborators.

## why?

unless it's something very, very big, normally it's hard to find people to work with you on your OSS projects. When multiple people collaborate on the same thing, programming becomes more fun, faster and more efficient. This app is not meant to be a repository hosting service, but rather a sharing platform.

## how?

the backend is built with [go](https://golang.org/) (specifically [echo](https://echo.labstack.com/)). The frontend is yet to be started, but it will probably be build with something very simple like [vue](https://vuejs.org/) or go templates + htmx.

## how to run

### prerequisites

- [go](https://golang.org/)
- [postgresql](https://www.postgresql.org/)
- any smtp service (eg. gmail)

### steps

1. clone the repo
2. `cd` into the repo's `db` directory
3. change `setup.go` to your needs (password, port, etc.) and run it with `go run setup.go`
4. cd back into the main repo's directory
5. copy `.env.template` to `.env` and add your credentials
6. run `go run main.go`

Done! you can now open `localhost:3000` in your browser.

## project status

see [TODO.md](TODO.md)

## contributing

see [CONTRIBUTING.md](CONTRIBUTING.md)

## project structure

```js
.
├── ... <- configuration files, licenses and other miscelaneous files
├── `pkg`
│   ├── `api` <- the api
│   │   └── `auth` <- authentication routes
│   ├── `db` <- database stuff
│   ├── `util` <- utility functions
│   └── `web` <- the frontend TODO
├── `db` <- database initial configurations
└── `docs` <- api docs
```

## license

this project is licensed under either the [MIT license](LICENSE-MIT) or the [Apache 2.0 license](LICENSE-APACHE), at your option.

## api usage & documentation

all the api docs can be found in [docs/README.md](docs/README.md)

```
happy coding
```
