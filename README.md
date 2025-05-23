# Kill The Port 🔍

A cross-platform CLI tool to inspect and manage port usage on your system. Perfect for developers who frequently encounter "port already in use" errors.

## Output

![KillThePort](https://github.com/user-attachments/assets/aa4c82b5-8f6a-4c77-9568-c965d9c83acc)

## Features ✨

- List all processes using network ports
- Filter by specific port number
- Kill processes occupying ports
- Interactive process selection
- Cross-platform support (Windows & Unix)
- Lightweight alternative to `lsof`/`netstat` + `grep`

## Basic Commands
### Command	Description
| Command | Description |
|---------|-------------|
| `portsleuth` | List all active network connections |
| `portsleuth --kill` | Interactive mode to select and kill processes |
| `portsleuth --kill :3000` | Kill all processes using port 3000 |

