#!/bin/bash
#
# Simple script to run web service for shortr
#

usage () {
  echo -e "\033[1mUsage:\033[0m $( basename $0 ) [-c CONFIG_FILE] [-l LOG_PATH] [stop]"
}

LOG_PATH_GIVEN=0
while [ $# -gt 0 ]; do
  case $1 in
      -c )
          if [[ -r "$2" && ! -d "$2" ]]; then
              CONFIG_FILE=$2
              source "$2"
              shift 2
          else
              echo -e "Unreadable config file »$2«."
              exit 1
          fi
          ;;
      -l )
          LOG_PATH_GIVEN=1
          LOG_PATH="$2"
          shift 2
          ;;
      stop)
          if [ -f $HOME/.shortr-pid ]; then
            kill -15 `cat $HOME/.shortr-pid`
            rm -f $HOME/.shortr-pid
            echo -e "Stopping service shortr"
            exit 0;
          else
            echo -e "Service shortr is not running. Nothing to stop"
            exit 2
          fi
          ;;
      * )
          echo -e "Unknown Option »$1«"
          usage
          exit 2
          ;;
  esac
done

if [[ $LOG_PATH_GIVEN -eq 1 ]]; then
  nohup ./shortr --serve 2> $LOG_PATH/.shortr.log < /dev/null &
else
  nohup ./shortr --serve < /dev/null &
fi

echo $! > $HOME/.shortr-pid

exit 0
