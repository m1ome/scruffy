# Scruffy
[![Build Status](https://travis-ci.org/m1ome/scruffy.svg?branch=master)](https://travis-ci.org/m1ome/scruffy)
[![Coverage Status](https://coveralls.io/repos/github/m1ome/scruffy/badge.svg?branch=master)](https://coveralls.io/github/m1ome/scruffy?branch=master)
> Simple [ApiBlueprint](apiblueprint.com) builder.

![Scruffy](http://vignette4.wikia.nocookie.net/en.futurama/images/1/10/Scruffy.png/revision/latest?cb=20170123190905)

# Installation
```
go get github.com/m1ome/scruffy
```

# Usage
```
$ scruffy --help

NAME:
   Scruffy - build your blueprints from mess to order!

USAGE:
   scruffy [global options] command [command options] [arguments...]

VERSION:
   0.0.0

COMMANDS:
     publish  Publish/Preview your public blueprint
     build    Build your blueprint
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value  apiary.io token
   --help, -h     show help
   --version, -v  print the version

```

# Configuraion file
```
source: src
token: your_token_in_here

public:
  name: public
  preview: publicpreview
  env:
    Title: Public title
    Name: Public Name

private:
  name: private
  env:
    Title: Private title
    Name: Private Name
```