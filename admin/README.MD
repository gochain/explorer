# The admin command line tool

## How to launch on a server

This cli tool requires to have direct access to explorer's mongo database.
There are a couple of way how it could be launched but the following way is recommended (launch it on the same server where explorer is running)

```
docker exec -it web admin --mongo-url mongodb:27017 -testnet reload c 0xCB69Cc6c28131Cb4768913DD56B15C4026776dd3
```
