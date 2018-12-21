# Did

> What did you do?

A simple command-line-tool to keep track of what you did.

## Installation

```sh
go get go.htdvisser.nl/did/cmd/did
```

## Usage

Register something you just did: `did [description]`

```sh
$ did publish my sideproject to Github
```

View what you did: `did [today|yesterday|monday|tuesday|wednesday|thursday|friday|saturday|sunday]`

```sh
$ did today
12:34: publish my sideproject to Github
```
