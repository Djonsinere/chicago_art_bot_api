usage:
  add your tg_bot TOKEN
  
  cmd: go build . && go run chicago_bot
prod:
 docker build -t my-bot .
 docker run -d -v $(pwd)/chicago_users.db:/usr/src/app/chicago_users.db my-bot
 docker update --restart unless-stopped $(docker ps -q)
