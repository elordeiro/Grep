#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
echo -n 'apple pie, apple and pie' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E "(\w+) (\w+), \1 and \2$"