# Concurrency

- no `panic` in routines -> whole program crashes
  - log err
  - channel for err
  - err within message channel

## SyncGroup
- easily coordinate go-routines
- `wg.Add(int)`: mark how many go-routines are running
- `defer wg.Done()`: mark go-routine finish -> within go-routine
- `wg.Wait()`: wait for every go-routine finished

## Channels
- buffered channels: `messages := make(chan T, size)` -> now how many messages/goroutines are spawned
- unbuffered channels: sender **blocks** till somebody reads -> in same gorutine **Deadlock**
  - channel filled -> `close(c)`
- `for range size`: gather results with `msg, ok := <- messages`
  - ok to test if closed
  - after size times of receiving channel gets close
- `for range messages`:
  - get messages in channel
  - but **doesn't** close channel when all received

## Worker pools
- workers are own goroutines
- goroutine sends jobs and closes job-channel
- worker read job-channel and send results-channel
- separate goroutine waits for workers done and closes results-channel
- main reads results

### Deadlocks
- `wg.Wait()` of worker before reading -> worker can't send because nobody reads
- `for range channel`: reads till channel **closed**

## Pipeline
- Pattern: `func(<-chan T) <-chan T`
- stages encapsulate own goroutine

## Fan-Out/Fan-In
- Fan-Out: goroutines read from same channel (distribute work like worker-poll) but have **own** result channel 
- Fan-In: wait-group counts result channels, merge into one channel and close **after** waitgroup done 

## Worker vs Pipeline
- Pipeline: sequential transformation 
- Worker: parallel processing

**Avoid**:
- send and receive in different goroutines (unbuffered)
- filler of channel also closes
- close result channel only after **all** worker are done -> wait and then close
