# Spotify Top Ten

Designed to be hooked up to a crontab or something so that a designated playlist is constantly updated with a your top 10 songs.

How to run:
1. `git clone https://github.com/AntonyThomasCGI/spotify-top-ten.git`
2. Add required env variables to .env OR set environment variables however you would like.
3. `make build` and run the binary located dist/topten

First time you run it will prompt you to authorize through spotify and your web browser. Future runs will save your access tokens in ~/.spotify/auth.yaml so that the app requires no user input.
