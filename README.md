## EnvTrack

A simple utility that listens to the Consul key-value store (and perhaps etcd and ZooKeeper in the future) and then
sticks the result into a local file and (optionally) commits it in Git.

The intent is to be able to use a tool like envconsul but maintain auditabilty and reproducibility.  I'll likely add
an option to read the file and output curl commands, say, that will make Consul's state match the file.  With this you
would easily be able to do rollbacks of your environmental configuration, or do rollouts by using one of these files
to initialize Consul's state.