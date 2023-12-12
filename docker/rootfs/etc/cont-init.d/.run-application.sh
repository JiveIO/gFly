#!/usr/bin/with-contenv bash

source /root/.bashrc
echo -e "${BLUE}---------- Run application ----------${NC}"

echo -e "${GREEN} 1. Start app with live mode ${NC}"
air main.go
