#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
CYAN='\033[0;36m'
RESET='\033[0m'

( ( (make -s test) && printf '%bTests - OK%b\n' "$GREEN" "$RESET" ) || printf '%bTests - failed%b\n' "$RED" "$RESET" )
( ( (make -s lint) && printf '%bLint  - OK%b\n' "$GREEN" "$RESET" ) || printf '%bLint  - failed%b\n' "$RED" "$RESET" )
printf '%bOn    - %s%b\n' "$CYAN" "$(date)" "$RESET"
