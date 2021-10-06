#!/bin/bash

BOT_TOKEN=
APP_ID=892321773981421568
GUILD_ID=
ENDPOINT="https://discord.com/api/v8/applications/${APP_ID}/guilds/${GUILD_ID}/commands"

while getopts ":g" optKey; do
  case "$optKey" in
    g)
      ENDPOINT="https://discord.com/api/v8/applications/${APP_ID}/commands"
      ;;
  esac
done
shift $((OPTIND - 1))

command=$1
if [ -z "${BOT_TOKEN}" ]; then
  echo "BOT_TOKEN is required."
  exit 1;
fi
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
