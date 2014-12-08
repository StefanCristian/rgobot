GoBot
========


## Description

> GoBot is darthlukan's work with a few personal tweaks.

## Installation
```
    $ go get github.com/blacknoxis/rgobot
```
> Or:

```
    $ mkdir -p $GOPATH/src/github.com/blacknoxis
    $ cd $GOPATH/src/github.com/blacknoxis
    $ git clone https://github.com/BlackNoxis/rgobot.git
    $ cd gobot
```

## Usage

> Edit the config.json file located in $GOPATH/src/github.com/blacknoxis/gobot to your preferences.

> After those variables have been edited, you can run:
```
    $ go install .    # Note the '.'
    $ gobot
```

## in-channel interaction

> For now, there aren't really any "real" commands, but if you prefix with "!", the bot will spit out a message.

> Example:
```
    !slap SomeUser really hard
    >> *$botNick slaps SomeUser really hard, FOR SCIENCE!
```

## TODO

- 1. Add commands that actually do something useful
- 2. Google Search
- 3. Logging
- 4. Tests would be nice >.>

## License

> GPLv2, see LICENSE file.
