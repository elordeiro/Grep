#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
echo -n 'cog' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E "^d.g"