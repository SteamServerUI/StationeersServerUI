# Socket Server Usage

This guide explains how to interact with the socket-based `SSUI-API`, which exposes the same HTTP endpoints as the http server over a _Unix socket on Linux_ or a _named pipe on Windows_. This allows local, lightweight access to the API without network overhead or authentication, mainly for the planned plugin system but also for advancedscripting or external local tools. The server is built with Go’s standard library on Linux and `microsoft/go-winio` on Windows for cross-platform compatibility.

## Overview

The socket server reuses the HTTP routes defined in `routes.go` but serves them over:
- **Linux**: A Unix socket at `/tmp/ssui.sock`.
- **Windows**: A named pipe at `\\.\pipe\ssui`.

This is similar to how Docker uses `/var/run/docker.sock` for local API access.

## Notes
- The socket server skips authentication, assuming local access is secure (like Docker’s socket).

## Prerequisites

- **Linux**:
  - `curl` installed for testing.
- **Windows**:
  - PowerShell 7+ for testing.

## Testing the Socket Server

### Linux: Using a Unix Socket

Use `curl` with the `--unix-socket` flag to send HTTP requests to the socket. Example for the `/api/v2/settings` endpoint:

```bash
curl --unix-socket /tmp/ssui.sock http://localhost/api/v2/settings
```

**Expected Output**: JSON response from the `settings.RetrieveSettings` handler, e.g.:
```json
{
  "settings": {
    "some_key": "some_value"
  }
}
```

### Windows: Using a Named Pipe

Use the following PowerShell 7 script to send an HTTP request to the named pipe. Save as `test_namedpipe.ps1` and run it.

```powershell
# test_namedpipe.ps1

$pipeName = "\\.\pipe\ssui"
$endpoint = "/api/v2/server/status"
$hostHeader = "localhost"  # Dummy host for HTTP request

# Create the HTTP request
$request = "GET $endpoint HTTP/1.1`r`nHost: $hostHeader`r`nConnection: close`r`n`r`n"

try {
    # Connect to the named pipe
    $pipe = New-Object System.IO.Pipes.NamedPipeClientStream(".", "ssui-api", [System.IO.Pipes.PipeDirection]::InOut)
    $pipe.Connect(5000)  # 5-second timeout

    # Convert request to bytes and send
    $requestBytes = [System.Text.Encoding]::UTF8.GetBytes($request)
    $pipe.Write($requestBytes, 0, $requestBytes.Length)
    $pipe.Flush()

    # Read response
    $reader = New-Object System.IO.StreamReader($pipe)
    $response = ""
    while ($null -ne ($line = $reader.ReadLine())) {
        $response += "$line`n"
    }

    # Output the response
    Write-Output $response
}
catch {
    Write-Error "Error connecting to named pipe or reading response: $_"
}
finally {
    # Clean up
    if ($null -ne $pipe) {
        $pipe.Close()
        $pipe.Dispose()
    }
    if ($null -ne $reader) {
        $reader.Close()
        $reader.Dispose()
    }
}
```

Run it:
```powershell
.\test_namedpipe.ps1
```

**Expected Output**: HTTP response with headers and JSON body, e.g.:
```
Response from \\.\pipe\ssui/api/v2/server/status:
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 123

{"status":"running","players":0}
```

## Testing Other Endpoints

All routes from `routes.go` (e.g., `/api/v2/server/start`, `/api/v2/backups`) are available. Modify the endpoint in the commands:

- **Linux**:
  ```bash
  curl --unix-socket /tmp/ssui.sock http://localhost/api/v2/backups
  ```
- **Windows**: Edit `$endpoint` in `test_namedpipe.ps1`, e.g., `$endpoint = "/api/v2/backups"`.

For POST requests (e.g., `/api/v2/server/start`), add a JSON payload:
- **Linux**:
  ```bash
  curl --unix-socket /tmp/ssui.sock -X POST -H "Content-Type: application/json" -d '{"action":"start"}' http://localhost/api/v2/server/start
  ```
- **Windows**: Update the PowerShell script’s `$request`:
  ```powershell
  $body = '{"action":"start"}'
  $request = "POST $endpoint HTTP/1.1`r`nHost: $hostHeader`r`nContent-Type: application/json`r`nContent-Length: $($body.Length)`r`n`r`n$body"
  ```
