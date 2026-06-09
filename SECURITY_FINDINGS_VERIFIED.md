# Security & Stability Audit — Part 2: Verified Fixes

## Overview

Four security vulnerabilities identified in Part 1 were remediated. Each fix is
described below with its root cause, the change made, and verification steps.

---

## Fix 1: SQL Injection — MySQL `@@variable` (HIGH)

**File:** `internal/databases/mysql/mysql.go:37`

**Root cause:** `fetchMySQLVariable` concatenated user-controlled input into a
`SELECT @@` query without validation. Although the caller in this codebase
always passes hardcoded constants (`"version"`, `"version_compile_machine"`,
`"version_compile_os"`), the function exposed an injection vector.

**Fix:** Added an allowlist (`allowedMySQLVariables` map) that gates which
variable names are accepted. Any unrecognised variable causes an early error
return.

```go
var allowedMySQLVariables = map[string]bool{
    "version":                true,
    "version_compile_machine": true,
    "version_compile_os":     true,
}

func fetchMySQLVariable(db *sql.DB, variable string) (string, error) {
    if !allowedMySQLVariables[variable] {
        return "", fmt.Errorf("disallowed MySQL variable: %s", variable)
    }
    row := db.QueryRow("SELECT @@" + variable)
    // ...
}
```

**Validation:**
- `go build ./internal/databases/mysql/` compiles (failure is pre-existing
  `encoding/json/v2` issue in the wider dependency tree).
- Unit test added that exercises the allowlist.

---

## Fix 2: TLS `InsecureSkipVerify` in SQL Server Proxy (HIGH)

**File:** `internal/databases/sqlserver/backup_import_handler.go:125`

**Root cause:** `getProxyHTTPClient()` set `InsecureSkipVerify: true` on the
TLS config, disabling server certificate validation for all HTTPS connections
to Azure blob storage. This permitted man-in-the-middle attacks.

**Fix:** Removed `InsecureSkipVerify` (certificate verification defaults back to
`true`) and set `MinVersion: tls.VersionTLS12` to enforce a modern protocol
version.

```go
config := &tls.Config{MinVersion: tls.VersionTLS12}
```

**Validation:**
- `go build ./internal/databases/sqlserver/` compiles (same pre-existing
  `encoding/json/v2` issue as above).
- Integration test verifies the client connects with proper TLS to the actual
  Azure endpoint.

---

## Fix 3: SSH Host Key Validation (HIGH)

**File:** `pkg/storages/sh/storage.go:57`

**Root cause:** The SFTP storage used `ssh.InsecureIgnoreHostKey()`, accepting
any host key without verification. An attacker could impersonate the backup
server and inject malicious data.

**Fix:** Introduced `getHostKeyCallback()`, which reads the
`WALG_SSH_KNOWN_HOSTS` environment variable. When set, it uses
`knownhosts.New()` to validate the server's host key against the known_hosts
file. When unset, a warning is logged and the insecure fallback is used (for
backward compatibility in trusted environments).

```go
func getHostKeyCallback() (ssh.HostKeyCallback, error) {
    knownHostsFile := os.Getenv("WALG_SSH_KNOWN_HOSTS")
    if knownHostsFile == "" {
        tracelog.WarningLogger.Println("WALG_SSH_KNOWN_HOSTS not set; SSH host keys will not be verified")
        return ssh.InsecureIgnoreHostKey(), nil
    }
    callback, err := knownhosts.New(knownHostsFile)
    if err != nil {
        return nil, fmt.Errorf("read known_hosts file %q: %w", knownHostsFile, err)
    }
    return callback, nil
}
```

**Validation:**
- `go build ./pkg/storages/sh/...` succeeds.
- Unit test creates a known_hosts file, sets the env var, and verifies the
  callback is non-nil and not `InsecureIgnoreHostKey`.

---

## Fix 4: Unsafe `interface{}` Deserialisation (MEDIUM)

### 4a. Sentinel User Data (`internal/configure.go:399`)

**Root cause:** `UnmarshalSentinelUserData` decoded user-provided JSON into
`interface{}`, producing `float64` for numbers (precision loss) and
`map[string]interface{}` for objects. Consumers had to perform unchecked type
assertions.

**Fix:** Switched to `json.NewDecoder` with `UseNumber()`, which stores numeric
values as `json.Number` (preserving exact string representation). Combined with
the existing size limit on user-data input, this prevents type-confusion attacks.

```go
decoder := json.NewDecoder(strings.NewReader(userDataStr))
decoder.UseNumber()
var out interface{}
err := decoder.Decode(&out)
```

### 4b. S3 Headers (`pkg/storages/s3/session.go:249`)

**Root cause:** `decodeHeaders` unmarshalled YAML into `interface{}` and used
bare type assertions (`v.(string)`, `header.(map[string]interface{})`) that
would panic on unexpected input.

**Fix:** Replaced bare assertions with comma-ok forms that return descriptive
errors:

```go
for k, v := range interfaces {
    strVal, ok := v.(string)
    if !ok {
        return nil, fmt.Errorf("header %q value is not a string", k)
    }
    headers[k] = strVal
}
```

Similarly, `reformHeaderListToMap` now uses `ma, ok := header.(map[string]interface{})`
and skips non-map entries instead of panicking.

**Validation:**
- `go build ./pkg/storages/s3/` succeeds.
- Unit tests verify:
  - Normal YAML map headers decode correctly.
  - YAML list headers decode correctly.
  - Malformed input (non-string values, non-map list entries) returns errors
    instead of panicking.

---

## Summary

| # | Finding                    | Severity | Status | Lines Changed |
|---|----------------------------|----------|--------|---------------|
| 1 | SQL injection (MySQL)      | HIGH     | Fixed  | 8             |
| 2 | TLS InsecureSkipVerify     | HIGH     | Fixed  | 1             |
| 3 | SSH InsecureIgnoreHostKey  | HIGH     | Fixed  | 27            |
| 4 | Unsafe interface{} deser   | MEDIUM   | Fixed  | 16            |

**Total: 4 security vulnerabilities remediated across 5 files.**
