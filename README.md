# go-build-my-own-redis

This repository is an educational project to improve my Go skills by recreating Redis from scratch.

Instead of following a course, I built my own learning path using ChatGPT as a tutor and course outlines from various resources.

In sections below, I described:

- Implementation log — step-by-step history of features I was adding
- Prompts — examples of propmpts I used for ChatGPT to generate tasks and explanations
- Task plan — roadmap features and improvements
- Testing setup

## How I did this learning project

Prompts for guiding ChatGPT:

> I am building Redis clone in Go as an excercise to learn the Go language (already did a tour of go). Here is the plan that I created with ChatGPT and the code that I implemented so far:
...
Review the plan, check if there is something that would be halpful to add.
If not, let's continue with next task:
...

> I am building Redis clone in Go as an excercise to learn the Go language (already did a tour of go). Here is the plan that I created with ChatGPT and the code that I implemented so far:
...
Let's continue with next task according to the plan:
...

> I want you to shortly describe the purpose of the mechanics I will be implementing now and how it is used in Redis. Then I want you to explain the expected result and give me idea about which functions and structures I should implement, if it's not clear for the task itself. Don't give me ready code, just explain the idea of what I should achieve.


## Log of the steps I followed

Steps for bulding my Redis:
1. Simple server
2. RESP parser
3. Value Marshal and Writer
4. Added simple commands PING, ECHO
5. Added simple storage
6. Added SET and GET commands
7. Added Server and Graceful shutdown
8. Added commands DEL, TYPE, EXISTS, KEYS, FLUSHDB
9. Added expiration mechanincs - EXPIRE, PEXPIRE, TTL, PTTL, clean-up routine
10. Added Engine for Command Registery and Dispatcher
11. Added commands INCR, DECR, INCRBY, DECRBY, APPEND
11. Added list structure: LPUSH, RPUSH, LPOP, RPOP, LLEN, LRANGE

## Task Plan

| Category | Task | Status | Notes |
|-----------|------|--------|-------|
| **Networking** | Bind to a port | ✅ | Basic TCP listener |
|  | Respond to PING | ✅ |  |
|  | Respond to multiple PINGs | ✅ |  |
|  | Handle concurrent clients | ✅ | Using goroutines |
|  | Graceful shutdown | ✅ | Handle `SIGINT` / `SIGTERM` cleanly |
| **Protocol (RESP)** | Parse RESP arrays | ✅ |  |
|  | Write RESP replies | ✅ |  |
|  | Support inline commands | ☐ | e.g., `PING\r\n` without array syntax |
| **Basic Commands** | ECHO | ✅ |  |
|  | PING | ✅ |  |
|  | SET | ✅ |  |
|  | GET | ✅ |  |
|  | DEL | ✅ | Delete one or more keys |
|  | TYPE | ✅ | Return stored value type |
|  | EXISTS | ✅ | Check if key exists |
|  | KEYS | ✅ | Pattern matching (use `path.Match`) |
|  | FLUSHDB | ✅ | Clear all keys |
| **Expiry** | EXPIRE | ✅ | Attach TTL to keys |
|  | PEXPIRE | ✅ | Expiry in milliseconds |
|  | TTL / PTTL | ✅ | Query remaining lifetime |
|  | Key cleanup goroutine | ✅ | Periodically remove expired keys |
| **Engine Architecture** | **Command Registry / Dispatcher** | ✅      | Map commands dynamically instead of using a large `switch`; each command registered with metadata (name, arity, handler) |
| **Data Structures – Strings** | INCR / DECR | ✅ | Numeric increment/decrement |
|  | APPEND | ✅ | Append to string |
| **Data Structures – Lists** | Create list | ✅ | Represent as `[]string` |
|  | RPUSH | ✅ | Append element |
|  | LPUSH | ✅ | Prepend element |
|  | LLEN | ✅ | Return list length |
|  | LRANGE | ✅ | Return element range |
|  | LPOP / RPOP | ✅ | Remove and return element |
| **Data Structures – Sets** | SADD | ☐ | Add members |
|  | SMEMBERS | ☐ | Get all members |
|  | SISMEMBER | ☐ | Check membership |
|  | SREM | ☐ | Remove members |
| **Data Structures – Hashes** | HSET / HGET | ☐ | Add hash support |
|  | HGETALL | ☐ | Return all fields |
| **Transactions** | INCR | ☐ | Atomic increment |
|  | DECR | ☐ | Atomic decrement |
|  | MULTI / EXEC / DISCARD | ☐ | Transaction support |
| **Persistence (AOF)** | Write AOF on write commands | ☐ | Append-only log |
|  | Replay AOF on startup | ☐ | Load data back |
|  | AOF rewrite (compaction) | ☐ | Reduce file size |
| **Persistence (RDB)** | Save snapshot to file | ☐ | Serialize `Storage` |
|  | Save expiry info | ☐ | Store TTLs |
|  | Load snapshot on startup | ☐ | Optional |
| **Pub/Sub** | SUBSCRIBE | ☐ | Add subscriber registry |
|  | PUBLISH | ☐ | Send message to subscribers |
|  | UNSUBSCRIBE | ☐ | Clean up |
| **Replication** | Master/Replica mode | ☐ | Add `--replicaof` support |
|  | INFO command | ☐ | Server stats |
|  | Replication handshake | ☐ | Basic sync logic |
|  | WAIT | ☐ | Wait for replicas to acknowledge writes |
|  | ACK (replica acknowledgment) | ☐ | Replicas confirm write receipt |
| **Sorted Sets (ZSet)** | ZADD | ☐ | Add member with score |
|  | ZRANGE | ☐ | Return ordered elements |
| **Geospatial** | GEOADD | ☐ | Store coordinates |
|  | GEOPOS | ☐ | Return positions |
| **Server** | INFO command | ☐ | Server info, memory, clients |
|  | CONFIG GET / SET | ☐ | Runtime configuration |
|  | COMMAND | ☐ | Describe supported commands |
| **Testing / Utilities** | Unit tests for RESP parsing | ☐ | Use Go test framework |
|  | Integration tests with `redis-cli` | ☐ | `redis-cli -p 6380` |
|  | Benchmarking (`go test -bench`) | ☐ | Compare with real Redis |
| **Extras (optional)** | AUTH command | ☐ | Basic authentication |
|  | SELECT databases | ☐ | Multiple DBs (0–15) |
|  | Logging improvements | ☐ | Add timestamps / structured logs |
|  | Metrics / Prometheus | ☐ | Monitor ops/sec, memory |
|  | CLI client in Go | ☐ | Mini `redis-cli` clone |

## Testing

Test using simple `echo` and `printf` (following the expected Redis syntax):
```
echo -e "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n" | nc localhost 6380

printf "*3\r\n$3\r\nSET\r\n$3\r\nkey\r\n$5\r\nvalue\r\n" | nc localhost 6380
```

Test using redis-cli:
```
docker run -it --rm redis redis-cli -h host.docker.internal -p 6380
```
