# How to launch

There are 2 ways to launch this program.

# The easy way:

If you are on windows with powershell installed, simply right click the file Start2Peers.ps1 and select "run with powershell"

# The manual way:

1. Open 2 different terminals.
2. Write in the first terminal: go run . 0
3. Write in the second terminal: go run . 1
4. Demo data will automatically be inserted and printed
5. If you want to add, write 'add key value' to put a key-value pair in the dictionary
6. If you want to read, write 'read key' to get the value of a key
7. If you want to print the distributed dictionary, write 'printdictionary'
8. If you want to print the requests handled map, write 'printrequests'
9. If you want to force crash the client, enter 'crash' or use ctrl+c. Clients will automatically figure our new leader
