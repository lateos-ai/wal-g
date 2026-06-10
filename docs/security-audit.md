# Security & Stability Audit Report

**Branch:** `security-stability-audit`  
**Date:** 2026-06-09  
**Tools:** govulncheck, semgrep, golangci-lint, go-licenses, go list -u -m all, manual grep

---

## 1. Dependency Vulnerability Scan (govulncheck)

**Command:** `govulncheck -scan=package ./...`

### Remaining Vulnerabilities (3 total, all N/A fix)

| ID | Module | Found In | Fixed In | Severity |
|---|---|---|---|---|
| GO-2026-4887 | github.com/docker/docker | v28.5.2+incompatible | N/A | HIGH |
| GO-2026-4883 | github.com/docker/docker | v28.5.2+incompatible | N/A | HIGH |
| GO-2026-4518 | github.com/jackc/pgproto3/v2 | v2.3.3 | N/A | MEDIUM |

**Action:** None. All earlier vulnerabilities (x/net, x/crypto, stdlib) were fixed in PR #8. These 3 have no available fix yet.

---

## 2. SAST Scan (Semgrep)

**Command:** `semgrep --config "p/golang" --config "p/security-audit" .`  
**Rules run:** 143 | **Findings:** 37

### CRITICAL (0)

### HIGH (5)

| Finding | File | Line | Issue |
|---|---|---|---|
| Unsafe deserialization into `interface{}` | `internal/configure.go` | 399-400 | JSON unmarshal into `interface{}` allows arbitrary types (CWE-502) |
| Unsafe deserialization into `interface{}` | `pkg/storages/s3/session.go` | 249-250 | YAML unmarshal into `interface{}` (CWE-502) |
| SHA1 hash for crypto | `internal/crypto/envelope/enveloper.go` | 34 | `sha1.Sum()` — not collision-resistant |
| SQL injection (string-formatted query) | `internal/databases/mysql/mysql.go` | 37 | `"SELECT @@" + variable` |
| SQL injection (string-formatted queries) | `internal/databases/sqlserver/*.go` | Multiple | `fmt.Sprintf("BACKUP DATABASE %s TO %s", ...)` — 7 occurrences across backup/restore/log handlers |

### MEDIUM (16)

| Finding | File | Line | Issue |
|---|---|---|---|
| `math/rand` used (non-production) | `internal/storagetools/check.go` | 6 | Should use `crypto/rand` |
| `math/rand` used | `internal/multistorage/stats/alive_checker.go` | 8 | Should use `crypto/rand` |
| `math/rand` used (non-production) | `internal/profile.go` | 4 | Should use `crypto/rand` |
| `math/rand` used (test files) | 9 test files | Multiple | Tests only, lower priority |
| MD5 hash | `pkg/storages/s3/folder.go` | 58 | `md5.Sum([]byte(sseCustomerKey))` |
| MD5 hash | `pkg/storages/storage/storage.go` | 33 | `md5.New()` |
| MD5 hash (test) | `pkg/storages/s3/uploader_test.go` | 55 | Test only |
| Unsafe pointer in Windows code | `internal/multistorage/stats/cache/flock_windows.go` | 27-36 | `unsafe.Pointer` for syscall args (needed for Windows API) |
| Unsafe pointer in Windows code | `internal/diskwatcher/disk_watcher_windows.go` | 27-30 | `unsafe.Pointer` for syscall args (needed for Windows API) |
| Unsafe pointer for page verification | `internal/databases/postgres/paged_file_verifier.go` | 95 | `unsafe.Pointer` to cast page bytes (needed for low-level PG page checksum) |

### LOW (16)

| Finding | File | Line | Issue |
|---|---|---|---|
| Dockerfile ends with `USER root` | `docker/cloudberry_tests/Dockerfile` | 23 | Container runs as root |
| Dockerfile ends with `USER root` | `docker/etcd_tests/Dockerfile` | 39 | Container runs as root |
| Dockerfile ends with `USER root` | `docker/gp_tests/Dockerfile` | 20 | Container runs as root |
| TLS MinVersion not set | `internal/databases/mysql/mysql.go` | 191-193 | TLS 1.2 default, should pin TLS 1.3 |
| TLS MinVersion not set + InsecureSkipVerify | `internal/databases/sqlserver/backup_import_handler.go` | 125 | Should pin TLS 1.3 |
| Bind to all interfaces | `internal/databases/mongo/binary/mongod_runner.go` | 114 | `net.Listen("tcp", ":0")` |
| SSH InsecureIgnoreHostKey | `pkg/storages/sh/storage.go` | 57 | Host key verification disabled (MITM risk) |
| `math/rand` in test/benchmark files | 12 test files | Multiple | Low priority |

---

## 3. Anti-Slop Code Audit (golangci-lint)

**Command:** `golangci-lint run ./...`  
**Linters used:** gci, gofmt, goimports, govet, errcheck, ineffassign, misspell, revive, staticcheck, unconvert, whitespace, gocritic, gocyclo, dupl, funlen, lll, nakedret, unparam, unused, bodyclose, copyloopvar, asciicheck, makezero

### Findings (7)

| Severity | Linter | File | Line | Issue |
|---|---|---|---|---|
| MEDIUM | gci | `cmd/common/common.go` | 1 | Import not properly formatted |
| MEDIUM | gci | `internal/databases/redis/archive/sharded.go` | 1 | Import not properly formatted |
| MEDIUM | gci | `pkg/storages/gcs/folder.go` | 1 | Import not properly formatted |
| MEDIUM | staticcheck | `internal/multistorage/stats/cache/flock_windows.go` | 28 | `syscall.Syscall6` deprecated, use `SyscallN` |
| MEDIUM | unconvert | `internal/multistorage/stats/cache/flock_windows.go` | 29 | Unnecessary `uintptr()` conversion |
| LOW | whitespace | `internal/diskwatcher/disk_watcher_windows.go` | 11 | Leading newline |
| LOW | whitespace | `internal/diskwatcher/disk_watcher_windows.go` | 47 | Trailing newline |

