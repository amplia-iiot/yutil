#!/usr/bin/env bash

CYAN='\033[0;36m'
RESET='\033[0m'

make -s check
printf '%b[%s] watching...%b\n' "$CYAN" "$(date "+%T")" "$RESET"
