# The Real-Life Guide to Data Structures & Algorithms

> Explained like you're five, coded like you're a Go engineer — and, for
> every topic, coded like a C engineer too.

## Why does this exist?

Imagine your toys are scattered all over your room. **Data Structures** are the different types of toy boxes you use to organize them — a bin for blocks, a shelf for books, a hook for backpacks. **Algorithms** are the steps you follow to *do something* with the toys — like "how do I find my favorite car the fastest?" or "how do I line up my stuffed animals from smallest to biggest?"

Every data structure is just: **a box + rules for how you're allowed to put things in and take things out.**
Every algorithm is just: **a recipe — a list of steps that always gets you the right answer.**

That's it. Everything below is that idea, over and over, in fancier clothes.

We'll go **Beginner → Intermediate → Advanced**. Every topic has:
- 🧸 **ELI5** — the toy-box explanation
- 🌍 **Real life** — where you *already* use this without knowing it
- 📝 **Pseudocode** — the recipe in plain steps
- 🐹 **Go code** — the real thing, with the standard library doing the heavy lifting
- 🇨 **C code** — the *same* structure with no standard library to hide behind: you write the `malloc`, the pointer chasing, and the array-index math yourself
- ⏱️ **Complexity** — how the recipe scales when you have way more toys

**Why both languages?** Go's `map`, `container/heap`, and garbage collector do
real work you don't see. Writing the same data structure in C forces every
piece of that work into view — which is exactly why this guide's C examples
are paired with **[`wiki/c-lang/real-life-c-guide.md`](../c-lang/real-life-c-guide.md)**,
a full C field guide: its Part III ("Data structures from scratch," §17–19)
builds the linked lists, hash tables, and trees used below in more depth,
and its Part VIII covers the memory-safety discipline (leak patterns,
ownership rules) every one of this guide's `malloc`/`free` pairs relies on.
Each 🇨 block below links to the exact section that goes deeper.

---

# 📏 Big-O Notation — "How Slow Does This Get?"

Before we meet any data structure, we need a way to talk about **how a recipe behaves when the pile of toys gets huge**. That's all Big-O is.

🧸 **ELI5**: Imagine you have 3 toys — checking every single one to find your favorite is quick, maybe 3 seconds. Now imagine you have 3 *million* toys. Big-O is just the answer to: **"if I 1,000x the toys, does the search take 1,000x longer, the same time, or something in between?"** It's not a stopwatch — it's a *shape of growth*. We use the letter `n` to mean "however many things you have."

🌍 **Real life**: Big-O is why a phone book with a million names is only slightly slower to search than one with a thousand names (binary search — `O(log n)`), but why checking every student against every other student in a class for "who has the same birthday" gets painful fast as the class grows (`O(n²)`).

Here's the full scale, from best to worst, each with something you already do in real life:

| Big-O | Name | ELI5 | Real-world example |
|---|---|---|---|
| `O(1)` | Constant | Doesn't matter how many toys exist — this is always instant. | Swiping a keycard at a door. Checking if a light switch is on. Hash map lookup. |
| `O(log n)` | Logarithmic | Cuts the problem in half every step. | Guessing a number 1-100 in ≤7 tries. Looking up a word in a paper dictionary. Binary search. |
| `O(n)` | Linear | You touch every item exactly once. | Reading every page of a book once. Counting every kid in a classroom. Linear search. |
| `O(n log n)` | Linearithmic | Split the pile up, then do a little work on each piece. | Efficiently sorting a deck of cards (merge sort / quick sort). |
| `O(n²)` | Quadratic | Compare every item to every other item. | Checking every student against every other student for matching birthdays. Bubble sort. |
| `O(2ⁿ)` | Exponential | Every new item *doubles* your work — tries every possible combination. | Trying every possible pizza topping combination (add 1 topping = the choices double). Naive recursive Fibonacci. |
| `O(n!)` | Factorial | Tries every possible *ordering* of everything. | Trying every possible seating arrangement for `n` wedding guests. Brute-force Traveling Salesman. |

**Why `n` matters at scale** — say you have 1,000,000 items:

```
O(1)        ->  1 step            (instant, always)
O(log n)    ->  ~20 steps         (binary search on a million records)
O(n)        ->  1,000,000 steps   (scan everything once)
O(n log n)  ->  ~20,000,000 steps (a good sort)
O(n²)       ->  1,000,000,000,000 steps  (bubble sort on a million items — do NOT try this)
```

That jump from `O(n log n)` to `O(n²)` is *the* reason engineers care so much about picking the right data structure — it's the difference between a page loading instantly and a server timing out.

📝 **Pseudocode** (spotting Big-O by counting loops):
```
function example(list):        // n = length(list)
    for item in list:          // one loop over n items -> O(n)
        print(item)

function example2(list):
    for a in list:              // loop 1: n
        for b in list:          // loop 2 nested inside: n
            print(a, b)          // nested loops -> O(n * n) = O(n²)
```

🐹 **Go** (same shapes, real code):
```go
package main

import "fmt"

// O(1) — constant time, no matter how big `nums` is
func firstItem(nums []int) int {
	return nums[0]
}

// O(n) — linear, one pass over every item
func sum(nums []int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// O(n²) — quadratic, a loop inside a loop
func hasDuplicatePair(nums []int) bool {
	for i := range nums {
		for j := range nums {
			if i != j && nums[i] == nums[j] {
				return true
			}
		}
	}
	return false
}

func main() {
	nums := []int{4, 8, 15, 16, 23, 42}
	fmt.Println(firstItem(nums))        // O(1)
	fmt.Println(sum(nums))              // O(n)
	fmt.Println(hasDuplicatePair(nums)) // O(n²)
}
```

⏱️ **The rule of thumb**: every topic below lists its own `⏱️ Complexity`. When you're choosing a data structure or algorithm for a real project, ask **"how big can `n` get?"** — `O(n²)` is perfectly fine for 50 items and a real problem for 5,000,000.

---

# 🟢 BEGINNER

Big-O gave you the ruler for measuring "how slow does this get." This tier
gives you the actual toy boxes — the containers and the two slowest ways to
search and sort — that every fancier structure later in this guide is built
on top of. None of these seven topics are hard; each one just solves exactly
one "how do I organize/find/undo this" problem.

## 1. Arrays

🧸 **ELI5**: An array is an egg carton. It has a fixed number of slots, all lined up, each with a number on it (0, 1, 2, 3...). You can look inside slot #4 instantly because you know exactly where it is — you don't have to open every slot first.

🌍 **Real life**: A parking garage with numbered spots. A row of mailboxes in an apartment building. Seats in a movie theater (Row C, Seat 12).

📝 **Pseudocode**:
```
array = [10, 20, 30, 40]
read array[2]        // instantly returns 30
```

🐹 **Go**:
```go
package main

import "fmt"

func main() {
	parkingSpots := [5]string{"Car", "", "Bike", "", "Truck"}
	fmt.Println(parkingSpots[2]) // "Bike" — instant lookup
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §6](../c-lang/real-life-c-guide.md#6-arrays-pointers-and-the-arithmetic-that-explains-everything-else)):
```c
#include <stdio.h>

int main(void) {
    const char *parkingSpots[5] = {"Car", "", "Bike", "", "Truck"};
    printf("%s\n", parkingSpots[2]); // "Bike" — instant lookup, same address-math trick
    return 0;
}
```

⏱️ **Complexity**: Read by index = `O(1)` (instant). Insert/delete in the middle = `O(n)` (you have to shift everyone over, like asking every car to move up one spot).

---

## 2. Strings

🧸 **ELI5**: A string is just an array of letter-beads on a necklace string. "CAT" is 3 beads in a row: C, A, T.

🌍 **Real life**: A text message. A license plate. A word in your favorite book.

📝 **Pseudocode**:
```
word = "cat"
for each letter in word:
    print(letter)
```

🐹 **Go**:
```go
package main

import "fmt"

