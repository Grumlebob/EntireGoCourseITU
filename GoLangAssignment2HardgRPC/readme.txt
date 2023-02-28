Help from: https://github.com/NaddiNadja/grpc101

To launch:
open terminal in server folder:
go run .

Open another terminal in client folder:
go run .

a) What are packets in your implementation? 
a2)What data structure do you use to transmit data and meta-data?
Our packets consist of a message to illustrate the 3 handshakes, but primarily 2 primary integer types, being the Acknowledgement number and the Sequence number. We do not send actual data.


b) Does your implementation use threads or processes? 
b2) Why is it not realistic to use threads?
A thread won't spin across networks so to say. But if our implementation uses processes, we can containerize them, and horizontally scale them, by using the same image multiple times. If we were to use threads, we would have to vertical scale, by adding more CPU and memory, which is so to say not scalable in extreme terms, due to limits of techonology, which is why it isn't realistic to use threads.


c) How do you handle message re-ordering? & d) How do you handle message loss?
We don't, as we focused all our time on understanding TCP, especially the three way handshake. If we were to make an implementation, we would ensure that the server actually checks the Sequence number, compared to the expected sequence number. If the server lost some data, it would print the last Acknowledgement number it had before it went wrong, and the client would have to re-send all the data, from that acknowledgement number.


e) Why is the 3-way handshake important?
Opposed to UDP, TCP is connection-oriented, meaning we want to be sure that we are actually having a conversation with someone, who is answering. So by having the handshakes, we negotiate the starting sequence numbers, and make sure the connection is established before sending data.
