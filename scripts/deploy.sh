#!/bin/bash
set -e

MINI="skyra@10.0.0.50"

echo "deploying to $MINI..."

# read keys from local keychain
echo "reading keys from keychain..."
OPENROUTER_KEY=$(security find-generic-password -s "OpenRouter_API_KEY" -a "mikepersonal" -w)

# push local changes first
echo "pushing to github..."
git push

# write keys to restricted .env on mini
echo "syncing keys..."
printf 'OPENROUTER_API_KEY=%s\n' "$OPENROUTER_KEY" | ssh "$MINI" "cat > ~/skyra/.env && chmod 600 ~/skyra/.env"

# pull, build, restart on the mini
ssh "$MINI" bash -ls <<'REMOTE'
set -e
cd ~/skyra
git pull --ff-only

cd skyra-v.05
echo "building..."
go build -o skyra-bin .
echo "build ok"

# stop existing instance if running
if pgrep -f "skyra-bin" > /dev/null 2>&1; then
    echo "stopping existing instance..."
    pkill -f "skyra-bin"
    sleep 1
fi

echo "done. run manually with: cd ~/skyra/skyra-v.05 && ./skyra-bin"
REMOTE

echo "deployed."
