# Skyra Deployment Pipeline — Glowdun

## Overview

Use Skyra's multi-being coordination to automate deployment for Gene's webapps. Three beings, each responsible for one stage. They coordinate through Skyra's threading system. Vercel handles hosting.

## Beings

### builder
- **Job:** Pull branch, run build, catch compile/lint errors
- **Device:** shell
- **Fails →** routes error back to Gene's AI via AI_COMMS.md

### crawler
- **Job:** Hit the Vercel preview URL, walk the frontend, flag broken pages/assets/console errors
- **Device:** MCP (browser tools)
- **Fails →** routes broken pages/errors back, blocks promotion

### deployer
- **Job:** Trigger Vercel preview deploy, promote to production once builder and crawler sign off
- **Device:** shell (vercel CLI) or Vercel API directly
- **Succeeds →** updates AI_COMMS.md with production URL

## Flow

1. Gene's AI opens a PR on `glowdun`
2. Pipeline triggers (webhook or cron)
3. `builder` pulls and builds — if errors, route back to Gene's AI
4. `deployer` triggers Vercel preview deploy
5. `crawler` crawls the preview URL via MCP browser tools
6. If clean, `deployer` promotes to production
7. Result posted to AI_COMMS.md

## What Needs to Be Built

1. **Shell device** — execute commands, return output (unblocks builder + deployer)
2. **MCP device** — connect to MCP server with browser tools (unblocks crawler)
3. **Trigger mechanism** — GitHub webhook listener or cron being that watches for PRs
4. **AI_COMMS.md integration** — beings write status updates back to the repo

## Stack

- **Runtime:** Skyra v.05
- **Hosting:** Vercel
- **Repo:** glowdun (GitHub)
- **Inference:** OpenRouter → Claude
- **Frontend testing:** MCP browser tools
