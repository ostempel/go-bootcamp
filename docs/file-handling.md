# File handling

## Files
- close file after done `defer f.Close()`
- simple file handling
  - `os.ReadFile`: read file as `[]byte`
  - `os.WriteFile`: truncate file or create and write `[]byte`
  - `os.IsNotExist(err)`: check file existence

## bufio
- `Scanner` read file line by line -> custom scanner possible (e.g.: `ScanWords`)
  - `s.Scan()` - loop: reads till next split token 
  - `s.Error`: after loop check for errors
