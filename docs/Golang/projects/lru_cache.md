
This plan is structured into **Objective**, **Architecture**, **Implementation Milestones**, and **Testing**.

---

# 💻 Project Plan: LRU Cache Implementation in Golang

## ⭐ Goal and Objectives

**Goal:** To design and implement a thread-safe Least Recently Used (LRU) cache data structure in Go that provides $O(1)$ average time complexity for both the `Get` (read) and `Put` (write) operations.

**Success Criteria:**
1.  The cache must correctly enforce the capacity limit, evicting the least recently used item when full.
2.  Read/Write operations (`Get`/`Put`) must operate in constant time $O(1)$.
3.  The structure must be thread-safe, preventing race conditions in concurrent use.

## 🧱 Phase I: Architecture and Data Structures Design (The Blueprint)

### A. Core Components Analysis

To achieve $O(1)$ complexity for both lookups and reordering/eviction, we need a combination of two standard computer science structures:

1.  **Hash Map (`map[KeyType] *Node`):** Used for fast key lookup ($O(1)$). It allows us to immediately access the associated node in our list given any key.
2.  **Doubly Linked List:** Used to maintain the order of usage (Most Recently Used $\rightarrow$ ... $\rightarrow$ Least Recently Used). Because it is *doubly* linked, we can move a node from its current position and re-insert it at the head in $O(1)$ time.

### B. Proposed Go Data Structures

We will define three primary structs:

**1. `Node` (The element in the list):**
*   `Key`: The type of the key used for the map lookup.
*   `Value`: The associated data value stored in the cache.
*   `Prev *Node`: Pointer to the previous node in the list.
*   `Next *Node`: Pointer to the next node in the list.

**2. `LRUCache` (The container/cache itself):**
*   `capacity int`: The maximum number of items the cache can hold.
*   `map map[KeyType]*Node`: Maps keys to their corresponding Node pointer for $O(1)$ access.
*   `head *Node`: Pointer to the head of the linked list (MRU - Most Recently Used).
*   `tail *Node`: Pointer to the tail of the linked list (LRU - Least Recently Used).
*   ***Concurrency Control:*** `mu sync.Mutex`: Required mutex to ensure thread safety.

## 🏗️ Phase II: Implementation Milestones (The Code Flow)

This phase breaks down the implementation into logical, testable steps.

### Milestone 1: Basic Structure and Initialization (Setup)
*   **Action:** Define the `Node` and `LRUCache` structs.
*   **Function:** Implement a constructor function (`NewLRUCache(capacity int) *LRUCache`).
*   **Check:** Initialize the map, set `capacity`, and ensure the list pointers are nil or correctly pointing to sentinel nodes (optional but robust).

### Milestone 2: List Management Helpers (The $O(1)$ Engine)
Before implementing core cache methods, we need helper functions that manage the linked list structure efficiently. This is the most critical step.
*   **`moveToHead(node *Node)`:** Given a node pointer, removes it from its current position in the list and re-links it to be directly after the `head`. (Complexity: $O(1)$)
*   **`popTail()`:** Removes the tail node (the LRU item) and returns its key/value. This handles eviction. (Complexity: $O(1)$)

### Milestone 3: Implementing `Put` (Write Operation)
The `Put` operation must handle three cases atomically:
1.  **Key Exists (Hit):** If the key is in the map, retrieve the node. Call `moveToHead()` on that node to mark it as recently used. Update its value if necessary.
2.  **New Key:** Create a new `Node`. Add the key/value pair to the map and insert the node at the head of the list (MRU).
3.  **Capacity Check & Eviction:** If, after adding the new item, the cache size exceeds `capacity`, call `popTail()` to remove the LRU node before returning.

### Milestone 4: Implementing `Get` (Read Operation)
The `Get` operation must also maintain usage order:
1.  **Lookup:** Check if the key exists in the internal map.
2.  **Miss:** If not found, return a zero value/error immediately.
3.  **Hit:** If found, retrieve the node pointer. Call `moveToHead()` on this node to promote it (marking it as recently used). Return the node's value.

### Milestone 5: Concurrency Safety (Polish)
*   **Action:** Wrap all public methods (`Get`, `Put`) with lock/unlock logic using `sync.Mutex`.
*   `func (c *LRUCache) Get(...) string { c.mu.Lock(); defer c.mu.Unlock(); ... }`

## 🧪 Phase III: Testing and Validation Strategy

Testing must cover functional correctness, edge cases, and concurrent behavior.

### A. Unit Test Cases (Functional Logic)
| Test Case | Description | Expected Outcome | Complexity Check |
| :--- | :--- | :--- | :--- |
| **Basic Put/Get** | Add item A, retrieve item A. | Successful retrieval of A's value. | $O(1)$ Get |
| **Capacity Limit** | Set capacity to 2. Insert A, B, C. | A is successfully evicted; only B and C remain. | Correct eviction (LRU) |
| **Hit Update Order** | Cache: [A (MRU), B, C (LRU)]. Get A. | New usage order: [A (MRU), B, C] $\rightarrow$ [A (MRU), B, C]. The list structure must be updated *without* changing the size. | $O(1)$ List Reordering |
| **Overwrite** | Cache: [B, C]. Put B (new value). | B's node is moved to MRU, its value is updated. Eviction check skipped if capacity allows. | $O(1)$ Update & Move |

### B. Stress/Concurrency Testing
*   **Test:** Run multiple goroutines concurrently calling `Put` and `Get` on the same cache instance (e.g., 50 concurrent reads, 20 concurrent writes).
*   **Validation:** Ensure that no panics occur due to race conditions, and the final state of the cache size and contents remains consistent with the operations performed.

## 🗓️ Deliverables Checklist

| Status | Component | Description | Notes |
| :---: | :--- | :--- | :--- |
| $\square$ | `Node` Struct | Defines key/value pairs and necessary pointers (`Prev`, `Next`). | Core data structure. |
| $\square$ | `LRUCache` Struct | Holds the map, head, tail, capacity, and mutex. | Container logic. |
| $\square$ | `NewLRUCache()` | Constructor function. | Initialization logic. |
| $\square$ | `moveToHead(node *Node)` | Helper method for list re-linking. | Essential $O(1)$ operation. |
| $\square$ | `Get(key KeyType) (ValueType, bool)` | Retrieves item and updates usage order. | Includes mutex lock. |
| $\square$ | `Put(key KeyType, value ValueType)` | Inserts/updates item, handles eviction when capacity is exceeded. | Includes mutex lock & size check. |
| $\square$ | Unit Tests (Go testing package) | Comprehensive tests covering all edge cases and concurrency. | Verification of correctness. |