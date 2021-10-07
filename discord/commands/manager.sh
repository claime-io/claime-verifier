#!/bin/bash

BOT_TOKEN=
APP_ID=892321773981421568
TEST_GUILD_ID=892441777808765049
BASE_URL="https://discord.com/api/v8/applications/${APP_ID}"
ENDPOINT="${BASE_URL}/guilds/${TEST_GUILD_ID}/commands"

if [ -z "${BOT_TOKEN}" ]; then
  echo "BOT_TOKEN is required."
  exit 1;
fi

OPTIONS='gld:'

while getopts "${OPTIONS}" option; do
  case "$option" in
    g)
      ENDPOINT="${BASE_URL}/commands"
      echo "global"
      ;;
  esac
done

OPTIND=1

while getopts "${OPTIONS}" option; do
  case "$option" in
    l)
      curl -X GET -H "Content-Type: application/json" \
      -H "Authorization: Bot ${BOT_TOKEN}" \
      "${ENDPOINT}"
      exit 0
      ;;
    d)
      COMMAND_ID=$OPTARG
      ENDPOINT="${ENDPOINT}/${COMMAND_ID}"
      curl -X DELETE -H "Content-Type: application/json" \
        -H "Authorization: Bot ${BOT_TOKEN}" \
        "${ENDPOINT}"
      exit 0
      ;;
  esac
done
shift $((OPTIND - 1))

command=$1
if [ -z "${command}" ]; then
  echo "command name is required."
  grep -rl . | grep json | sed -e "s/\.json//"
  exit 1;
fi

body=$(cat "${command}.json" | tr -d '\n')

curl -X POST -H "Content-Type: application/json" \
-H "Authorization: Bot ${BOT_TOKEN}" \
-d "${body}" \
"${ENDPOINT}"
