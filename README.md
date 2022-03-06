![example workflow](https://github.com/asahnoln/mesproc/actions/workflows/go.yml/badge.svg)

# What is it?

This is a Bot Story Processor. It was called `mesproc` because it was named `Message Processor` in the beginning. The repo name is going to be changed (probably).

## Why?

This project was created specifically for "Advanced Knitting Techniques" audioplay production by Kate Dzvonik. The Audience uses Telegram to follow the story.

## How does it work?

The common process is:

1. User sends text to the bot.
2. Bot sends it to the Story module.
3. Story module figures out if expected text for the current step in the story is right.
   1. If it's right - the Story returns right responses.
   2. If it's wrong - the Story returns fail responses.
4. Bot sends whatever the Story gives back to the user.