---

## 4. Cryptographic Audit

### HMAC
- **Not used** — no direct HMAC operations found.
- Envelope encryption uses key-wrapping (AWS KMS, YC KMS, OpenPGP).

### AES
- **Not used directly** — no `aes.NewCipher` calls in codebase.
- Encryption is delegated to external providers (AWS KMS, YC KMS, OpenPGP).

### Random Number Generation
- `crypto/rand` used in: `internal/fsutil/direct_io_reader_test.go`, `internal/crypto/awskms/key.go`
- `math/rand` used in: 12 non-production/test files, 3 production files (`internal/storagetools/check.go`, `internal/multistorage/stats/alive_checker.go`, `internal/profile.go`)

### SHA1 Usage (1 finding)
- `internal/crypto/envelope/enveloper.go:34` — `sha1.Sum(encryptedKey.Data)` for key fingerprinting. Used only for identification/logging, not for cryptographic verification. Low severity.

### MD5 Usage (2 production + 1 test)
- `pkg/storages/s3/folder.go:58` — `md5.Sum([]byte(sseCustomerKey))` — used for SSE-C key hashing by AWS S3 SDK requirement. Not used for security.
- `pkg/storages/storage/storage.go:33` — `md5.New()` — used for ETag/checksum computation. Lower severity as it's not used for authentication.
- `pkg/storages/s3/uploader_test.go:55` — Test only.

### TLS/SSL
- `internal/databases/mysql/mysql.go:191` — TLS config missing `MinVersion` (defaults to TLS 1.2; should pin to 1.3)
- `internal/databases/sqlserver/backup_import_handler.go:125` — `InsecureSkipVerify: true` + missing `MinVersion`
- `pkg/storages/sh/storage.go:57` — `ssh.InsecureIgnoreHostKey()` — disables host key verification for SFTP/SSH storage

### Key Storage
- AWS KMS keys: configured via environment variables (standard)
- YC KMS keys: configured via environment variables (standard)
- OpenPGP keys: via environment variable or file (standard)
- No hardcoded keys found

### Password Handling
- No plaintext password logging found

---

## 5. License Compliance

**Command:** `go-licenses csv ./utility ./internal/crypto ./internal/checksum`

| Dependency | License | Compatible |
|---|---|---|
| github.com/lateos-ai/wal-g | MIT | ✅ |
| github.com/pkg/errors | BSD-2-Clause | ✅ |
| github.com/wal-g/tracelog | Apache-2.0 | ✅ |

**No GPL/AGPL licenses found.** Full scan failed on `encoding/json/v2` (experimental stdlib), but sampled packages show clean licenses.

---

## 6. Dependency Age & Maintenance

**Command:** `go list -u -m all`  
**Total outdated dependencies:** ~170

### Key production dependencies 1+ year stale

| Dependency | Current Version | Latest | Age |
|---|---|---|---|
| `cloud.google.com/go` | v0.65.0 | v0.123.0 | ~5 years |
| `google.golang.org/api` | v0.30.0 | v0.283.0 | ~5 years |
| `google.golang.org/genproto` | v0.0.0-20211021150943 | v0.0.0-202606 | ~5 years |
| `gopkg.in/ini.v1` | v1.67.0 | v1.67.3 | Ok (minor) |
| `github.com/prometheus/client_golang` | v1.12.1 | v1.23.2 | ~3 years |
| `github.com/spf13/cobra` | v1.7.0 | v1.10.2 | ~2 years |
| `github.com/aws/aws-sdk-go` | v1.55.7 | v1.55.8 | Ok (recent) |
| `github.com/Azure/azure-sdk-for-go/sdk/azcore` | v1.21.1 | v1.22.0 | Ok (recent) |
| `github.com/klauspost/compress` | v1.18.5 | v1.18.6 | Ok (recent) |

### Abandoned or low-maintenance dependencies
- `github.com/3rf/mongo-lint` — last updated 2014 (fork of golint)
- `github.com/ryanuber/columnize` — last updated 2016
- `github.com/bgentry/speakeasy` — last updated 2015
- `github.com/jstemmer/go-junit-report` — fork, last updated 2016
- `github.com/kr/logfmt` — last updated 2017
- `github.com/pascaldekloe/goe` — last updated 2019

Most of these are transitive dev dependencies and don't affect production.

---

## Summary

| Category | CRITICAL | HIGH | MEDIUM | LOW |
|---|---|---|---|---|
| govulncheck | 0 | 2 | 1 | 0 |
| Semgrep SAST | 0 | 5 | 16 | 16 |
| golangci-lint | 0 | 0 | 5 | 2 |
| Crypto audit | 0 | 0 | 4 | 1 |
| **Total** | **0** | **7** | **26** | **19** |

**Top priorities for Part 2 (Manual Review):**
1. SQL injection in `mysql.go` and `sqlserver/*.go` — verify `quoteName()` is safe
2. TLS `InsecureSkipVerify` in SQL Server backup import
3. SSH `InsecureIgnoreHostKey` in SFTP storage
4. Unsafe deserialization in `configure.go` and `s3/session.go`
5. `ssh.InsecureIgnoreHostKey` in SSH storage
