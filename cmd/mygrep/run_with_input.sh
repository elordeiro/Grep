#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
cd cmd/mygrep && echo -n 'apple' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E \"\\d\"