RED='\033[1;31m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
BLUE='\033[1;34m'
NC='\033[0m' # No Color

export TERM=xterm
# shell prompt: green host, blue current path
export PS1="\[$GREEN\]\\h\[$NC\]\[$BLUE\] \\w \\$\[$NC\] "
source /root/.aliases
source /home/gfly/.aliases

echo -e "${BLUE}---------- Build application ----------${NC}"
echo -e "${GREEN} 1. Run command \"build_app\" to build app ${NC}"
echo -e "${GREEN} 2. Run command \"build_artisan\" to build CLI tool ${NC}"

echo -e "${BLUE}---------- Release application ----------${NC}"
echo -e "${GREEN} 1. Run command \"release_windows_64\" to build for Windows ${NC}"
echo -e "${GREEN} 2. Run command \"release_mac_amd64\" to build for Mac (Intel chip) ${NC}"
echo -e "${GREEN} 3. Run command \"release_mac_arm64\" to build for Mac (Silicon chip) ${NC}"
echo -e "${GREEN} 4. Run command \"release_linux_amd64\" to build for Linux (AMD) ${NC}"
echo -e "${GREEN} 5. Run command \"release_linux_arm64\" to build for Linux (ARM) ${NC}"
