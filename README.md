## Configured with

- [go-gorp](https://github.com/go-gorp/gorp): Go Relational Persistence
- [jwt-go](https://github.com/golang-jwt/jwt): JSON Web Tokens (JWT) as middleware
- [go-redis](https://github.com/go-redis/redis): Redis support for Go
- Go Modules
- Feature **PostgreSQL 12** with JSON/JSONB queries & trigger functions
- Enviroment support

### Installation

```
$ go get github.com/Zayank/group-chats
```

```
$ cd $GOPATH/src/github.com/Zayank/group-chats
```

```
$ go mod init
```

```
$ go install
```

You will find the **database.sql** in `db/database.sql`

And you can import the postgres database using this command:

```
$ psql -U postgres -h localhost < ./db/database.sql
```

> Make sure to change the values in .env for your databases

```
$ go run *.go
```

## Building Your Application

```
$ go build -v
```

```
$ ./group-chats
```
## Credit

https://github.com/gin-gonic/examples/tree/master/realtime-chat
https://github.com/Massad/gin-boilerplate

---

## License

(The MIT License)

Permission is hereby granted, free of charge, to any person obtaining
a copy of this software and associated documentation files (the
'Software'), to deal in the Software without restriction, including
without limitation the rights to use, copy, modify, merge, publish,
distribute, sublicense, and/or sell copies of the Software, and to
permit persons to whom the Software is furnished to do so, subject to
the following conditions:

The above copyright notice and this permission notice shall be
included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED 'AS IS', WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
