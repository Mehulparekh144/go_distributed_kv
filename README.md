
# Distributed Scalable Key-Value Store

## Features Implemented So Far

### 1. Basic Key-Value Store

* Simple in-memory storage of key-value pairs with TTL (time-to-live) support.
* Supports `PUT`, `GET`, and `DELETE` operations over HTTP.

### 2. Consistent Hashing for Node Distribution

* Uses CRC32 hashing to assign nodes to positions on a hash ring.
* Keys are assigned to nodes based on their hashed position on the ring.
* Dynamically handles multiple nodes in the system by mapping keys to responsible nodes.

### 3. Request Forwarding

* If a node receives a request for a key it’s not responsible for, it forwards the request to the correct node.
* Supports forwarding for `PUT`, `GET`, and `DELETE` requests to ensure correctness in a distributed environment.

### 4. Mutex Locks for Concurrent Access

* Uses mutex locking in the key-value store to handle concurrent requests safely and prevent race conditions.

---

## How to Run / Test

* Run multiple server instances on different ports to simulate a cluster.
* Use HTTP clients (like `curl` or Postman) to send requests to any node.
* Requests are forwarded to the correct node based on consistent hashing.

---

## Future Enhancements (Planned)

* Replication for fault tolerance and availability.
* Retry and fallback logic for node failures.
* Persistent storage to prevent data loss on node restart.
* Health checks and dynamic node membership management.
