# Zookeeper over Scalog
#### Building Fault Tolerant Applications without implementing fault tolerance

This project proposes a streamlined approach to instill fault tolerance in distributed applications by leveraging a fault-tolerant shared log as foundational infrastructure. This shared log serves as a resilient distributed logging service, capturing state changes and enabling the restoration of application states in the face of system failures. Demonstrating the utility of this abstraction, we implement resilient data structures in a Zookeeper service, where each node's internal structure benefits from the inherent fault tolerance of the distributed shared log. A consensus module acts as a cohesive link between the application and the shared log, ensuring a comprehensive fault-tolerant framework. Our project aims to showcase the simplicity of developing fault-tolerant applications using this foundational base. Through experiments comparing our Zookeeper implementation with a ZAB-consensus-based system, we affirm the viability and efficiency of our approach, validating the practicality of developing fault-tolerant applications with ease.

##### Architecture
<img width="922" alt="Arch" src="https://github.com/shahidikram0701/FaultTolerantDSOverSL/assets/27767537/80fb9573-ec8c-4feb-ad77-619e72f1a5e3">

As illustrated in the Figure, we have multiple applications accessing the Zookeeper. The Zookeeper is an ensemble of replicated services. The ensemble is present only for scaling the read operations and does not provide fault tolerance. Each of the Zookeeper service has a hierarchical namespace in the form of a trie. This trie maintains the internal state of the zookeeper data in a format that provides efficient lookups. This namespace closely mirrors the format in which the namespace is maintained in the original Zookeeper. We now have an additional shared log component that is backing up the Zookeeper. The shared log is transparent to the applications. The consensus module in each of the individual services provides a medium for the Zookeeper to communicate with the underlying shared log. We use scalog as our underlying shared log implementation. The scalog configuration in our experiments consists of 2 data shards, each in a primary-backup replica mode and 3 ordering nodes that internally run the raft protocol to order the data appended to the log.


##### Consensus module
<img width="755" alt="ConsensusModule" src="https://github.com/shahidikram0701/FaultTolerantDSOverSL/assets/27767537/919bdc86-10da-4adf-a11a-dc86cd5544c3">

The consensus module is a pipe into the shared log. This module is tasked with managing the interface between Zookeeper and the shared log. Its state encompasses three key elements: the local sequence number (LSN), representing the most recent operation from the log applied to the in-memory trie; the metadata component, housing specific auxiliary information vital for a log implementation; and the shared log client, furnishing APIs to interact with the shared log. Importantly, our design enables flexibility by allowing a seamless transition to a new underlying shared log. This adaptability is facilitated by instantiating a new shared log and seamlessly connecting the client to the consensus module. The metadata component further proves invaluable, offering a repository for auxiliary information crucial to supporting the state of the consensus module's interaction with the shared log.

##### Metadata management

Encountering a challenge with scalog as the underlying log implementation revealed a notable issue: clients writing to the log directly interacted with the data nodes. Consequently, when multiple Zookeeper nodes, functioning as clients of the shared log system, appended data to the shared log, the data shard that contains the data appended by one node, would not known by the other. This underscored the need for a mechanism to communicate this information among nodes, a process we term "metadata management."

Every Zookeeper node keeps track of a list detailing all shards containing data associated with the sequence numbers it appended to the log. In the background, these nodes actively communicate with other Zookeeper nodes to keep their tables updated. If, at any point, shard information for a particular sequence number cannot be retrieved, the node can perform a direct inquiry of the log to obtain the necessary details. Maintaining a mapping of sequence numbers to corresponding shards could lead to significant table growth, raising concerns about its scalability. To address scalability concerns with mapping entries, we adopt a proactive pruning approach. Entries are removed once operations tied to their sequence numbers are applied. However, to prevent premature pruning and ensure critical mappings are retained, we implement a threshold on the mapping table's size. With each entry consuming only 13 bytes, storing up to 200 entries is negligible in terms of network bandwidth, ensuring efficient metadata management.

#### Implementation

##### CreateZNode(path, data)
<img width="970" alt="CreateZNode" src="https://github.com/shahidikram0701/FaultTolerantDSOverSL/assets/27767537/3b51743e-f02a-41a7-949b-f71d5295e8ff">

The CreateZNode() API provides an interface for the clients to create an entry in the shared log which is later used by Zookeeper's hierarchical data structure. It takes 2 variables as input arguments: path and data. The path is used to indicate the node in the Zookeeper data structure, whereas the data consists of the operation performed (and not just the actual data). Whenever a client requests an operation, it gets logged into the underlying shared log as is, and the client is acknowledged with a success of a failure message accordingly. Based on the read policy (defined in the Reads sub-section), the entire log record (operation and data) gets processed and stored in a Zookeeper node.

##### GetZNode(path)
<img width="967" alt="GetZNode" src="https://github.com/shahidikram0701/FaultTolerantDSOverSL/assets/27767537/b8412ce0-b0fb-4fc0-82aa-568e431022a9">

The GetZNode() API serves as a client interface for data retrieval in our system. By taking a 'path' argument as input, representing the location where the data is stored, clients can request data from Zookeeper nodes. These nodes respond based on their internal state. If the specified path exists within the tree-like data structure of a Zookeeper node, it promptly reads and returns the associated data. Alternatively, if the path is not present, the node reads all unprocessed log records, sequentially applies the operations, and then returns the relevant data value. To uphold strong consistency guarantees, we augment the log with read operations. Integrating the Get operation into the log establishes a sequential order for operations. This sequencing enables us to determine the point until which all operations must be applied before providing the result to the client. Consequently, once the result is returned, the client is assured of consistently observing this value or a fresher one.

#### Reads
While we have established writes to the log are fast, reads can be slow as the Zookeeper data structure may not be potentially up-to-date when a read query arrives. This section explains the various read policies that we have in place to speed up the read operations.

##### Lazy Reads
The Lazy Reads policy states that writes to the Zookeeper Trie are deferred indefinitely until a read query arrives. This significantly reduces the overhead of synchronization of the Trie across nodes. This policy is preferred in a write-heavy workload scenario, however may perform substandard when there are a mixture of reads and writes. 

##### Asynchronous updates
The Asynchronous updates policy states that writes to the Zookeeper Trie are deferred for a short duration of time and then applied to the Trie in a parallel thread. This ensures the data structure to be up-to-date when the read arrives while still ensuring the write acknowledgements does not suffer. The defer duration variable needs to be tuned for efficiency. If the defer duration is too long then reads may arrive before the before the Trie has been updated, and if the duration is too short then the write performance suffers.

##### Optimizing asynchronous updates
To improve upon the asynchronous update policy model we employ batching and parallelism. When the asynchronous defer duration is triggered, we batch fetch all the pending serial numbers that are waiting on a update to be applied. Post that, a thread is a spawned for each of them in parallel to fetch the operations associated with them from the shared log. When all the threads join, the operations are ordered and applied to the Zookeeper trie.
This approach provides significant benefit over the vanilla asynchronous update model.

#### Failures
Due to the nature of our fault tolerance mechanism employed - shared log, recovery from failures is a straightforward task. When a node fails, no data is lost as it has already been committed to the scalog layer and other replica nodes can function without disruption. On a node recovery, or an addition of a new node; the scalog is queried and all the operations are replayed populating the Trie in-order.
