# go-build-my-own-redis

Steps for bulding my Redis:
1. Simple server
2. RESP parser
3. Value Marshal and Writer
4. Add simple commands (PING, ECHO)
5. Add simple storage
6. Add SET and GET commands

Test using redis-cli
```
docker run -it --rm redis redis-cli -h host.docker.internal -p 6380
```
