# Fahhh: Your workplace catharsis

## How to run
1. Intall dependencies: `go mod tidy`
2. Build locally: `go build -o fahhh .`
3. Run: `./fahhh` and **fahhh** away your frustration on every keystroke
4. `Ctrl + C` when your manager is starting to notice

## Grant macOS Accessibility Permission (First Run Only)
Because the binary captures keystrokes globally, macOS will prompt for permission:

1. Open System Settings → Privacy & Security → Accessibility.
2. Toggle your terminal application (Terminal, iTerm2, or Cursor) to ON.
3. Rerun `./fahhh`.
