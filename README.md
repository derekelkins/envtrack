## EnvTrack

A simple utility that listens to the Consul key-value store (and perhaps etcd and ZooKeeper in the future) and then
sticks the result into a local file and (optionally) commits it in Git.

The intent is to be able to use a tool like envconsul but maintain auditabilty and reproducibility.  There is
an option to read the file and output curl commands that will make Consul's state match the file.  With this you
would easily be able to do rollbacks of your environmental configuration or do rollouts by using one of these files
to initialize Consul's state.

Feel free to create issues for desired features or for bugs, or, of course, pull requests.

### Example

```bash
envtrack -backend=git consul://consul.service.consul:8500/tracked/
```

Assuming the directory this is executed in is a Git repository, this will make a file called `config` with the values of all of the keys under `/tracked`.  Any time a key is updated, the `config` file will be updated and committed.

```bash
envtrack -script consul://consul.service.consul:8500/tracked/
```

This command will print a collection of `curl` commands to set the keys in Consul.  This is also useful to base64 decode the keys.
