# Skyra Explained (For Friends)

## What Skyra Is

Skyra is a personal AI assistant I am building to help with day-to-day life and projects.

Think of it like a custom "Jarvis" that can:

- listen and talk naturally
- remember important context
- help with coding, planning, and problem-solving
- run tasks on different machines depending on how heavy the work is

## The Simple Version of How It Works

Skyra has 3 "brains" working together:

1. **Fast Listener (Raspberry Pi)**
- Always on.
- Hears wake word, converts speech to text, gives quick responses.
- Great for speed and instant feedback.

2. **Main Brain (Mac mini)**
- The decision center.
- Handles memory, planning, and task coordination.
- Makes the final/authoritative decisions.

3. **Heavy Brain (GPU Machine)**
- Used only for hard, expensive thinking.
- Handles deeper reasoning and bigger coding/analysis tasks.

## Why It Feels Fast

Skyra can give a quick "provisional" response from the fast listener, then ask the main brain to verify and finish the job.

So in practice:

- you get a fast response right away
- then you get the final, correct answer from the main system

## Why This Matters

Most assistants are generic. Skyra is personal and built for real daily use.

Goals:

- reduce mental load
- keep momentum on work/school/projects
- make complex tasks feel manageable
- be reliable over the long run

## Current Status

Skyra is actively being built and tested.

What already exists:

- architecture and reliability design
- voice listener service
- control-plane API foundation
- event reliability pipeline (so requests are not lost)

What is being improved next:

- better task formation and routing
- stronger memory and context handling
- smoother voice experience and faster response quality

## In One Sentence

Skyra is my personal, always-improving AI partner designed to help me think clearer, build faster, and stay consistent.
