# DistHashTable
1. What system model are you assuming in your implementation. Write a full desciption.
I assume the system is using active replication with 2 nodes. When a new Put or Get request has been made the client, working also as a frontend, replicates the requests for the servers. Both of the serves have a HashTable and the system prevents them to have inconsisteny between them buy locking the critical section while updating or getting a value from the server 

2. What is the minimal number of nodes in your system to fullfill the requirements? Why?
The minimal number of nodes in the system is 3, 1 client and 2 servers. In active replication the frontend has been merged to the client so that the client node is distributing the messages to all servers. We need two servers to keep the system alive and responsive, although one server craches 

3. Explain how your system recovers from crash failure.
If one server craches, the remaining server will stay alive to act and respond to the incoming requests. When making a PUT request to the server, the crached server will result to a error which will be indicated to the user with a "request failed" message
If the client crashes, the other clients will stay alive and keep working as usual.

4. Explain how you achieve the Property 1 requirement.
The system will iterate through the HashTable and check if the given key exists in the HashTable. If the condition (key == GetRequest.GetKey()) is true then the value attached to it will be the val. The two replicas have their own HashTables and the system makes them synchronized by locking the server when updating it. If the server is updating it while doing a GET request, it will not return a value to the client until the requested servers are outputting the same values. I do this by comparing the values and returning the val once the vals from the two tables are in sync. 
5. Explain how you achieve the Property 2 requirement.
The system has a condition where it will return a value if the given key is in the HashTable. However, if the key doesn't exsist in the HashTable it will return 0
6. Explain how you achieve the Property 3 requirement.
Each request is replicated to each server. In order to keep the system consistent, the system compares the two values from the get method and does not return a value until the values are equal. If the values are not equal, it means that the system is still updating a put request.
7. Explain how you achieve the Property 4 - Liveness requirement.
The system will return false from a put request if a server has crached. Since the system will have the remaining server, this server will be returning true to the put request given that only one server has crached.
If the servers are busy they will queue the requests and take the next request from the queue. 
8. Explain how you achieve the Property 5 - Reliability requirement.
With active replication, if a node craches, the remaining replica(s) will be still working and it will provide responses to the client given that there has been 2 server in the beginning.
