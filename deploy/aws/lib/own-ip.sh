#!/usr/bin/env bash

curl --silent 'https://api.ipify.org' | xargs -I {} printf '{ "myip_addr": "%s" }' {}
