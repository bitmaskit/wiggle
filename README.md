# wiggle

A simple Go tool to keep your computer awake by moving the mouse in a circular motion after a period of inactivity.

## Installation

```sh
go install github.com/bitmaskit/wiggle@latest
```

## Usage

You can set the idle timeout (in seconds) using either an environment variable or a command-line argument:

```sh
# Using environment variable
IDLE_TIMEOUT=10 wiggle

# Using command-line argument
wiggle 10
```
### If neither is set, the default idle timeout of 60 seconds will be used.
