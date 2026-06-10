# SharkEyes Go Library

[![Go Reference](https://pkg.go.dev/badge/github.com/SharkEyes-Dev/sharkeyes.go.svg)](https://pkg.go.dev/github.com/SharkEyes-Dev/sharkeyes.go)
[![Go Report Card](https://goreportcard.com/badge/github.com/SharkEyes-Dev/sharkeyes.go)](https://goreportcard.com/report/github.com/SharkEyes-Dev/sharkeyes.go)
[![License](https://img.shields.io/github/license/SharkEyes-Dev/sharkeyes.go)](LICENSE)

Official Go library for integrating with the SharkEyes web form bot protection service. This package allows you to easily verify security tokens on your backend application.

> **Security First & Privacy Friendly.** We protect your forms while preserving user confidentiality. Privacy is our priority; audit logs are anonymized and feature short retention cycles.

---

## Installation

Ensure you have Go 1.24 or higher installed, then run the following command in your project terminal:

```bash
go get [github.com/SharkEyes-Dev/sharkeyes.go@v1.0.1](https://github.com/SharkEyes-Dev/sharkeyes.go@v1.0.1)

```

---

## Quick Start

To begin using the library, you need to obtain an API Key from your [SharkEyes](https://sharkeyes.dev) dashboard.

```go
package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SharkEyes-Dev/sharkeyes.go"
)

func init() {
	// 1. Initialize the library with your secret API key before starting the server
	sharkeyes.Configure("your_secret_api_key_here", "", 5*time.Second)
}

func handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// 2. Extract the verification token sent from the frontend
	token := r.FormValue("sharkeyes_token")

	// 3. Verify the token against the SharkEyes API
	isValid, reason := sharkeyes.Verify(token, r)
	if !isValid {
		// If it is a bot or the token is invalid, block the request
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "Access denied: %s\n", reason)
		return
	}

	// 4. Processing continues normally if verification succeeds
	fmt.Fprintln(w, "Form submitted successfully!")
}

func main() {
	http.HandleFunc("/submit", handler)
	fmt.Println("Server running on :8080")
	http.ListenAndServe(":8080", nil)
}

```

---

## API Reference

### `func Configure(key string, url string, t time.Duration)`

Configures the global SharkEyes client parameters. This function must be executed prior to any calls to `Verify`.

| Parameter | Type | Description |
| --- | --- | --- |
| `key` | `string` | Your private API key. Required. |
| `url` | `string` | Custom API endpoint. Passing `""` defaults to `https://api.sharkeyes.dev/api/v2/server-verify`. |
| `t` | `time.Duration` | Timeout duration for HTTP verification requests. Passing `0` sets the default timeout of 5 seconds. |

---

### `func Verify(token string, r *http.Request) (bool, string)`

Dispatches the provided token to the SharkEyes servers to evaluate potential bot activity.

**Return Values:**

* `bool` — Returns `true` if the request is legitimate (human activity), and `false` if it represents a bot or if a technical failure occurred.
* `string` — Returns a descriptive string indicating the reason for rejection (e.g., `"Enable JavaScript."`, `"Technical error..."`), or an empty string upon successful validation.

> **Warning:** Failing to execute `Configure` or supplying an empty API key before invoking `Verify` triggers a runtime panic. This safety mechanism prevents bypassing bot verification due to server misconfiguration.

---

## Resilience and Error Handling

* **Logging:** In the event of network disruptions or JSON parsing issues, the library prevents application failure. It logs errors via the standard Go logger with a `[SharkEyes]` prefix for system monitoring compatibility, while returning a safe fallback value of `false`.
* **Goroutine Protection:** The integration of a configurable timeout mechanism (defaulting to 5 seconds) guarantees that remote API latency will not block backend server goroutines or degrade request handling capabilities.

---
## Documentation
For advanced configurations, custom middleware implementations, and detailed integration guides, please visit our official documentation at https://docs.sharkeyes.dev
## License

This project is licensed under the MIT License - see the [LICENSE](https://www.google.com/search?q=LICENSE) file for details. Developed by the **SharkEyes** team.

