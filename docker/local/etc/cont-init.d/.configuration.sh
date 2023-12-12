#!/usr/bin/with-contenv bash

source /home/gfly/.bashrc
echo -e "${BLUE}--- Create common config ---${NC}"

# Generate .env file from template .env.docker
if [ ! -f .env ]; then
  cp .env.docker .env
fi