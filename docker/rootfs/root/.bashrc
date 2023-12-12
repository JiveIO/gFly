RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

export TERM=xterm
# shell prompt: green host, blue current path
export PS1="\[$YELLOW\]\\h\[$NC\]\[$BLUE\] \\w \[$RED\]#\[$NC\] "
source /root/.aliases
