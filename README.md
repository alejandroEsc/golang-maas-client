# Another Golang MAAS API Client library

A lighter version fo the `MAAS API client library` originally from [https://github.com/juju/gomaasapi](https://github.com/juju/gomaasapi]).
The Developers there have done and continue doing amazing work. Please review their solution before using this.


Here we present a light version for my own purposes. Feel free to contribute. Because it is light
calls are more direct and there are/will be less helper code. We also leverage golang a bit more
making use of json marshalling and tags.


# Project Status
Alpha

Current status is alpha. Expect breaking changes for a while as things progress.


# Run sample code
First export variables that are particular to your maas setup, e.g.,

```
$ export MAAS_API_VERSION=2.0
$ export MAAS_API_URL=http://192.168.4.2:5240/MAAS/
$ export MAAS_API_KEY=<SOME_KEY>
```

and then you can simply run
```
$ go run ./cmd/sample.go
```

# Reference

[1] [https://github.com/juju/gomaasapi](https://github.com/juju/gomaasapi])