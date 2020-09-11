# go-promise

[![Go Report Card](https://goreportcard.com/badge/github.com/mrdulin/go-promise)](https://goreportcard.com/report/github.com/mrdulin/go-promise)

Go language implementation of [Promise](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise) static method in JavaScript.
In order to use the promise concurrency model for concurrent programming. Friendly for JavaScript programmers.

### Features

- [Promise.all()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/all)
- [Promise.allSettled()](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/allSettled)
- [Promise.race](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Promise/race)
- [Promise.raceAll](https://stackoverflow.com/a/48578424/6463558)

Support more Promise extension APIs in the future.

### Testing

```bash
go test -v -race
```
