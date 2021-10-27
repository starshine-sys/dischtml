#!/bin/sh
npx tailwindcss -i ./discord.css -o ./embed/discord.css
uglifycss embed/discord.css --output embed/discord.css
