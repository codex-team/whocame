# Whocame – bot that knows who come to the office

## Installation

With local Go:
```
go build -o build .
```

With docker:
```
docker-compose build
docker-compose up
```

## Configuration

Use sample from `storage.yml`:
```
interface: en0
logLevel: debug
webhook: https://notify.bot.codex.so/u/ABCDEF00
goneAfter: 30m
members:
  n0str:
    - 00:00:5e:00:53:af
    - FF:C1:87:7D:A2:5F
  khaydarovm:
    - 2c-54-91-88-c9-e3
```

* interface - interface to listen
* webhook - CodeX Bot webhook to the chat
* goneAfter - Make person left after 30m
* Members - nicknames and their MACs

## Running

You can run it manually, but it's better to use Systemd.
```
chmod +x ./whocame
./whocame
```