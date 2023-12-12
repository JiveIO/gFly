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
