#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
echo -n 'dogs1' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E "dogs?\d"