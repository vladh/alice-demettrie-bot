# Alice DeMettrie

A simple Telegram bot for vladh.

Works on Linux and Windows. On Windows, requires `nircmd` to be in your `PATH`.

## How to run

Create a `config.go` file with the following format:

```
package main

const TOKEN = "..."
```

Compile with:

```
make
```

If you want Alice to run on startup, install her and add her to systemd:

```
sudo cp alice /usr/bin/
sudo cp alice.service /etc/systemd/user/
```

## Commands

```
* setvol xx%: Set volume to xx percent.
* sleep: Put the comptuer into sleep mode
* help: Print the help text
```
