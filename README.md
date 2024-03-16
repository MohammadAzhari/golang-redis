# TCP Durable Redis Server Clone

This project is a clone of a TCP durable Redis server implemented in Go. The application follows the architecture shown in the following diagram:

![Architecture Diagram](/screenshots/arc.png)

## RESP Package Test Coverage

The RESP package, responsible for handling the Redis Serialization Protocol (RESP), has full test coverage to ensure its reliability and correctness.

## Running the Application

To run the application, execute the following command in your terminal:

```
go run *.go
```

## Testing the Application

To test the application, you can use a Redis client. You can download a Redis client from [Redis Downloads](https://redis.io/download). After installing the client, you can connect to the server and test its functionality.

---

The project idea implemented after reading this blog https://www.build-redis-from-scratch.dev/en/introduction