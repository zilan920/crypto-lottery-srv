# crypto lottery sever


- functions
  - generate random private key
  - store generated key on logs / aws s3
  - check address on ethereum if any balance

- todo
  - support multi chain
  - support telegram bot
  - customize goroutine number
  - support eth endpoints (etherscan api currently)
  - important: implement distributed system to calculated on different machine, like zookeeper ?

- Thinking...
  - generate bloom filter to check balance instead of using endpoint
  - what is best way to store rainbow table of blockchain ?
  - essentially private key is just a big number, how would those been using private key distributed over all key range ? this could help us design better `random` algorithm 
  