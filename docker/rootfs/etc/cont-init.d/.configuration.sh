#!/usr/bin/with-contenv bash

source /root/.bashrc
echo -e "${BLUE}---------- Configuration ----------${NC}"

echo -e "${GREEN} 1. Generate .env file from template .env.docker ${NC}"
if [ ! -f .env ]; then
  cp .env.docker .env
fi