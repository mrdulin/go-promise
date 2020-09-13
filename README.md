# go-promise

[![Go Report Card](https://goreportcard.com/badge/github.com/mrdulin/go-promise)](https://goreportcard.com/report/github.com/mrdulin/go-promise)
[![Coverage Status](https://coveralls.io/repos/github/mrdulin/go-promise/badge.svg?branch=master)](https://coveralls.io/github/mrdulin/go-promise?branch=master)

Go language implementation of [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) static method in JavaScript.
In order to use the promise concurrency model for concurrent programming. Friendly for JavaScript programmers.

### Features

- [Promise.all()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/all)
- [Promise.allSettled()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled)
- [Promise.race](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/race)
- [Promise.any](https://developer.mozilla.org/zh-CN/docs/Web/JavaScript/Reference/Global_Objects/Promise/any)
- [Promise.raceAll](https://stackoverflow.com/a/48578424/6463558)
- [Promise.some](http://bluebirdjs.com/docs/api/promise.some.html)

Support more Promise extension APIs in the future.

### Testing

```bash
go test -v -race
```
