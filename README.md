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
   0.0.1

COMMANDS:
     publish  Publish/Preview your public blueprint
     build    Build your blueprint
     help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --token value  apiary.io token
   --help, -h     show help
   --version, -v  print the version
```

# Build
Building changes will be available at `<source>\build\apiary.apib`
```
$ scruffy build --help

NAME:
   scruffy build - Build your blueprint

USAGE:
   scruffy build [command options] [arguments...]

OPTIONS:
   --config scruffy.yml  application configuration in yaml scruffy.yml
   --env value           Environment that have been set in config
   --watch false         Watch changes and reload on file change false

```

# Publish
Publishing changes to apiary.io
```
$ scruffy publish --help

NAME:
   scruffy publish - Publish/Preview your public blueprint

USAGE:
   scruffy publish [command options] [arguments...]

OPTIONS:
   --config scruffy.yml  application configuration in yaml scruffy.yml
   --env value           Environment that have been set in config
   --release false       Release changes in production doc false
   --watch false         Watch changes and reload on file change false
```

# Configuraion file
```
source: source
token: <YOUR_TOKEN>

environments:
  public:
        release: scruffypublic
        preview: scruffypublicpreview
        env:
          Title: Hello, user!
          Token: 519503441186ceb64b433cbc6455d2e7

  private:
        release: scruffyprivate
        preview: scruffyprivatepreview
        env:
          Title: Hello, world.
          Token: 519503441186ceb64b433cbc6455d2e7
```