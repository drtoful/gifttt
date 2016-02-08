# gifttt

gifttt is a simple rule engine in go that uses Lisp S-Expressions to define its rules.

## Building

    go get github.com/drtoful/gifttt

## Running

You can run gifttt by invoking its binary in your $GOPATH/bin. This will read in all rule files in the current directory and start the API server on port 4200.

### Command line options

    -cpuprofile string
          write cpu profile to file
    -db string
          path to the database store (default "gifttt.db")
    -ip string
          ip to bind the api server to
    -port string
          port for api server (default "4200")
    -ruledir string
          path to rule files (default "./")

## Quick start

Have a look in doc/quick.md for a small tutorial on how to operate with gifttt.

## License

gifttt is licensed under the BSD License. See LICENSE for more information

### Third-Party Libraries

* [bolt](github.com/boltdb/bolt): MIT License
* [negroni](github.com/codegangsta/negroni): MIT License
* [context](github.com/gorilla/context): BSD License
* [mux](github.com/gorilla/mux): BSD License
* [twik](github.com/drtoful/twik): LGPLv3
