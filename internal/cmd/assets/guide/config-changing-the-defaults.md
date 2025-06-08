omm allows you to change the some of its behavior via configuration, which it
will consider in the order listed below:

- CLI flags (run `omm -h` to see details)
- Environment variables (eg. `$OMM_EDITOR`)
- A TOML configuration file (run `omm -h` to see where this lives; you can
    change this via the flag `--config-path`)

omm will consider configuration in the order laid out above, ie, CLI flags will
take the highest priority.
