#!/usr/bin/env bash

set -e

echo "Removing cortex..."
sudo rm -f /usr/local/bin/cortex
rm -rf ~/.cortex

echo "cortex has been successfully uninstalled."