func main() {
	word := "cat"
	for i, letter := range word {
		fmt.Printf("bead %d is %c\n", i, letter)
	}
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §7](../c-lang/real-life-c-guide.md#7-strings-in-c-buffers-ownership-and-why-gets-was-removed)):
```c
#include <stdio.h>

int main(void) {
    const char *word = "cat";       // in C, a "string" IS just an array of letter-beads + '\0'
    for (int i = 0; word[i] != '\0'; i++) {
        printf("bead %d is %c\n", i, word[i]);
    }
    return 0;
}
```

⏱️ **Complexity**: Reading a character by position = `O(1)`. Scanning the whole word = `O(n)`.

---

## 3. Linked Lists

🧸 **ELI5**: Instead of an egg carton, imagine a **treasure hunt**. Each clue (node) tells you two things: the treasure at this spot, and where to find the *next* clue. You can't jump straight to clue #5 — you must follow clue #1 → #2 → #3 → #4 → #5.

🌍 **Real life**: A conga line / train cars — each car is only connected to the one behind and in front. A scavenger hunt. Your browser history's "back" chain (kind of).

📝 **Pseudocode**:
```
node = { value, next }
head -> node(1) -> node(2) -> node(3) -> null

function printAll(head):
    current = head
    while current != null:
        print(current.value)
        current = current.next
```

🐹 **Go**:
```go
package main

import "fmt"

type Node struct {
	Value int
	Next  *Node
}

func printAll(head *Node) {
	for cur := head; cur != nil; cur = cur.Next {
		fmt.Println(cur.Value)
	}
}

func main() {
	third := &Node{Value: 3}
	second := &Node{Value: 2, Next: third}
	first := &Node{Value: 1, Next: second}
	printAll(first) // 1, 2, 3
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §18.1](../c-lang/real-life-c-guide.md#181-singly-linked-list)) — this is the language linked lists were *born* in; a `Node` here is a `malloc`'d block plus a pointer, nothing more:
```c
#include <stdio.h>
#include <stdlib.h>

typedef struct Node {
    int value;
    struct Node *next;
} Node;

Node *node_new(int value, Node *next) {
    Node *n = malloc(sizeof(Node));
    n->value = value;
    n->next = next;
    return n;
}

void print_all(Node *head) {
    for (Node *cur = head; cur != NULL; cur = cur->next) {
        printf("%d\n", cur->value);
    }
}

void list_free(Node *head) {   // unlike Go, nothing frees this for you — Section 10.3
    while (head) { Node *next = head->next; free(head); head = next; }
}

int main(void) {
    Node *third = node_new(3, NULL);
    Node *second = node_new(2, third);
    Node *first = node_new(1, second);
    print_all(first); // 1, 2, 3
    list_free(first);
    return 0;
}
```

⏱️ **Complexity**: Access by position = `O(n)` (follow the clues one by one). Insert at the front = `O(1)` (just point to a new first clue) — much faster than an array's `O(n)` shift.

---

## 4. Stacks

🧸 **ELI5**: A stack is a stack of pancakes. You can only add a pancake to the **top**, and you can only eat from the **top**. You can never grab the bottom pancake without removing all the ones above it first. **Last one in, first one out** (LIFO).

🌍 **Real life**: The "Undo" button (Ctrl+Z) — it undoes your *most recent* action first. Your browser's "Back" button. A stack of trays in a cafeteria.

📝 **Pseudocode**:
```
stack = []
push(stack, "plate1")
push(stack, "plate2")
pop(stack)   // removes "plate2" (the last one added)
```

🐹 **Go**:
```go
package main

import "fmt"

type Stack struct {
	items []string
}

func (s *Stack) Push(item string) {
	s.items = append(s.items, item)
}

func (s *Stack) Pop() (string, bool) {
	if len(s.items) == 0 {
		return "", false
	}
	last := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return last, true
}

func main() {
	undoHistory := &Stack{}
	undoHistory.Push("typed 'hello'")
	undoHistory.Push("typed 'world'")
	action, _ := undoHistory.Pop()
	fmt.Println("undo:", action) // undo: typed 'world'
}
```

🇨 **C** (a fixed-capacity array-backed stack — no dynamic growth needed for a stack of pancakes):
```c
#include <stdio.h>
#include <string.h>

#define STACK_CAP 32

typedef struct {
    char items[STACK_CAP][64];
    int  top; // -1 means empty
} Stack;

void stack_push(Stack *s, const char *item) {
    if (s->top + 1 >= STACK_CAP) return;
    strncpy(s->items[++s->top], item, sizeof s->items[0] - 1);
}

int stack_pop(Stack *s, char *out, size_t out_cap) {
    if (s->top < 0) return 0;
    strncpy(out, s->items[s->top--], out_cap - 1);
    return 1;
}

int main(void) {
    Stack undoHistory = { .top = -1 };
    stack_push(&undoHistory, "typed 'hello'");
    stack_push(&undoHistory, "typed 'world'");

    char action[64];
    if (stack_pop(&undoHistory, action, sizeof action)) {
        printf("undo: %s\n", action); // undo: typed 'world'
    }
    return 0;
}
```

⏱️ **Complexity**: Push/Pop = `O(1)` — always fast, no matter how tall the stack.

---

## 5. Queues

🧸 **ELI5**: A queue is the line at an ice cream truck. Whoever got in line **first** gets their ice cream **first**. No cutting! You join at the back, you leave from the front. **First one in, first one out** (FIFO).

🌍 **Real life**: A printer queue (the first document you sent prints first). People waiting for a rollercoaster. Customer support tickets.

📝 **Pseudocode**:
```
queue = []
enqueue(queue, "Alice")
enqueue(queue, "Bob")
dequeue(queue)   // removes "Alice" (the first one who joined)
```

🐹 **Go**:
```go
package main

import "fmt"

type Queue struct {
	items []string
}

func (q *Queue) Enqueue(item string) {
	q.items = append(q.items, item)
}

func (q *Queue) Dequeue() (string, bool) {
	if len(q.items) == 0 {
		return "", false
	}
	first := q.items[0]
	q.items = q.items[1:]
	return first, true
}

func main() {
	iceCreamLine := &Queue{}
	iceCreamLine.Enqueue("Alice")
	iceCreamLine.Enqueue("Bob")
	served, _ := iceCreamLine.Dequeue()
	fmt.Println("served:", served) // served: Alice
}
```

🇨 **C** (a circular buffer — `head`/`tail` indices wrap around with `%`, so unlike Go's slice trick, nothing ever gets shifted in memory):
```c
#include <stdio.h>
#include <string.h>

#define QUEUE_CAP 32

typedef struct {
    char items[QUEUE_CAP][64];
    int  head, tail, count;
} Queue;

void queue_enqueue(Queue *q, const char *item) {
    if (q->count >= QUEUE_CAP) return;
    strncpy(q->items[q->tail], item, sizeof q->items[0] - 1);
    q->tail = (q->tail + 1) % QUEUE_CAP;
    q->count++;
}

int queue_dequeue(Queue *q, char *out, size_t out_cap) {
    if (q->count == 0) return 0;
    strncpy(out, q->items[q->head], out_cap - 1);
    q->head = (q->head + 1) % QUEUE_CAP;
    q->count--;
    return 1;
}

int main(void) {
    Queue iceCreamLine = {0};
    queue_enqueue(&iceCreamLine, "Alice");
    queue_enqueue(&iceCreamLine, "Bob");

    char served[64];
    if (queue_dequeue(&iceCreamLine, served, sizeof served)) {
        printf("served: %s\n", served); // served: Alice
    }
    return 0;
}
```

⏱️ **Complexity**: Enqueue/Dequeue = `O(1)` (with a proper ring buffer/linked list; Go's slice trick above is `O(1)` amortized but not perfectly optimal).

---

## 6. Linear Search

🧸 **ELI5**: You lost your favorite toy car in a toy box full of 50 toys. You pick up toy #1, check if it's the car, put it down, pick up toy #2, check... You keep going until you find it. No shortcuts.

🌍 **Real life**: Scanning every song in an unsorted playlist looking for one title. Checking every locker in a hallway for your friend's name because lockers aren't sorted by name.

📝 **Pseudocode**:
```
function linearSearch(list, target):
    for i from 0 to length(list) - 1:
        if list[i] == target:
            return i
    return -1
```

🐹 **Go**:
```go
package main

import "fmt"

func linearSearch(toys []string, target string) int {
	for i, toy := range toys {
		if toy == target {
			return i
		}
	}
	return -1
}

func main() {
	toys := []string{"ball", "robot", "car", "doll"}
	fmt.Println(linearSearch(toys, "car")) // 2
}
```

🇨 **C**:
```c
#include <stdio.h>
#include <string.h>

int linear_search(const char *toys[], int n, const char *target) {
    for (int i = 0; i < n; i++) {
        if (strcmp(toys[i], target) == 0) return i;
    }
    return -1;
}

int main(void) {
    const char *toys[] = {"ball", "robot", "car", "doll"};
    printf("%d\n", linear_search(toys, 4, "car")); // 2
    return 0;
}
```

⏱️ **Complexity**: `O(n)` — worst case, you check every single toy.

---

## 7. Bubble Sort

🧸 **ELI5**: Line up your friends by height. You compare the first two — if the shorter one is standing *behind* the taller one, they swap. You keep walking down the line doing this, and each pass, the tallest kid "bubbles up" to the end. Repeat until nobody needs to swap anymore.

🌍 **Real life**: Manually sorting playing cards by comparing two adjacent cards at a time and swapping. Simple, but slow for a big deck.

📝 **Pseudocode**:
```
function bubbleSort(list):
    for i from 0 to length(list) - 1:
        for j from 0 to length(list) - i - 2:
            if list[j] > list[j+1]:
                swap(list[j], list[j+1])
```

🐹 **Go**:
```go
package main

import "fmt"

func bubbleSort(heights []int) {
	n := len(heights)
	for i := 0; i < n; i++ {
		for j := 0; j < n-i-1; j++ {
			if heights[j] > heights[j+1] {
				heights[j], heights[j+1] = heights[j+1], heights[j]
			}
		}
	}
}

func main() {
	kids := []int{140, 120, 160, 130}
	bubbleSort(kids)
	fmt.Println(kids) // [120 130 140 160]
}
```

🇨 **C**:
```c
#include <stdio.h>

void bubble_sort(int *heights, int n) {
    for (int i = 0; i < n; i++) {
        for (int j = 0; j < n - i - 1; j++) {
            if (heights[j] > heights[j + 1]) {
                int tmp = heights[j];
                heights[j] = heights[j + 1];
                heights[j + 1] = tmp;
            }
        }
    }
}

int main(void) {
    int kids[] = {140, 120, 160, 130};
    bubble_sort(kids, 4);
    for (int i = 0; i < 4; i++) printf("%d ", kids[i]); // 120 130 140 160
    printf("\n");
    return 0;
}
```

⏱️ **Complexity**: `O(n²)` — fine for a handful of kids, terrible for a stadium of them.

---

# 🟡 INTERMEDIATE

The beginner tier gave you raw containers and the two slowest ways to search
(check every slot) and sort (compare neighbors, over and over). This tier is
where those get *fast* — binary search, hash maps with near-instant lookup —
and introduces the two shapes, trees and graphs, that most real systems
(file systems, org charts, social networks, road maps) are actually built
from underneath.

## 8. Binary Search

🧸 **ELI5**: You're playing a guessing game: "I'm thinking of a number between 1 and 100." Instead of guessing 1, 2, 3, 4... you guess **50** (the middle). Too high? Now you know it's between 1-49, so you guess **25**. Too low? Now you know it's between 26-49, so you guess... You cut the possibilities in **half** every single time.

🌍 **Real life**: Looking up a word in a paper dictionary — you don't start at "A," you flip to the middle, see you're past your word, flip back halfway, and narrow in. This *only* works if the list is already **sorted**.

📝 **Pseudocode**:
```
function binarySearch(sortedList, target):
    low = 0
    high = length(sortedList) - 1
    while low <= high:
        mid = (low + high) / 2
        if sortedList[mid] == target:
            return mid
        else if sortedList[mid] < target:
            low = mid + 1
        else:
            high = mid - 1
    return -1
```

🐹 **Go**:
```go
package main

import "fmt"

func binarySearch(sorted []int, target int) int {
	low, high := 0, len(sorted)-1
	for low <= high {
		mid := (low + high) / 2
		switch {
		case sorted[mid] == target:
			return mid
		case sorted[mid] < target:
			low = mid + 1
		default:
			high = mid - 1
		}
	}
	return -1
}

func main() {
	pages := []int{2, 8, 15, 23, 42, 71, 99}
	fmt.Println(binarySearch(pages, 42)) // 4
}
```

🇨 **C**:
```c
#include <stdio.h>

int binary_search(const int *sorted, int n, int target) {
    int low = 0, high = n - 1;
    while (low <= high) {
        int mid = low + (high - low) / 2;   // avoids (low+high) overflowing for huge arrays
        if (sorted[mid] == target) return mid;
        if (sorted[mid] < target) low = mid + 1;
        else high = mid - 1;
    }
    return -1;
}

int main(void) {
    int pages[] = {2, 8, 15, 23, 42, 71, 99};
    printf("%d\n", binary_search(pages, 7, 42)); // 4
    return 0;
}
```

⏱️ **Complexity**: `O(log n)` — doubling the toy box size only adds *one more* guess. 1,000,000 items? Just ~20 guesses.

---

## 9. Recursion

🧸 **ELI5**: Recursion is a set of Russian nesting dolls (matryoshka). To find the smallest doll, you open a doll, and inside is... a smaller version of the *exact same problem*: "open a doll and look inside." You keep doing the same tiny step until you hit the smallest doll (the **base case**) that doesn't open anymore — then you're done.

🌍 **Real life**: Standing between two mirrors facing each other — the reflection repeats itself smaller and smaller. Folders inside folders inside folders on your computer.

📝 **Pseudocode**:
```
function factorial(n):
    if n == 0:               // base case — the smallest doll
        return 1
    return n * factorial(n - 1)   // same problem, smaller size
```

🐹 **Go**:
```go
package main

import "fmt"

func factorial(n int) int {
	if n == 0 { // base case: stop opening dolls
		return 1
	}
	return n * factorial(n-1) // same task, one size smaller
}

func main() {
	fmt.Println(factorial(5)) // 120
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §5.1](../c-lang/real-life-c-guide.md#51-what-actually-happens-on-a-function-call)) — this is the guide where "the call stack" stops being an abstraction: every recursive call really is a new stack frame, and you can watch it overflow:
```c
#include <stdio.h>

long factorial(int n) {
    if (n == 0) return 1;          // base case: stop opening dolls
    return n * factorial(n - 1);   // same task, one size smaller
}

int main(void) {
    printf("%ld\n", factorial(5)); // 120
    return 0;
}
```

⏱️ **Complexity**: Depends on the problem — `factorial` is `O(n)` calls deep. Watch out: every call uses memory (the "call stack"), so too much recursion without a base case = **stack overflow** (imagine infinite nesting dolls that never end — you run out of room).

---

## 10. Hash Maps (a.k.a. Dictionaries / Hash Tables)

🧸 **ELI5**: A hash map is a **coat check** at a party. You hand over your coat, they give you ticket #42. Later, you show ticket #42, and they go **directly** to hook #42 and grab your coat — no searching through every coat in the building.

🌍 **Real life**: Your phone's contacts app — you type a name, and it instantly finds the number (it's not scanning every contact one by one). A locker with your name mapped to locker #17.

📝 **Pseudocode**:
```
map = {}
map["Alice"] = "555-1234"
map["Bob"]   = "555-5678"
print(map["Alice"])   // instantly -> "555-1234"
```

🐹 **Go**:
```go
package main

import "fmt"

func main() {
	contacts := make(map[string]string)
	contacts["Alice"] = "555-1234"
	contacts["Bob"] = "555-5678"

	number, found := contacts["Alice"]
	if found {
		fmt.Println(number) // 555-1234
	}
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §18.2](../c-lang/real-life-c-guide.md#182-a-real-hash-table-separate-chaining-djb2-hash)) — C has no built-in map, so this coat-check counter is the *actual mechanism* Go's `map` and Python's `dict` hide from you: a hash function plus separate-chaining buckets:
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define BUCKETS 16

typedef struct HTEntry {
    char *key, *value;
    struct HTEntry *next;
} HTEntry;

static HTEntry *buckets[BUCKETS];

static unsigned long djb2(const char *s) {   // Section 18.2's hash function
    unsigned long hash = 5381;
    int c;
    while ((c = (unsigned char)*s++)) hash = ((hash << 5) + hash) + (unsigned long)c;
    return hash;
}

void ht_set(const char *key, const char *value) {
    unsigned long idx = djb2(key) % BUCKETS;
    for (HTEntry *e = buckets[idx]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) { free(e->value); e->value = strdup(value); return; }
    }
    HTEntry *e = malloc(sizeof(HTEntry));
    e->key = strdup(key);
    e->value = strdup(value);
    e->next = buckets[idx];   // insert at the head of this ticket number's coat rack
    buckets[idx] = e;
}

const char *ht_get(const char *key) {
    unsigned long idx = djb2(key) % BUCKETS;
    for (HTEntry *e = buckets[idx]; e; e = e->next) {
        if (strcmp(e->key, key) == 0) return e->value;
    }
    return NULL;
}

int main(void) {
    ht_set("Alice", "555-1234");
    ht_set("Bob", "555-5678");

    const char *number = ht_get("Alice");
    if (number) printf("%s\n", number); // 555-1234
    return 0;
}
```
*(compile with `-std=gnu17`, not strict `-std=c17` — `strdup` is POSIX, hidden under strict ANSI on Linux/glibc; see the C guide's §2.2 note.)*

⏱️ **Complexity**: Insert/lookup/delete = `O(1)` average — the magic of "hashing" the key straight to its slot, like ticket #42 pointing straight to hook #42.

---

## 11. Trees

🧸 **ELI5**: A tree is your **family tree**. You (the root) have parents, they have parents (grandparents), and so on. Or flip it: you have children, they might have children. Everyone is connected, but there's exactly **one path** from you to any relative — no loops.

🌍 **Real life**: Folders and subfolders on your computer (`Documents > Photos > Vacation > Beach.jpg`). A company's org chart (CEO → VPs → Managers → Employees). The `<div>` structure of a webpage (HTML DOM).

📝 **Pseudocode**:
```
node = { value, children[] }

function printTree(node, depth=0):
    print(indent(depth) + node.value)
    for child in node.children:
        printTree(child, depth + 1)
```

🐹 **Go**:
```go
package main

import (
	"fmt"
	"strings"
)

type TreeNode struct {
	Name     string
	Children []*TreeNode
}

func printTree(node *TreeNode, depth int) {
	fmt.Println(strings.Repeat("  ", depth) + node.Name)
	for _, child := range node.Children {
		printTree(child, depth+1)
	}
}

func main() {
	grandpa := &TreeNode{Name: "Grandpa"}
	dad := &TreeNode{Name: "Dad"}
	uncle := &TreeNode{Name: "Uncle"}
	me := &TreeNode{Name: "Me"}

	grandpa.Children = []*TreeNode{dad, uncle}
	dad.Children = []*TreeNode{me}

	printTree(grandpa, 0)
	// Grandpa
	//   Dad
	//     Me
	//   Uncle
}
```

🇨 **C** (children stored as a `realloc`-grown array of pointers — Section 16.1's dynamic-array trick, applied per-node):
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

typedef struct TreeNode {
    char name[32];
    struct TreeNode **children;
    int   num_children;
} TreeNode;

TreeNode *tree_new(const char *name) {
    TreeNode *n = calloc(1, sizeof(TreeNode));
    strncpy(n->name, name, sizeof n->name - 1);
    return n;
}

void tree_add_child(TreeNode *parent, TreeNode *child) {
    parent->children = realloc(parent->children,
                                (size_t)(parent->num_children + 1) * sizeof(TreeNode *));
    parent->children[parent->num_children++] = child;
}

void print_tree(TreeNode *node, int depth) {
    for (int i = 0; i < depth; i++) printf("  ");
    printf("%s\n", node->name);
    for (int i = 0; i < node->num_children; i++) print_tree(node->children[i], depth + 1);
}

void tree_free(TreeNode *node) {
    for (int i = 0; i < node->num_children; i++) tree_free(node->children[i]);
    free(node->children);
    free(node);
}

int main(void) {
    TreeNode *grandpa = tree_new("Grandpa");
    TreeNode *dad = tree_new("Dad");
    TreeNode *uncle = tree_new("Uncle");
    TreeNode *me = tree_new("Me");

    tree_add_child(grandpa, dad);
    tree_add_child(grandpa, uncle);
    tree_add_child(dad, me);

    print_tree(grandpa, 0);
    // Grandpa / Dad / Me / Uncle
    tree_free(grandpa);
    return 0;
}
```

⏱️ **Complexity**: Traversing every node = `O(n)`. Height of a balanced tree with `n` nodes = `O(log n)`.

---

## 12. Binary Search Trees (BST)

🧸 **ELI5**: Imagine organizing your library books on a special shelf: for every book, **all books with a smaller title go to its left**, and **all books with a bigger title go to its right**. Now finding any book is like the guessing game (binary search) — go left or right at every shelf.

🌍 **Real life**: Auto-complete style organized dictionaries. Any "sorted lookup" structure like a phone book organized as a tree instead of a flat list.

📝 **Pseudocode**:
```
function insert(node, value):
    if node == null: return newNode(value)
    if value < node.value: node.left = insert(node.left, value)
    else: node.right = insert(node.right, value)
    return node

function search(node, value):
    if node == null: return false
    if node.value == value: return true
    if value < node.value: return search(node.left, value)
    return search(node.right, value)
```

🐹 **Go**:
```go
package main

import "fmt"

type BSTNode struct {
	Value       int
	Left, Right *BSTNode
}

func insert(node *BSTNode, value int) *BSTNode {
	if node == nil {
		return &BSTNode{Value: value}
	}
	if value < node.Value {
		node.Left = insert(node.Left, value)
	} else {
		node.Right = insert(node.Right, value)
	}
	return node
}

func search(node *BSTNode, value int) bool {
	if node == nil {
		return false
	}
	if node.Value == value {
		return true
	}
	if value < node.Value {
		return search(node.Left, value)
	}
	return search(node.Right, value)
}

func main() {
	var root *BSTNode
	for _, v := range []int{50, 30, 70, 20, 40} {
		root = insert(root, v)
	}
	fmt.Println(search(root, 40)) // true
	fmt.Println(search(root, 99)) // false
}
```

🇨 **C** ([deep dive: real-life-c-guide.md §18.3](../c-lang/real-life-c-guide.md#183-a-binary-search-tree-for-completeness)):
```c
#include <stdio.h>
#include <stdlib.h>

typedef struct BSTNode {
    int value;
    struct BSTNode *left, *right;
} BSTNode;

BSTNode *bst_insert(BSTNode *node, int value) {
    if (!node) {
        BSTNode *n = malloc(sizeof(BSTNode));
        n->value = value; n->left = n->right = NULL;
        return n;
    }
    if (value < node->value) node->left = bst_insert(node->left, value);
    else                     node->right = bst_insert(node->right, value);
    return node;
}

int bst_search(BSTNode *node, int value) {
    if (!node) return 0;
    if (node->value == value) return 1;
    return value < node->value ? bst_search(node->left, value) : bst_search(node->right, value);
}

void bst_free(BSTNode *node) {
    if (!node) return;
    bst_free(node->left);
    bst_free(node->right);
    free(node);
}

int main(void) {
    BSTNode *root = NULL;
    int values[] = {50, 30, 70, 20, 40};
    for (int i = 0; i < 5; i++) root = bst_insert(root, values[i]);

    printf("%d\n", bst_search(root, 40)); // 1 (true)
    printf("%d\n", bst_search(root, 99)); // 0 (false)
    bst_free(root);
    return 0;
}
```

⏱️ **Complexity**: Balanced tree: search/insert/delete = `O(log n)`. Worst case (all books added in sorted order, forming a straight line): `O(n)` — this is why balanced trees (AVL, Red-Black) exist.

---

## 13. Merge Sort & Quick Sort

🧸 **ELI5 (Merge Sort)**: You have a huge pile of shuffled cards. Split the pile in half, then split those halves in half, keep splitting until each pile has just **1 card** (a single card is always "sorted"). Now merge pairs of tiny sorted piles back together in order, then merge those bigger piles, until you have one giant sorted pile. **Divide, then combine.**

🧸 **ELI5 (Quick Sort)**: Pick one card as the "pivot" (say, a 7). Put all smaller cards in one pile to the left, all bigger cards in a pile to the right. Now do the *same trick* on each smaller pile. Eventually every pile has 1 card and everything is sorted.

🌍 **Real life**: How you'd actually sort a huge stack of exam papers by grade if you had friends helping — split the stack among friends, each sorts their small stack, then you merge the stacks back together in order.

📝 **Pseudocode (Merge Sort)**:
```
function mergeSort(list):
    if length(list) <= 1: return list
    mid = length(list) / 2
    left = mergeSort(list[0:mid])
    right = mergeSort(list[mid:])
    return merge(left, right)

function merge(left, right):
    result = []
    while left and right are not empty:
        append smaller of left[0], right[0] to result
    append any leftovers from left or right
    return result
```

🐹 **Go**:
```go
package main

import "fmt"

func mergeSort(list []int) []int {
	if len(list) <= 1 {
		return list
	}
	mid := len(list) / 2
	left := mergeSort(list[:mid])
	right := mergeSort(list[mid:])
	return merge(left, right)
}

func merge(left, right []int) []int {
	result := make([]int, 0, len(left)+len(right))
	i, j := 0, 0
	for i < len(left) && j < len(right) {
		if left[i] <= right[j] {
			result = append(result, left[i])
			i++
		} else {
			result = append(result, right[j])
			j++
		}
	}
	result = append(result, left[i:]...)
	result = append(result, right[j:]...)
	return result
}

func main() {
	exams := []int{88, 45, 92, 67, 71, 30}
	fmt.Println(mergeSort(exams)) // [30 45 67 71 88 92]
}
```

🇨 **C** (in-place on one array, with two small `malloc`'d scratch buffers per merge — no garbage collector to clean up the intermediate slices Go's version allocates):
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

void merge(int *arr, int lo, int mid, int hi) {
    int n1 = mid - lo + 1, n2 = hi - mid;
    int *left = malloc((size_t)n1 * sizeof(int));
    int *right = malloc((size_t)n2 * sizeof(int));
    memcpy(left, arr + lo, (size_t)n1 * sizeof(int));
    memcpy(right, arr + mid + 1, (size_t)n2 * sizeof(int));

    int i = 0, j = 0, k = lo;
    while (i < n1 && j < n2) arr[k++] = (left[i] <= right[j]) ? left[i++] : right[j++];
    while (i < n1) arr[k++] = left[i++];
    while (j < n2) arr[k++] = right[j++];

    free(left);
    free(right);
}

void merge_sort(int *arr, int lo, int hi) {
    if (lo >= hi) return;
    int mid = lo + (hi - lo) / 2;
    merge_sort(arr, lo, mid);
    merge_sort(arr, mid + 1, hi);
    merge(arr, lo, mid, hi);
}

int main(void) {
    int exams[] = {88, 45, 92, 67, 71, 30};
    merge_sort(exams, 0, 5);
    for (int i = 0; i < 6; i++) printf("%d ", exams[i]); // 30 45 67 71 88 92
    printf("\n");
    return 0;
}
```

⏱️ **Complexity**: Merge Sort: `O(n log n)` always, but needs extra space (`O(n)`). Quick Sort: `O(n log n)` average, `O(n²)` worst case (bad pivot choices), but sorts **in place** (little extra memory) — that's why it's often the default in real libraries.

---

## 14. Graphs

🧸 **ELI5**: A tree is like a family — no loops. A **graph** is like your group of **friends** — you can be friends with anyone, friendships can loop around (your friend's friend can also be your friend directly), and there's no single "boss" at the top.

🌍 **Real life**: A social network (people = dots/"nodes", friendships = lines/"edges"). A road map (cities = nodes, roads = edges). Flight routes between airports.

📝 **Pseudocode**:
```
graph = {
    "Alice": ["Bob", "Carol"],
    "Bob": ["Alice"],
    "Carol": ["Alice"]
}
```

🐹 **Go**:
```go
package main

import "fmt"

type Graph struct {
	adjacency map[string][]string
}

func NewGraph() *Graph {
	return &Graph{adjacency: make(map[string][]string)}
}

func (g *Graph) AddFriendship(a, b string) {
	g.adjacency[a] = append(g.adjacency[a], b)
	g.adjacency[b] = append(g.adjacency[b], a)
}

func main() {
	social := NewGraph()
	social.AddFriendship("Alice", "Bob")
	social.AddFriendship("Alice", "Carol")
	fmt.Println(social.adjacency["Alice"]) // [Bob Carol]
}
```

🇨 **C** — without a built-in map, people become integer indices (0, 1, 2...) into a fixed-size names/adjacency table; this same integer-indexed shape is reused for BFS/DFS right below:
```c
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define MAX_NODES 16

typedef struct {
    char names[MAX_NODES][32];
    int  adj[MAX_NODES][MAX_NODES];       // adj[i][j] = 1 if a friendship exists
    int  adj_list[MAX_NODES][MAX_NODES];  // neighbor indices, in insertion order
    int  adj_count[MAX_NODES];
    int  num_nodes;
} Graph;

int graph_add_node(Graph *g, const char *name) {
    strncpy(g->names[g->num_nodes], name, 31);
    return g->num_nodes++;
}

void graph_add_friendship(Graph *g, int a, int b) {
    if (!g->adj[a][b]) { g->adj[a][b] = 1; g->adj_list[a][g->adj_count[a]++] = b; }
    if (!g->adj[b][a]) { g->adj[b][a] = 1; g->adj_list[b][g->adj_count[b]++] = a; }
}

int main(void) {
    Graph social = {0};
    int alice = graph_add_node(&social, "Alice");
    int bob   = graph_add_node(&social, "Bob");
    int carol = graph_add_node(&social, "Carol");

    graph_add_friendship(&social, alice, bob);
    graph_add_friendship(&social, alice, carol);

    for (int i = 0; i < social.adj_count[alice]; i++) {
        printf("%s ", social.names[social.adj_list[alice][i]]); // Bob Carol
    }
    printf("\n");
    return 0;
}
```

⏱️ **Complexity**: Storage = `O(V + E)` (V = number of people/nodes, E = number of friendships/edges).

---

## 15. BFS & DFS (Graph/Tree Traversal)

🧸 **ELI5 (BFS - Breadth First Search)**: You're spreading a rumor at school. First you tell your close friends (1 step away). Then *they* tell *their* friends (2 steps away). You spread outward in **rings**, layer by layer — like ripples in a pond.

🧸 **ELI5 (DFS - Depth First Search)**: You're exploring a maze. You pick a path and go **as deep as possible** until you hit a dead end, then backtrack to the last fork and try a different path. You commit fully to one direction before trying another.

🌍 **Real life**: BFS = "how many people are within 2 friendship-steps of me?" (LinkedIn's "2nd degree connections"), finding the **shortest path** in an unweighted maze. DFS = solving a maze by hand, exploring a file system recursively, detecting cycles.

📝 **Pseudocode (BFS)**:
```
function bfs(graph, start):
    visited = {start}
    queue = [start]
    while queue is not empty:
        current = dequeue(queue)
        print(current)
        for neighbor in graph[current]:
            if neighbor not in visited:
                visited.add(neighbor)
                enqueue(queue, neighbor)
```

📝 **Pseudocode (DFS)**:
```
function dfs(graph, current, visited={}):
    if current in visited: return
    visited.add(current)
    print(current)
    for neighbor in graph[current]:
        dfs(graph, neighbor, visited)
```

🐹 **Go (BFS)**:
```go
package main

import "fmt"

func bfs(graph map[string][]string, start string) {
	visited := map[string]bool{start: true}
	queue := []string{start}

	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		fmt.Println("visiting:", current)

		for _, neighbor := range graph[current] {
			if !visited[neighbor] {
				visited[neighbor] = true
				queue = append(queue, neighbor)
			}
		}
	}
}

func main() {
	friendGraph := map[string][]string{
		"Alice": {"Bob", "Carol"},
		"Bob":   {"Alice", "Dave"},
		"Carol": {"Alice"},
		"Dave":  {"Bob"},
	}
	bfs(friendGraph, "Alice") // Alice, Bob, Carol, Dave (ring by ring)
}
```

🐹 **Go (DFS)**:
```go
package main

import "fmt"

func dfs(graph map[string][]string, current string, visited map[string]bool) {
	if visited[current] {
		return
	}
	visited[current] = true
	fmt.Println("visiting:", current)

	for _, neighbor := range graph[current] {
		dfs(graph, neighbor, visited)
	}
}

func main() {
	maze := map[string][]string{
		"Start": {"A", "B"},
		"A":     {"Start", "DeadEnd"},
		"B":     {"Start", "Exit"},
	}
	dfs(maze, "Start", map[string]bool{}) // dives deep down one path first
}
```

🇨 **C** (both BFS and DFS on the same integer-indexed adjacency-list `Graph` from #14 — BFS's queue is a plain array with `head`/`tail` indices, exactly Section 23.2's ring-buffer idea; DFS is recursion, exactly #9's "same problem, smaller size"):
```c
#include <stdio.h>
#include <string.h>

#define MAX_NODES 16

typedef struct {
    char names[MAX_NODES][32];
    int  adj_list[MAX_NODES][MAX_NODES];
    int  adj_count[MAX_NODES];
    int  num_nodes;
} Graph;

int graph_add_node(Graph *g, const char *name) {
    strncpy(g->names[g->num_nodes], name, 31);
    return g->num_nodes++;
}

void graph_add_edge(Graph *g, int a, int b) {
    g->adj_list[a][g->adj_count[a]++] = b;
    g->adj_list[b][g->adj_count[b]++] = a;
}

void bfs(Graph *g, int start) {
    int visited[MAX_NODES] = {0};
    int queue[MAX_NODES], head = 0, tail = 0;

    visited[start] = 1;
    queue[tail++] = start;

    while (head < tail) {
        int current = queue[head++];
        printf("visiting: %s\n", g->names[current]);
        for (int i = 0; i < g->adj_count[current]; i++) {
            int neighbor = g->adj_list[current][i];
            if (!visited[neighbor]) { visited[neighbor] = 1; queue[tail++] = neighbor; }
        }
    }
}

void dfs(Graph *g, int current, int *visited) {
    if (visited[current]) return;
    visited[current] = 1;
    printf("visiting: %s\n", g->names[current]);
    for (int i = 0; i < g->adj_count[current]; i++) {
        dfs(g, g->adj_list[current][i], visited);
    }
}

int main(void) {
    Graph friendGraph = {0};
    int alice = graph_add_node(&friendGraph, "Alice");
    int bob   = graph_add_node(&friendGraph, "Bob");
    int carol = graph_add_node(&friendGraph, "Carol");
    int dave  = graph_add_node(&friendGraph, "Dave");
    graph_add_edge(&friendGraph, alice, bob);
    graph_add_edge(&friendGraph, alice, carol);
    graph_add_edge(&friendGraph, bob, dave);

    printf("-- BFS --\n");
    bfs(&friendGraph, alice); // Alice, Bob, Carol, Dave (ring by ring)

    printf("-- DFS --\n");
    int visited[MAX_NODES] = {0};
    dfs(&friendGraph, alice, visited); // dives deep down one path first
    return 0;
}
```

⏱️ **Complexity**: Both BFS and DFS = `O(V + E)` — you visit every person/place and every friendship/path once.

---

# 🔴 ADVANCED

Everything so far had one obvious way to solve it. This tier is about
*choice* — problems where many valid paths exist and you need a strategy:
try everything but remember what you've already solved (DP), grab the
locally-best option and hope it's globally best too (Greedy), explore and
back up on dead ends (Backtracking), or always know the most urgent item
without re-sorting the whole pile (Heaps). These eight topics are also where
technical interviews spend most of their time, for exactly that reason.

## 16. Dynamic Programming (DP)

🧸 **ELI5**: Imagine you're climbing a staircase and someone asks "how many ways can you climb 30 steps if you can take 1 or 2 steps at a time?" Instead of recalculating from scratch every time (which is *super* slow), you **write the answer for small staircases on sticky notes** (2 steps = 2 ways, 3 steps = 3 ways...) and reuse those notes to quickly build up the answer for 30 steps. **Never solve the same sub-problem twice — remember it!**

🌍 **Real life**: A cashier making change with the *fewest coins possible* remembers the best way to make smaller amounts and builds up. GPS apps remembering the best sub-routes to avoid recalculating the same road segment's cost repeatedly. Netflix computing "best next recommendation" using cached sub-results.

📝 **Pseudocode (Fibonacci with memoization)**:
```
memo = {}
function fib(n):
    if n <= 1: return n
    if n in memo: return memo[n]     // check the sticky note first!
    memo[n] = fib(n-1) + fib(n-2)    // write a new sticky note
    return memo[n]
```

🐹 **Go**:
```go
package main

import "fmt"

func fib(n int, memo map[int]int) int {
	if n <= 1 {
		return n
	}
	if val, found := memo[n]; found { // check the sticky note
		return val
	}
	result := fib(n-1, memo) + fib(n-2, memo)
	memo[n] = result // write the sticky note for next time
	return result
}

func main() {
	memo := make(map[int]int)
	fmt.Println(fib(30, memo)) // 832040 — instant, thanks to memoization
}
```

🇨 **C** (the "sticky notes" are a plain array here — `have[n]` marks whether `memo[n]` has been written yet, since a C array has no built-in "key not found" like a Go map does):
```c
#include <stdio.h>

#define MEMO_SIZE 64
static long memo[MEMO_SIZE];
static int  have[MEMO_SIZE];

long fib(int n) {
    if (n <= 1) return n;
    if (have[n]) return memo[n];        // check the sticky note first!
    long result = fib(n - 1) + fib(n - 2);
    memo[n] = result;                    // write a new sticky note
    have[n] = 1;
    return result;
}

int main(void) {
    printf("%ld\n", fib(30)); // 832040 — instant, thanks to memoization
    return 0;
}
```

⏱️ **Complexity**: Without memoization, naive Fibonacci is `O(2ⁿ)` (extremely slow — it re-solves the same problem millions of times). With memoization: `O(n)` — each sub-problem solved exactly once.

---

## 17. Greedy Algorithms

🧸 **ELI5**: A greedy algorithm always grabs the **best-looking option right now**, without worrying about the future. Like when a cashier gives you change: for $0.67, they hand you a quarter, then a quarter, then a dime, a nickel, and 2 pennies — always picking the **biggest coin that still fits**, one decision at a time, never going back to reconsider.

🌍 **Real life**: Making change with the fewest coins (works great with US coins!). Picking the meeting room schedule that lets you attend the most meetings (always pick the meeting that ends soonest). Google Maps sometimes uses greedy-like shortcuts for quick estimates.

⚠️ **Watch out**: Greedy is fast but doesn't *always* give the perfect answer for every problem — it works great for some problems (coin change with normal coins, interval scheduling) and fails for others (coin change with weird coin values like {1, 3, 4} for target 6 — greedy might pick 4+1+1=3 coins instead of the true best, 3+3=2 coins).

📝 **Pseudocode (coin change, greedy)**:
```
function greedyChange(coins, amount):
    sort coins descending
    result = []
    for coin in coins:
        while amount >= coin:
            result.add(coin)
            amount -= coin
    return result
```

🐹 **Go**:
```go
package main

import (
	"fmt"
	"sort"
)

func greedyChange(coins []int, amount int) []int {
	sort.Sort(sort.Reverse(sort.IntSlice(coins)))
	var used []int
	for _, coin := range coins {
		for amount >= coin {
			used = append(used, coin)
			amount -= coin
		}
	}
	return used
}

func main() {
	fmt.Println(greedyChange([]int{25, 10, 5, 1}, 67))
	// [25 25 10 5 1 1] -> 67 cents using 6 coins
}
```

🇨 **C** (coins pre-sorted descending — sorting them in C means either a hand-rolled sort like #7's bubble sort or `qsort` from Section 11.2, so this version takes the sort as a precondition to keep the greedy logic itself front and center):
```c
#include <stdio.h>

int greedy_change(const int *coins, int n, int amount, int *used, int used_cap) {
    int count = 0;
    for (int i = 0; i < n && count < used_cap; i++) {
        while (amount >= coins[i] && count < used_cap) {
            used[count++] = coins[i];
            amount -= coins[i];
        }
    }
    return count;
}

int main(void) {
    int coins[] = {25, 10, 5, 1}; // pre-sorted descending
    int used[16];
    int count = greedy_change(coins, 4, 67, used, 16);
    for (int i = 0; i < count; i++) printf("%d ", used[i]); // 25 25 10 5 1 1
    printf("\n");
    return 0;
}
```

⏱️ **Complexity**: Usually `O(n log n)` (dominated by the sort) — very fast, since it never backtracks or reconsiders.

---

## 18. Backtracking

🧸 **ELI5**: You're solving a maze (or a Sudoku puzzle). You try a path. If it leads to a dead end, you **walk back** (backtrack) to your last choice and try a *different* path instead. You keep exploring and undoing until you find a path that works — or discover none exist.

🌍 **Real life**: Solving a Sudoku puzzle — try a number, if it breaks a rule later, erase it and try the next number. Solving a maze on paper with a pencil, backing up when you hit a wall. A GPS trying alternate routes when one turns out to be blocked.

📝 **Pseudocode (N-Queens style: place items without conflicts)**:
```
function solve(board, row):
    if row == size(board): return true    // all placed successfully!
    for col in 0..size(board):
        if isSafe(board, row, col):
            place(board, row, col)         // try it
            if solve(board, row + 1): return true
            remove(board, row, col)        // backtrack — undo and try next
    return false                            // no option worked from here
```

🐹 **Go (N-Queens: place N queens on a chessboard so none attack each other)**:
```go
package main

import "fmt"

func isSafe(cols []int, row, col int) bool {
	for r := 0; r < row; r++ {
		c := cols[r]
		if c == col || r-c == row-col || r+c == row+col {
			return false // same column or same diagonal — queens attack!
		}
	}
	return true
}

func solve(cols []int, row, n int) bool {
	if row == n {
		return true // placed every queen safely
	}
	for col := 0; col < n; col++ {
		if isSafe(cols, row, col) {
			cols[row] = col       // try placing the queen here
			if solve(cols, row+1, n) {
				return true
			}
			// backtrack is implicit: next loop iteration overwrites cols[row]
		}
	}
	return false // dead end — go back further
}

func main() {
	n := 4
	cols := make([]int, n)
	if solve(cols, 0, n) {
		fmt.Println("Queen positions (row -> col):", cols)
	}
}
```

🇨 **C** — line for line the same recursion as the Go version; backtracking here is *literally* "the loop moves to its next iteration and overwrites `cols[row]`," with no explicit "undo" step needed:
```c
#include <stdio.h>

int is_safe(const int *cols, int row, int col) {
    for (int r = 0; r < row; r++) {
        int c = cols[r];
        if (c == col || r - c == row - col || r + c == row + col) return 0; // same column or diagonal
    }
    return 1;
}

int solve(int *cols, int row, int n) {
    if (row == n) return 1;               // placed every queen safely
    for (int col = 0; col < n; col++) {
        if (is_safe(cols, row, col)) {
            cols[row] = col;               // try placing the queen here
            if (solve(cols, row + 1, n)) return 1;
            // backtrack is implicit: next loop iteration overwrites cols[row]
        }
    }
    return 0;                               // dead end — go back further
}

int main(void) {
    int n = 4;
    int cols[4];
    if (solve(cols, 0, n)) {
        printf("Queen positions (row -> col):");
        for (int i = 0; i < n; i++) printf(" %d", cols[i]);
        printf("\n");
    }
    return 0;
}
```

⏱️ **Complexity**: Often exponential in the worst case (`O(kⁿ)`-ish) — but backtracking prunes dead ends early, so it's much faster in practice than trying every possibility blindly.

---

## 19. Heaps / Priority Queues

🧸 **ELI5**: Imagine a hospital emergency room. It's *not* first-come-first-served — the person with the **worst injury** gets seen first, no matter when they arrived. A heap is a special box that always lets you grab the "most urgent" item instantly, and it reorganizes itself automatically whenever something new comes in.

🌍 **Real life**: Hospital triage. Task schedulers in your OS deciding which program gets the CPU next. "Top K" trending topics on Twitter. Dijkstra's algorithm (next!) uses a heap to always process the "closest" unexplored city next.

📝 **Pseudocode**:
```
heap = new MinHeap()
heap.insert(5)
heap.insert(1)
heap.insert(3)
heap.extractMin()   // returns 1, the smallest/most urgent
```

🐹 **Go** (using the standard library's `container/heap`):
```go
package main

import (
	"container/heap"
	"fmt"
)

type PatientHeap []int // lower number = more urgent

func (h PatientHeap) Len() int            { return len(h) }
func (h PatientHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h PatientHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *PatientHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *PatientHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

func main() {
	er := &PatientHeap{5, 1, 3}
	heap.Init(er)
	heap.Push(er, 0) // a critical patient just arrived!

	for er.Len() > 0 {
		fmt.Println("treating urgency level:", heap.Pop(er))
	}
	// treating urgency level: 0, 1, 3, 5
}
```

🇨 **C** — C has no `container/heap` to import, so this *is* what that package does internally: an array where index `i`'s children live at `2i+1`/`2i+2`, kept valid by "bubbling" a value up or down after every change:
```c
#include <stdio.h>

#define HEAP_CAP 64

typedef struct { int items[HEAP_CAP]; int size; } MinHeap;

void heap_swap(MinHeap *h, int i, int j) {
    int t = h->items[i]; h->items[i] = h->items[j]; h->items[j] = t;
}

void heap_push(MinHeap *h, int value) {
    int i = h->size++;
    h->items[i] = value;
    while (i > 0) {                              // bubble up to restore the min-heap property
        int parent = (i - 1) / 2;
        if (h->items[parent] <= h->items[i]) break;
        heap_swap(h, parent, i);
        i = parent;
    }
}

int heap_pop(MinHeap *h) {
    int top = h->items[0];
    h->items[0] = h->items[--h->size];
    int i = 0;
    while (1) {                                    // bubble down to restore the min-heap property
        int left = 2 * i + 1, right = 2 * i + 2, smallest = i;
        if (left < h->size && h->items[left] < h->items[smallest]) smallest = left;
        if (right < h->size && h->items[right] < h->items[smallest]) smallest = right;
        if (smallest == i) break;
        heap_swap(h, i, smallest);
        i = smallest;
    }
    return top;
}

int main(void) {
    MinHeap er = {0};
    heap_push(&er, 5);
    heap_push(&er, 1);
    heap_push(&er, 3);
    heap_push(&er, 0);  // a critical patient just arrived!

    while (er.size > 0) printf("treating urgency level: %d\n", heap_pop(&er));
    // treating urgency level: 0, 1, 3, 5
    return 0;
}
```

⏱️ **Complexity**: Insert / extract-min = `O(log n)`. Peek at the top = `O(1)`.

---

## 20. Union-Find (Disjoint Set)

🧸 **ELI5**: Imagine everyone at a party starts in their own tiny friend group of 1. Whenever two people become friends, you **merge** their whole friend groups into one. Union-Find quickly answers: "are Alice and Bob in the same friend group?" and "merge these two groups together" — both almost instantly, even with millions of people.

🌍 **Real life**: Detecting if adding a road would create a cycle (used in building efficient road networks — Kruskal's algorithm). Grouping friends into "friend circles" on a social network. Detecting connected components in an image (like the "magic wand" tool in Photoshop grouping similar-colored pixels).

📝 **Pseudocode**:
```
function find(x):
    if parent[x] != x:
        parent[x] = find(parent[x])   // path compression: flatten the group
    return parent[x]

function union(x, y):
    rootX = find(x)
    rootY = find(y)
    if rootX != rootY:
        parent[rootX] = rootY          // merge the two groups
```

🐹 **Go**:
```go
package main

import "fmt"

type UnionFind struct {
	parent []int
}

func NewUnionFind(n int) *UnionFind {
	parent := make([]int, n)
	for i := range parent {
		parent[i] = i // everyone starts as their own group leader
	}
	return &UnionFind{parent: parent}
}

func (uf *UnionFind) Find(x int) int {
	if uf.parent[x] != x {
		uf.parent[x] = uf.Find(uf.parent[x]) // flatten the chain on the way up
	}
	return uf.parent[x]
}

func (uf *UnionFind) Union(x, y int) {
	rootX, rootY := uf.Find(x), uf.Find(y)
	if rootX != rootY {
		uf.parent[rootX] = rootY // merge the friend groups
	}
}

func main() {
	people := NewUnionFind(5) // people 0,1,2,3,4
	people.Union(0, 1)        // 0 and 1 become friends
	people.Union(1, 2)        // now 0,1,2 are one group

	fmt.Println(people.Find(0) == people.Find(2)) // true, same group
	fmt.Println(people.Find(0) == people.Find(3)) // false, different group
}
```

🇨 **C** — the "path compression" line (`uf->parent[x] = uf_find(uf, uf->parent[x])`) is the entire trick: every lookup flattens the chain a little more, so future lookups get faster:
```c
#include <stdio.h>

typedef struct { int parent[64]; } UnionFind;

void uf_init(UnionFind *uf, int n) {
    for (int i = 0; i < n; i++) uf->parent[i] = i; // everyone starts as their own group leader
}

int uf_find(UnionFind *uf, int x) {
    if (uf->parent[x] != x) uf->parent[x] = uf_find(uf, uf->parent[x]); // flatten on the way up
    return uf->parent[x];
}

void uf_union(UnionFind *uf, int x, int y) {
    int root_x = uf_find(uf, x), root_y = uf_find(uf, y);
    if (root_x != root_y) uf->parent[root_x] = root_y; // merge the friend groups
}

int main(void) {
    UnionFind people;
    uf_init(&people, 5); // people 0,1,2,3,4
    uf_union(&people, 0, 1); // 0 and 1 become friends
    uf_union(&people, 1, 2); // now 0,1,2 are one group

    printf("%d\n", uf_find(&people, 0) == uf_find(&people, 2)); // 1 (true), same group
    printf("%d\n", uf_find(&people, 0) == uf_find(&people, 3)); // 0 (false), different group
    return 0;
}
```

⏱️ **Complexity**: With path compression, find/union are nearly `O(1)` (technically `O(α(n))`, an incredibly slow-growing function — for all practical purposes, instant).

---

## 21. Trie (Prefix Tree)

🧸 **ELI5**: A trie is like a tree where each **branch is a single letter**, and words that share the same beginning **share the same branches**. "CAT" and "CAR" share the "CA" branch, then split into "T" and "R". This is exactly how your phone keyboard guesses what you're about to type.

🌍 **Real life**: Autocomplete/predictive text on your phone. Search engine "did you mean" and autosuggest dropdowns. Spell checkers.

📝 **Pseudocode**:
```
trieNode = { children: {}, isEndOfWord: false }

function insert(root, word):
    node = root
    for letter in word:
        if letter not in node.children:
            node.children[letter] = new trieNode()
        node = node.children[letter]
    node.isEndOfWord = true

function search(root, word):
    node = root
    for letter in word:
        if letter not in node.children: return false
        node = node.children[letter]
    return node.isEndOfWord
```

🐹 **Go**:
```go
package main

import "fmt"

type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
}

func newTrieNode() *TrieNode {
	return &TrieNode{children: make(map[rune]*TrieNode)}
}

type Trie struct {
	root *TrieNode
}

func NewTrie() *Trie {
	return &Trie{root: newTrieNode()}
}

func (t *Trie) Insert(word string) {
	node := t.root
	for _, letter := range word {
		if node.children[letter] == nil {
			node.children[letter] = newTrieNode()
		}
		node = node.children[letter]
	}
	node.isWord = true
}

func (t *Trie) Search(word string) bool {
	node := t.root
	for _, letter := range word {
		if node.children[letter] == nil {
			return false
		}
		node = node.children[letter]
	}
	return node.isWord
}

func main() {
	keyboard := NewTrie()
	keyboard.Insert("cat")
	keyboard.Insert("car")

	fmt.Println(keyboard.Search("cat")) // true
	fmt.Println(keyboard.Search("ca"))  // false (only a prefix, not a full word)
}
```

🇨 **C** — Go's version uses a `map[rune]*TrieNode` per node; C swaps that for a fixed 26-pointer array indexed by `letter - 'a'`, trading "any Unicode letter" for "no hash map needed":
```c
#include <stdio.h>
#include <stdlib.h>

typedef struct TrieNode {
    struct TrieNode *children[26];
    int is_word;
} TrieNode;

TrieNode *trie_new_node(void) { return calloc(1, sizeof(TrieNode)); }

void trie_insert(TrieNode *root, const char *word) {
    TrieNode *node = root;
    for (int i = 0; word[i]; i++) {
        int idx = word[i] - 'a';
        if (!node->children[idx]) node->children[idx] = trie_new_node();
        node = node->children[idx];
    }
    node->is_word = 1;
}

int trie_search(TrieNode *root, const char *word) {
    TrieNode *node = root;
    for (int i = 0; word[i]; i++) {
        int idx = word[i] - 'a';
        if (!node->children[idx]) return 0;
        node = node->children[idx];
    }
    return node->is_word;
}

void trie_free(TrieNode *node) {
    if (!node) return;
    for (int i = 0; i < 26; i++) trie_free(node->children[i]);
    free(node);
}

int main(void) {
    TrieNode *keyboard = trie_new_node();
    trie_insert(keyboard, "cat");
    trie_insert(keyboard, "car");

    printf("%d\n", trie_search(keyboard, "cat")); // 1 (true)
    printf("%d\n", trie_search(keyboard, "ca"));   // 0 (false — only a prefix, not a full word)

    trie_free(keyboard);
    return 0;
}
```

⏱️ **Complexity**: Insert/search = `O(k)` where `k` is the word length — independent of how many *other* words are stored! That's why autocomplete feels instant even with a huge dictionary.

---

## 22. Dijkstra's Algorithm (Shortest Path)

🧸 **ELI5**: You're using GPS to find the fastest way to grandma's house, and some roads are slower than others (traffic, distance). Dijkstra's algorithm always explores the **closest unvisited place first** (using a heap/priority queue, like #19!), gradually building the shortest-distance map to *everywhere* from your starting point.

🌍 **Real life**: GPS/Google Maps shortest route calculation. Network routers finding the fastest path for your internet data to travel. Flight-booking sites finding the cheapest connecting flights.

📝 **Pseudocode**:
```
function dijkstra(graph, start):
    distances = { all nodes: infinity }
    distances[start] = 0
    priorityQueue = [(0, start)]

    while priorityQueue is not empty:
        (dist, current) = extractMin(priorityQueue)
        for (neighbor, weight) in graph[current]:
            newDist = dist + weight
            if newDist < distances[neighbor]:
                distances[neighbor] = newDist
                priorityQueue.add((newDist, neighbor))
    return distances
```

🐹 **Go**:
```go
package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Edge struct {
	to     string
	weight int
}

type Item struct {
	node string
	dist int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(graph map[string][]Edge, start string) map[string]int {
	distances := make(map[string]int)
	for node := range graph {
		distances[node] = math.MaxInt32
	}
	distances[start] = 0

	pq := &PriorityQueue{{node: start, dist: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Item)
		for _, edge := range graph[current.node] {
			newDist := current.dist + edge.weight
			if newDist < distances[edge.to] {
				distances[edge.to] = newDist
				heap.Push(pq, Item{node: edge.to, dist: newDist})
			}
		}
	}
	return distances
}

func main() {
	roadMap := map[string][]Edge{
		"Home":   {{"School", 5}, {"Park", 2}},
		"Park":   {{"School", 1}},
		"School": {},
	}
	fmt.Println(dijkstra(roadMap, "Home"))
	// Home:0, Park:2, School:3 (Home->Park->School is faster than Home->School directly!)
}
```

🇨 **C** — this version picks the closest unvisited place with a plain `O(V)` scan each round instead of #19's heap, trading `O((V+E) log V)` for a simpler `O(V²)`; wiring the heap from #19 into this loop instead of the scan is exactly how you'd get the faster version:
```c
#include <stdio.h>
#include <limits.h>
#include <string.h>

#define MAX_NODES 8

typedef struct {
    char names[MAX_NODES][32];
    int  weight[MAX_NODES][MAX_NODES]; // 0 = no direct road
    int  num_nodes;
} RoadMap;

int roadmap_add_node(RoadMap *g, const char *name) {
    strncpy(g->names[g->num_nodes], name, 31);
    return g->num_nodes++;
}

void roadmap_add_road(RoadMap *g, int a, int b, int weight) {
    g->weight[a][b] = weight;   // directed; call twice for a two-way road
}

void dijkstra(RoadMap *g, int start, int *dist) {
    int visited[MAX_NODES] = {0};
    for (int i = 0; i < g->num_nodes; i++) dist[i] = INT_MAX;
    dist[start] = 0;

    for (int iter = 0; iter < g->num_nodes; iter++) {
        int current = -1;
        for (int i = 0; i < g->num_nodes; i++) {      // pick the closest unvisited place —
            if (!visited[i] && (current == -1 || dist[i] < dist[current])) current = i;
        }
        if (current == -1 || dist[current] == INT_MAX) break;
        visited[current] = 1;

        for (int next = 0; next < g->num_nodes; next++) {
            if (g->weight[current][next] > 0 && dist[current] + g->weight[current][next] < dist[next]) {
                dist[next] = dist[current] + g->weight[current][next];
            }
        }
    }
}

int main(void) {
    RoadMap roadMap = {0};
    int home = roadmap_add_node(&roadMap, "Home");
    int park = roadmap_add_node(&roadMap, "Park");
    int school = roadmap_add_node(&roadMap, "School");

    roadmap_add_road(&roadMap, home, school, 5);
    roadmap_add_road(&roadMap, home, park, 2);
    roadmap_add_road(&roadMap, park, school, 1);

    int dist[MAX_NODES];
    dijkstra(&roadMap, home, dist);

    for (int i = 0; i < roadMap.num_nodes; i++) printf("%s:%d ", roadMap.names[i], dist[i]);
    printf("\n"); // Home:0 Park:2 School:3 — Home->Park->School beats Home->School directly
    return 0;
}
```

⏱️ **Complexity**: `O((V + E) log V)` using a heap — fast enough for GPS to recalculate routes on the fly.

---

## 23. Topological Sort

🧸 **ELI5**: You can't put on your shoes before your socks. Some tasks **must** happen before others. Topological sort takes a pile of "this-before-that" rules and gives you **one valid order** to do everything in, respecting every rule.

🌍 **Real life**: Getting dressed in the right order (underwear → pants → belt). University course prerequisites (must take Calculus 1 before Calculus 2). Build systems compiling code (compile library A before the app that depends on it). Spreadsheet formulas that depend on other cells.

📝 **Pseudocode (Kahn's algorithm, using in-degrees + a queue)**:
```
function topologicalSort(graph):
    inDegree = count incoming edges for every node
    queue = all nodes with inDegree == 0     // no prerequisites
    order = []

    while queue is not empty:
        node = dequeue(queue)
        order.add(node)
        for neighbor in graph[node]:
            inDegree[neighbor] -= 1
            if inDegree[neighbor] == 0:
                enqueue(queue, neighbor)
    return order
```

🐹 **Go**:
```go
package main

import "fmt"

func topologicalSort(graph map[string][]string) []string {
	inDegree := make(map[string]int)
	for node := range graph {
		inDegree[node] = 0
	}
	for _, deps := range graph {
		for _, dep := range deps {
			inDegree[dep]++
		}
	}

	var queue []string
	for node, degree := range inDegree {
		if degree == 0 {
			queue = append(queue, node) // no prerequisites — start here
		}
	}

	var order []string
	for len(queue) > 0 {
		current := queue[0]
		queue = queue[1:]
		order = append(order, current)

		for _, next := range graph[current] {
			inDegree[next]--
			if inDegree[next] == 0 {
				queue = append(queue, next)
			}
		}
	}
	return order
}

func main() {
	gettingDressed := map[string][]string{
		"Underwear": {"Pants"},
		"Pants":     {"Belt"},
		"Socks":     {"Shoes"},
		"Belt":      {},
		"Shoes":     {},
	}
	fmt.Println(topologicalSort(gettingDressed))
	// e.g. [Underwear Socks Pants Shoes Belt] — a valid order respecting all rules
}
```

🇨 **C** (Kahn's algorithm, same integer-indexed `Graph` shape as #14/#15, with the ready-to-schedule queue as a plain array):
```c
#include <stdio.h>
#include <string.h>

#define MAX_NODES 8

typedef struct {
    char names[MAX_NODES][32];
    int  adj_list[MAX_NODES][MAX_NODES];
    int  adj_count[MAX_NODES];
    int  num_nodes;
} Graph;

int graph_add_node(Graph *g, const char *name) {
    strncpy(g->names[g->num_nodes], name, 31);
    return g->num_nodes++;
}

void graph_add_prereq(Graph *g, int before, int after) {
    g->adj_list[before][g->adj_count[before]++] = after; // "before" must happen before "after"
}

int topological_sort(Graph *g, int *order) {
    int in_degree[MAX_NODES] = {0};
    for (int i = 0; i < g->num_nodes; i++) {
        for (int j = 0; j < g->adj_count[i]; j++) in_degree[g->adj_list[i][j]]++;
    }

    int queue[MAX_NODES], head = 0, tail = 0, order_len = 0;
    for (int i = 0; i < g->num_nodes; i++) {
        if (in_degree[i] == 0) queue[tail++] = i; // no prerequisites — start here
    }

    while (head < tail) {
        int current = queue[head++];
        order[order_len++] = current;
        for (int i = 0; i < g->adj_count[current]; i++) {
            int next = g->adj_list[current][i];
            if (--in_degree[next] == 0) queue[tail++] = next;
        }
    }
    return order_len;
}

int main(void) {
    Graph gettingDressed = {0};
    int underwear = graph_add_node(&gettingDressed, "Underwear");
    int pants     = graph_add_node(&gettingDressed, "Pants");
    int belt      = graph_add_node(&gettingDressed, "Belt");
    int socks     = graph_add_node(&gettingDressed, "Socks");
    int shoes     = graph_add_node(&gettingDressed, "Shoes");

    graph_add_prereq(&gettingDressed, underwear, pants);
    graph_add_prereq(&gettingDressed, pants, belt);
    graph_add_prereq(&gettingDressed, socks, shoes);

    int order[MAX_NODES];
    int n = topological_sort(&gettingDressed, order);
    for (int i = 0; i < n; i++) printf("%s ", gettingDressed.names[order[i]]);
    printf("\n"); // e.g. Underwear Socks Pants Shoes Belt — a valid order respecting all rules
    return 0;
}
```

⏱️ **Complexity**: `O(V + E)` — visit every task and every "must happen before" rule once.

---

**That's the full toolkit** — 23 topics across three tiers, from egg-carton
arrays to Dijkstra's shortest path. What follows isn't more to learn; it's
the quick-reference for the moment that actually matters: staring at a new
problem and recognizing which of these 23 tools it's asking for.

# Cheat Sheet: When Do I Use What?

| Problem sounds like...                                  | Reach for...              |
|-----------------------------------------------------------|----------------------------|
| "I need instant lookup by a key/name"                    | Hash Map                  |
| "First come, first served"                                | Queue                     |
| "Undo my last action"                                     | Stack                     |
| "The list is sorted, find something fast"                 | Binary Search              |
| "Find the shortest path / fewest steps"                   | BFS or Dijkstra            |
| "Explore every possibility, backing up on dead ends"       | DFS / Backtracking          |
| "This must happen before that"                             | Topological Sort            |
| "I keep recalculating the same sub-answer"                | Dynamic Programming (memoize) |
| "Just grab the biggest/cheapest option each time"          | Greedy                     |
| "Always give me the most urgent/smallest/biggest item next"| Heap / Priority Queue       |
| "Are these two things in the same group?"                  | Union-Find                 |
| "Autocomplete as the user types"                            | Trie                       |
| "Organize hierarchical data (folders, org charts)"          | Tree                       |
| "Model relationships between things (people, cities)"       | Graph                      |

## The one-sentence version of everything

- **Array**: numbered boxes in a row, instant access by number.
- **String**: an array of letters.
- **Linked List**: a chain of clues, each pointing to the next.
- **Stack**: last in, first out — like a pile of pancakes.
- **Queue**: first in, first out — like a checkout line.
- **Linear Search**: check everything, one by one.
- **Bubble Sort**: repeatedly swap neighbors until sorted.
- **Binary Search**: cut the possibilities in half every guess.
- **Recursion**: solve a smaller version of the same problem.
- **Hash Map**: instant lookup via a coat-check ticket.
- **Tree**: a family tree — one path to every relative.
- **BST**: a tree where left is always smaller, right is always bigger.
- **Merge/Quick Sort**: divide the mess into tiny pieces, then combine/sort smartly.
- **Graph**: friendships that can loop around, unlike a family tree.
- **BFS/DFS**: explore ring-by-ring, or dive deep-first.
- **Dynamic Programming**: remember answers so you never redo work.
- **Greedy**: always take the best option available right now.
- **Backtracking**: try, and undo if it doesn't work out.
- **Heap**: always know the most urgent item instantly.
- **Union-Find**: quickly track which friend groups have merged.
- **Trie**: a tree of letters, for lightning-fast prefix lookups.
- **Dijkstra**: GPS logic — always expand the closest place first.
- **Topological Sort**: one valid order that respects all the "do this before that" rules.

---

# 📚 Further Reading

## Books

- **Grokking Algorithms** — Aditya Bhargava. The friendliest possible print introduction — hand-drawn illustrations, no heavy math. Basically this guide's spirit, in book form.
- **A Common-Sense Guide to Data Structures and Algorithms** — Jay Wengrow. Plain-English explanations with a strong focus on *why* Big-O matters in real code.
- **Cracking the Coding Interview** — Gayle Laakmann McDowell. The standard reference once you're comfortable with the basics and want interview-style practice problems.
- **Introduction to Algorithms (CLRS)** — Cormen, Leiserson, Rivest, Stein. The rigorous, formal textbook (a.k.a. "the DSA bible") — dense, but the definitive reference once you want proofs and deep detail.
- **The Algorithm Design Manual** — Steven Skiena. Great for connecting algorithms to real engineering problems, with a "war stories" section of applied case studies.
- **Algorithms, 4th Edition** — Robert Sedgewick & Kevin Wayne. Clean Java-based explanations; pairs well with the free Princeton course below.

## Blogs & Written Guides

- [**VisuAlgo**](https://visualgo.net/) — visualizes data structures and algorithms step-by-step; excellent for building intuition before reading code.
- [**GeeksforGeeks — DSA**](https://www.geeksforgeeks.org/data-structures/) — huge library of explanations, complexity breakdowns, and practice problems per topic.
- [**NeetCode Roadmap**](https://neetcode.io/roadmap) — a curated, structured path through the classic interview problems, organized by pattern (not just by topic).
- **The Go blog & `container/heap`/`sort` package docs** ([pkg.go.dev/container/heap](https://pkg.go.dev/container/heap), [pkg.go.dev/sort](https://pkg.go.dev/sort)) — see how these exact structures are implemented in Go's own standard library.
- [**Big-O Cheat Sheet**](https://www.bigocheatsheet.com/) — a single-page reference for the time/space complexity of nearly every common data structure and algorithm.
- **[`wiki/c-lang/real-life-c-guide.md`](../c-lang/real-life-c-guide.md)** — this guide's C companion. Every 🇨 block above is a stripped-down version of a pattern that guide covers in full: Part III (§17–23) builds linked lists, hash tables, BSTs, and thread-safe structures from scratch; Part VI walks the *exact same data structures* inside real SQLite/Redis/curl source; Part VIII is the memory-leak and defensive-C checklist behind every `malloc`/`free` pair used here.

## Video Tutorials

- [**MIT 6.006 — Introduction to Algorithms**](https://ocw.mit.edu/courses/6-006-introduction-to-algorithms-spring-2020/) (MIT OpenCourseWare) — free, full-rigor university course, lecture videos included.
- [**Abdul Bari — Algorithms Playlist**](https://www.youtube.com/@abdul_bari) (YouTube) — widely loved for extremely clear whiteboard walkthroughs of sorting, graphs, and DP.
- [**NeetCode — Coding Interview Solutions**](https://www.youtube.com/@NeetCode) (YouTube) — short, focused walkthroughs of individual problems mapped to the patterns in this guide.
- [**Back To Back SWE**](https://www.youtube.com/@BackToBackSWE) (YouTube) — deep, intuition-first explanations of dynamic programming, graphs, and heaps.
- [**Princeton Algorithms, Part I & II**](https://www.coursera.org/learn/algorithms-part1) (Coursera, free to audit) — the video companion to the Sedgewick & Wayne book above.
- [**CS50's Introduction to Computer Science**](https://cs50.harvard.edu/x/) (Harvard, free) — not DSA-only, but its data structures/algorithms weeks are a superb, beginner-friendly on-ramp.
