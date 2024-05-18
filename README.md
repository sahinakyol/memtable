## A "MemTable" implementation for educational purposes.

- In-memory
- Data is key-value
- All inserts and updates in here
- Soft delete with tombstone
- MemTables have default size(64M). It swaps with new MemTable, when it gets full

| Type | Key   | Value | 
|------|-------|-------|
| ADD  | hello | world |

### Flush
A thread persists immutable MemTables to disk in background. You can give some intervals.

Immutable MemTable ----Flush----> Level_0_SST


### SST - Sorted String Tables

### Compaction


### Lookup steps
1- Search the active memtable.

2-Search immutable memtables.

3-Search all SST files on L0 starting from the most recently flushed.

4-For L1 and below, find a single SST file that may contain the key and search the file.

Searching an SST file involves:

1-(optional) Probe the bloom filter.

2-Search the index to find the block the key may belong to.

3-Read the block and try to find the key there.



