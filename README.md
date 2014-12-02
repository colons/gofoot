# gofoot

Runs as 'trenchfoot' in #NekoDesu on rizon, among other places.

## Setup

If you've got go [set up nicely][gosetup], you should be able to do the
following:

```bash
git clone https://github.com/colons/gofoot.git
cd gofoot
go get
cp config.go.example config.go
nano config.go
# (☝ﾟ∀ﾟ)☝ useful comments in that file
go build
```

At this point, you should be able to run `./gofoot robot [network]` to fire up
an instance of the actual IRC robot, or `./gofoot server` to serve the
documentation. `[network]` in the robot invocation is the string you use as the
key for your network's settings.

I've probably missed something in these instructions or the example config
comments, but I'm also pretty much guaranteed to be idling as colons on rizon,
and I don't bite (although I might be asleep), so poke me if you're lost.

[gosetup]: http://golang.org/doc/install
