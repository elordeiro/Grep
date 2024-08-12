#!/bin/zsh
export PATH=$PATH:/usr/local/go/bin
echo -n 'sally has 12 apples' | dlv debug --headless --listen=:2345 --api-version=2 --accept-multiclient . -- -E "\d\\d\\d apples"