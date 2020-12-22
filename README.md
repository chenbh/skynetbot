# Usage

1. Install [ffmpeg](https://ffmpeg.org/), this will be needed to process audio files
1. Provide the required flags through command line or env vars:
  - `--admin-role` (`$DISCORD_ADMIN_ROLE`): ID for the role that will be required to perform some commands
  - `--token` (`$DISCORD_TOKEN`): Token for your bot, see https://discord.com/developers/applications


# Commands

Commands can be run in any channel the bot has read access to. Any message starting with `/` will be interpretted as a command.

- `/help`: display all available commands
- `/kill`: disable the bot, only works when run by user with admin role
- `/roll [NUM]: generate a random number between 0 (exclusive) and `NUM` (inclusive), in math notation this is [0, NUM). `NUM` defaults to 100.
- `/audio`:
  - `disconnect`: disconnect from current voice channel
  - `join`: join you current voice channel
  - `list`: list available sound clips
  - `play NAME`: plays a sound clip, see audio list
  - `remove NAME`: remove a sound clip, see audio list
  - `stop`: stop playing the current sound clip
  - `upload`: upload sound clip(s) using attachments, see audio list
  - `clip`: Creates a sound clip of the last 60 sec. of the current voice channel
- `/roles`
  - `assign @USER @ROLE`: assign a user to a role, requires admin role
  - `create NAME`: create a role, requires admin role
  - `join @ROLE`: join a role
  - `leave @ROLE`: leave a role
  - `list`: list all roles you currently have
  - `remove @ROLE`: remove a role, requires admin role
