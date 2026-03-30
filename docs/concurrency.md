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
- `for range size`: gather results with `msg, ok := <- messages`
  - ok to test if closed
  - after size times of receiving channel gets close
- `for range messages`:
  - get messages in channel
  - but **doesn't** close channel when all received
