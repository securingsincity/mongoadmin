This is a tool for infrastructure guys to let those annoying engineers add their own indexes.

We'll add more functionality as we ID more use cases.

# Setup
Just compile the binary in your server environment. 

Use the `sample.toml` as an example for how to configure the databases. The `label` is what shows up in the UI. `database` is the actual DB name, and `connectionString` is the string to connect to the MongoDB server. You can use replica sets, etc, whatever `mgo` will accept in `mgo.Dial`.

# Running
Just run the binary and pass the relative path to your config file as the only argument. Make sure you run it behind some sort of authentication wall, since there's no auth built-in. You can run it behind nginx running basic auth, for example. Maybe only in your network. Whatever you want. Just be careful.