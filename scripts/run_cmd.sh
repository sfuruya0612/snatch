#!/bin/sh

# RunCommand Test scripts
# e.g. snatch ssm run -f test/run_cmd.sh -t Role:test

# Uptime
uptime

# nginx status
systemctl status nginx

# Check 80 port listen
lsof -i:80
