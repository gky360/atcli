# atcli

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


## Installation

```
go get -u github.com/gky360/atcli
```


## Example usage

```
# List config
atcli config

# Set contest id and `atsrv` token
# and save config to ~/.atcli.yaml file.
# See https://github.com/gky360/atsrv for more details about `atsrv` token
atcli config -c arc090 -a xxxxxxxxxx

# Join contest
atcli join
# Create directories for tasks
atcli clone

# Build you source code
atcli build d
# Test your source code with sample cases downloaded by `atcli clone`
atcli test d               # always build your source code
atcli test d --skip-build  # skip build if possible
atcli test d 01            # run with a specified sample input

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

- fetching test cases used in judges on AtCoder
- templating your source code
- writing tests
