# coredns-recorder

Provides a CoreDNS plugin for recording queries.

By default the plugin will log to the CoreDNS log plugin.

You can optionally configure it to log to NATS.

If you log to NATS, you can also use the provided tool [`cmd/rec`](cmd/rec) to listen to NATS and log the records in CSV format to stdout or to a file.
Additional info in [doc/README.md](doc/README.md)
