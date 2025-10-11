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


I want you to shortly describe the purpose of the mechanics I will be implementing now and how it is used in Redis. Then I want you to explain the expected result and give me idea about which functions and structures I should implement, if it's not clear for the task itself. Don't give me ready code, just explain the idea of what I should achieve.

## 🧠 MyRedis Project Plan

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
|  | DEL | ☐ | Delete one or more keys |
|  | TYPE | ☐ | Return stored value type |
|  | EXISTS | ☐ | Check if key exists |
|  | KEYS | ☐ | Pattern matching (use `path.Match`) |
|  | FLUSHDB | ☐ | Clear all keys |
| **Expiry** | EXPIRE | ☐ | Attach TTL to keys |
|  | PEXPIRE | ☐ | Expiry in milliseconds |
|  | TTL / PTTL | ☐ | Query remaining lifetime |
|  | Key cleanup goroutine | ☐ | Periodically remove expired keys |
| **Data Structures – Strings** | INCR / DECR | ☐ | Numeric increment/decrement |
|  | APPEND | ☐ | Append to string |
| **Data Structures – Lists** | Create list | ☐ | Represent as `[]string` |
|  | RPUSH | ☐ | Append element |
|  | LPUSH | ☐ | Prepend element |
|  | LLEN | ☐ | Return list length |
|  | LRANGE | ☐ | Return element range |
|  | LPOP / RPOP | ☐ | Remove and return element |
| **Data Structures – Sets** | SADD | ☐ | Add members |
|  | SMEMBERS | ☐ | Get all members |
|  | SISMEMBER | ☐ | Check membership |
|  | SREM | ☐ | Remove members |
| **Data Structures – Hashes** | HSET / HGET | ☐ | Add hash support |
|  | HGETALL | ☐ | Return all fields |
| **Transactions** | INCR | ☐ | Atomic increment |
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
