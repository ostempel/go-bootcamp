# Concurrency

## SyncGroup
- easily coordinate go-routines
- `wg.Add(int)`: mark how many go-routines are running
- `defer wg.Done()`: mark go-routine finish -> within go-routine
- `wg.Wait()`: wait for every go-routine finished
