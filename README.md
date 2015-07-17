#MongoAdmin

This is a tool for infrastructure guys to let those annoying engineers add their own indexes and run simple find and find by id queries.


## Requirements

Bower is required to build the UI and Go required to build the binary.

To install bower, have node and npm install and run `npm install -g bower`

## Install & Setup

Get the Go depedencies

`go get .` 

Install the bower dependencies

`bower install`

Compile the binary in your server environment.

`go build`


Use the `sample.toml` as an example for how to configure the databases. The `label` is what shows up in the UI. `database` is the actual DB name, and `connectionString` is the string to connect to the MongoDB server. You can use replica sets, etc, whatever `mgo` will accept in `mgo.Dial`.

## Running
Just run the binary and pass the relative path to your config file as the only argument.

`./mongoadmin sample.toml`