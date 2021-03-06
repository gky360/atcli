# atcli

[![Build Status](https://travis-ci.org/gky360/atcli.svg?branch=master)](https://travis-ci.org/gky360/atcli)
[![Build status](https://ci.appveyor.com/api/projects/status/github/gky360/atcli?branch=master&svg=true)](https://ci.appveyor.com/project/gky360/atcli)
[![GoDoc](https://godoc.org/github.com/gky360/atcli?status.svg)](https://godoc.org/github.com/gky360/atcli)


```
        __           ___
       /\ \__       /\_ \    __
   __  \ \ ,_\   ___\//\ \  /\_\
 /'__'\ \ \ \/  /'___\\ \ \ \/\ \
/\ \L\.\_\ \ \_/\ \__/ \_\ \_\ \ \
\ \__/.\_\\ \__\ \____\/\____\\ \_\
 \/__/\/_/ \/__/\/____/\/____/ \/_/
```

Command line interface for AtCoder (unofficial)


## Requirements

- [atsrv](https://github.com/gky360/atsrv)


## Installation

```
go get -u github.com/gky360/atcli
```


## Example usage

```
# Configuration
# See also https://github.com/gky360/atsrv for more details about ATSRV_AUTH_TOKEN
export ATSRV_USER_ID=your_atcoder_user_id
export ATCLI_ROOT=~/atcoder
export ATCLI_CPP_TEMPLATE_PATH=$ATCLI_ROOT/templates/Main.cpp
export ATSRV_AUTH_TOKEN=$(cat /dev/urandom | base64 | fold -w 32 | head -n 1)


# Join contest
atcli join
# Create directories for tasks
atcli clone

# Build your source code
atcli build d
# Test your source code with sample cases downloaded by `atcli clone`
atcli test d               # always build your source code
atcli test d --skip-build  # skip build if possible
atcli test d sample01      # run with a specified sample input

# Submit your source code to AtCoder
atcli submit d

# Get info from AtCoder
atcli get
atcli get contest
atcli get task [d]
atcli get submission [2167890 | -t d]
```

See `atcli --help` for more deatils.


## Future work

- ~~fetching test cases used in judges on AtCoder~~ -> Done in v0.1.0
- ~~templating your source code~~ -> Done in v0.1.0
- writing tests
