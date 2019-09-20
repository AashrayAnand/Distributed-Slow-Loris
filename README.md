# Distributed Slow Loris

A distributed slow loris is an amalgamation of two well-known denial of service attack techniques. It utilizes a distributed architecture, where a master process instructs slave processes on an endpoint to attack. In turn, the slave processes execute the slow loris denial of service technique, initiating connections with an HTTP server, and opening, but never completing a single HTTP request. This is an especially powerful technique against servers which initiate a new thread to handle every connection (Apache).

## System Design

(version 1)[https://github.com/AashrayAnand/Distributed-Slow-Loris/tree/0af0743a584ee2ca8fa2f1a5cae69e419fe29b7a]
