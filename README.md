# Distributed Slow Loris

A distributed slow loris is an amalgamation of two well-known denial of service attack techniques. It utilizes a distributed architecture, where a master process instructs slave processes on an endpoint to attack. In turn, the slave processes execute the slow loris denial of service technique, initiating connections with an HTTP server, and opening, but never completing a single HTTP request. This is an especially powerful technique against servers which initiate a new thread to handle every connection (Apache).

## System Design

[version 0.1](https://github.com/AashrayAnand/Distributed-Slow-Loris/tree/0af0743a584ee2ca8fa2f1a5cae69e419fe29b7a):

Version 0.1 follows a single worker architecture, where attacks are manually executed from my local machine. These attacks provision a set of attacker threads, which individually execute attacks in separate goroutines.

[version 0.2](TODO ADD LINK TO V0.2):

Version 0.2 builds on v0.1, upgrading to a single-broadcaster, single-worker architecture. Using this design, attacks are requested from my local machine, which utilizes remote procedure calls to forward the work to an EC2 worker node. The worker node then manages the attack.

[version 0.3](TODO ADD LINK TO V0.3):

Version 0.3 continues to build off of previous iterations of the loris attacker, and utilizes a single-broadcaster, N-worker architecture. The broadcaster is an HTTP server, which serves attack requests, forwarding these requests to a specified number of worker nodes, which each execute attacks on the same target. The broadcaster node also serves termination requests, which are forwarded to worker nodes, which promptly terminate their ongoing attacks.

[version 1.0](TODO ADD LINK TO V1.0):

Version 1.0 is the final (for now) version of this slow loris attacker. It is a comprehensive system, which consists of the previously implemented single-broadcaster, N-worker architecture, but now also includes a well-defined client interface for users attempting to request an attack to be served.
