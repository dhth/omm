You can also import more than one task at a time by using the `import`
subcommand. For example:

```bash
cat << 'EOF' | omm import
orders: order new ACME rocket skates
traps: draw fake tunnel on the canyon wall
tech: assemble ACME jet-propelled pogo stick
EOF
```

omm will expect each line in stdin to hold one task's summary.
