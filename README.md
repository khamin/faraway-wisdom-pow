# About

Design and implement "Word of Wisdom" TCP server:

- TCP server should be protected from DDoS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), the challenge-response protocol should be used.
- The choice of the POW algorithm should be explained.
- After Proof Of Work verification, server should send one of the quotes from "word of wisdom" book or any other collection of the quotes.
- Docker file should be provided both for the server and for the client that solves the POW challenge.

# PoW Function Choice

Note: network-bound functions like [GTP](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol) are better choice but excluded due to infrastructural reasons of this demo.

The computing power of clients can reach a difference of 30 times in single-threaded mode, and 15 times in multi-threaded mode.
Also, the resolution of some tasks can be significantly accelerated using GPU or ASIC.
This makes CPU-bound functions inefficient in terms of client equality.

[Equihash](https://en.wikipedia.org/wiki/Equihash) [looks fair](whitepaper.pdf) and fits our needs.

# Run

```sh
docker-compose up --build
```
