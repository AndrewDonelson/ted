#!/bin/bash
# Clean test output script

# Ensure terminal is reset on exit
trap 'stty sane 2>/dev/null; tput reset 2>/dev/null' EXIT INT TERM

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

printf "${BOLD}${CYAN}Running tests...${NC}\n\n"

# Run tests and process output line by line
go test ./... -count=1 2>&1 | while IFS= read -r line; do
    # Clear any carriage returns
    line=$(echo "$line" | tr -d '\r')
    
    if [[ $line == ok* ]]; then
        package=$(echo "$line" | awk '{print $2}')
        time=$(echo "$line" | awk '{print $NF}')
        printf "${GREEN}✓${NC} %-50s %s\n" "$package" "$time"
    elif [[ $line == FAIL* ]]; then
        package=$(echo "$line" | awk '{print $2}')
        printf "${RED}✗ FAIL${NC} %-50s\n" "$package"
    elif [[ $line == "?"* ]]; then
        package=$(echo "$line" | awk '{print $2}')
        printf "${YELLOW}○${NC} %-50s [no tests]\n" "$package"
    elif [[ $line == *FAIL:* ]] || [[ $line == *Error:* ]]; then
        printf "${RED}  %s${NC}\n" "$line"
    fi
done

# Check exit code
EXIT_CODE=${PIPESTATUS[0]}

echo ""
if [ $EXIT_CODE -eq 0 ]; then
    echo -e "${GREEN}${BOLD}✓ All tests passed${NC}"
else
    echo -e "${RED}${BOLD}✗ Tests failed${NC}"
fi

exit $EXIT_CODE
