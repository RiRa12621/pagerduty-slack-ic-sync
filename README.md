# pagerduty-slack-ic-sync
Sync a given slack handle to update to who is oncall

## Installation

You can choose from the following three options:
1. go install
2. run in a container
3. build from source

### Go install

The go way, requires you to have a working go installation. After that you 
can run the following command:

```shell
go install github.com/RiRa12621/pagerduty-slack-ic-sync@latest
```

If you have you $GOPATH added to your $PATH, you will be able to use it 
right away. Otherwise you still need to move it to your path.

### Container

There is a container image available for you to use like so:

```shell
docker run -ti --rm quay.io/rira12621/pagerduty-slack-ic-sync <flags>
```

### From Source

You can build the tool from source and then run it:

```shell
git clone git@github.com:RiRa12621/pagerduty-slack-ic-sync.git
cd pagerduty-slack-ic-sync
go build -o pagerduty-slack-ic-sync
./pagerduty-slack-ic-sync
```

## Running

This tool requires a chunk of parameters that are all mandatory:

* --schedule=<ID of the PD Schedule>
* --pd-token=<Pagerduty access token>
* --slack-toke=<slack token to update the alias>
* --alias=<alias to update>


With those parameters set, you can run the tool:

```shell
pagerduty-slack-ic-sync --schedule=ABCD --pd-token=asda123541asdf --slack-token=asdf987987asdf987asdf --alias=ASDFF
```
