#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
echo -n 'log' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E "^log"