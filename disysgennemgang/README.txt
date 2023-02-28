Achitecture:
two servers: a main and a replica, either can crash and the system live on.
using active replication.

Get:
if client wants to get they get from either server since get does not 
mutate the hashtable no syncronisation is needed. Responce is an int32.

Put:
if client wants to put they try to get a lock on main server.
then they try for a lock on the replica server.
when they have both, or one and a failure, they send the put request those it 
got a lock on.
when a boolean response has been recieved from all alive servers the locks are 
relinquished.

since we need a lock on all alive servers to begin and we always ask for lock 
on main first there will be no deadlocks and the hashtables wil be kept in sync

