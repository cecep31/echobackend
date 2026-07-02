# Security Audit Report — echobackend

| Field | Value |
|-------|-------|
| **Project** | echobackend |
| **Commit audited** | `15f2c32` (main) |
| **Date** | 2026-07-02 |
| **Scope** | Full backend codebase: auth, JWT, sessions, passwords, OAuth, cookies, rate limiting, SQL, file upload, secrets handling, transport security, logging, and all `internal/platform/*` adapters. |
| **Out of scope** | Frontend (not in this repo); infrastructure provisioning; deployed instance runtime. |

---

## 1. Executive Summary

No critical vulnerabilities and **no known dependency CVEs** (`govulncheck` clean).
The codebase shows solid security fundamentals: JWT signing-method enforcement,
hashed refresh tokens, bcrypt password hashing, constant-time OAuth state
comparison, parameterized SQL, server-side file-type validation, and DB-backed
admin authorization (the `is_super_admin` JWT claim is not trusted).

**2 Medium** and **9 Low/Informational** findings were identified. The two
Mediums — (1) `ChangePassword` not invalidating sessions and (2) a weak/unused
password policy — are the most worthwhile to fix and are both low-risk code
changes.

### Severity counts

| Severity | Count |
|----------|-------|
| Critical | 0 |
| High     | 0 |
| Medium   | 2 |
| Low      | 6 |
| Informational | 3 |

---

## 2. Methodology

1. **Static analysis tools**
   - `govulncheck` (Go vulnerability database) → **No vulnerabilities found**.
   - `go vet ./...` → clean.
   - `go build ./...` → clean.
   - `go test ./...` → all passing.
   - `gosec` was listed in AGENTS.md but not installed in the environment; manual
     review covered the categories it would flag (hardcoded creds, weak crypto,
     SQL injection, error leakage).
2. **Manual code review** of all security-critical paths:
   `internal/service/auth_service.go`, `internal/middleware/auth_middleware.go`,
   `internal/handler/auth_handler.go`, `internal/routes/authRoutes.go`,
   `internal/platform/{database,storage,email,cache,queue}/*.go`,
   `config/config.go`, `pkg/response`, `pkg/market/*`, repository SQL.
3. **Framework semantics verification**: confirmed GORM `clause.Expression`
   interface requirements and Echo v5 middleware/routing semantics against
   vendored source where relevant.

---

## 3. Security Strengths (things done well)

| Area | Evidence |
|------|----------|
| Dependency vulnerabilities | `govulncheck` clean — no known CVEs in the dependency graph. |
| JWT signing method | HMAC method enforced in `validateToken` (`auth_middleware.go:124`); blocks `alg=none` confusion attacks. |
| JWT secret strength | Minimum 32 characters enforced (`config.go:272`); app panics if missing. |
| JWT expiry | `exp` claim set on issue (`auth_service.go:513`) and validated by the library. |
| Refresh token storage | Stored as **SHA-256 hash** in DB (`auth_service.go:529`), not plaintext; rotated on each refresh. |
| Password hashing | bcrypt `DefaultCost` used for register & reset & change (`auth_service.go:131,319`). |
| Login failure handling | Returns generic `ErrInvalidCredentials` for unknown user, nil password, and wrong password; logs failed-login activity. |
| OAuth state (CSRF) | `subtle.ConstantTimeCompare` (constant-time), HttpOnly + SameSite=Lax + conditional `Secure` cookie, **one-time use** (cleared before validation), path-scoped (`auth_handler.go:334-359`). |
| Auth rate limiting | Per-IP fixed-window on login/register/forgot/reset/refresh/exchange (`authRoutes.go:12-17`) with Redis-backed shared store and in-memory fallback. |
| User enumeration (forgot-password) | Returns a constant success message regardless of whether the email exists (`auth_handler.go:119`). |
| Error leakage | `response.InternalServerError` does **not** echo `err.Error()` to the client (`response.go:140-151`). |
| Admin authorization | `isSuperAdminFromDB` re-checks from the database; the JWT `is_super_admin` claim is **not** trusted for authorization decisions (`auth_middleware.go:98-104`). |
| SQL injection | All queries parameterized via GORM placeholders (`Where("col = ?", v)`); `ORDER BY` built only from whitelists (`post_repository.go:425`, `holding_repository.go:86-99`). |
| File upload | Server-side magic-byte content detection (`http.DetectContentType`), 1 MB size cap enforced twice (header + `LimitReader`), random hex object keys (no predictable paths), JPEG/PNG/WebP allowlist (`post_service.go:408-426`). |
| Secrets management | `.env` is gitignored; all secrets sourced from environment; no `fmt.Println`/`os.WriteFile` debug-leak patterns found in the codebase. |
| HTTP client hygiene | All external HTTP responses use `defer resp.Body.Close()`; timeouts configured on every client (Yahoo, RapidAPI, OpenRouter, GitHub). |
| Transport defaults | S3 `UseSSL=true` default (`config.go:211`); SMTP uses opportunistic STARTTLS with TLS 1.2 minimum and Go's `smtp.PlainAuth` localhost-plaintext guard (`email.go:147-153`). |

---

## 4. Findings

> Each finding lists severity, location, description, impact, and a recommended
> fix. "Code-actionable" fixes can be applied directly; "Operational" findings
> require deployment/environment configuration.

### 4.1 [MEDIUM] `ChangePassword` does not invalidate existing sessions

| | |
|---|---|
| **Location** | `internal/service/auth_service.go:304-333` |
| **Severity** | Medium |
| **Code-actionable** | Yes |

**Description**
After a successful password change, existing access tokens (valid up to
**3 hours**, `config.go:197`) and refresh tokens (30 days, `config.go:198`)
**remain valid**. `ResetPassword` correctly calls
`s.sessionRepo.DeleteByUserID(ctx, user.ID)` at line 268 to revoke all sessions,
but `ChangePassword` does not perform any session revocation.

**Impact**
If an account is compromised and the legitimate user changes their password to
regain control, the attacker's previously-issued tokens keep working for up to
3 hours (access token) / 30 days (refresh token). This defeats the user's
remediation step.

**Recommended fix**
After the successful `s.userRepo.Update(ctx, user)` in `ChangePassword`, add:
```go
_ = s.sessionRepo.DeleteByUserID(ctx, userID)
```
(mirroring `ResetPassword`). The access token cannot be revoked server-side
(stateless JWT), but its 3-hour window is short; revoking refresh tokens
prevents long-lived access. Optionally shorten `JWT_EXPIRY_HOURS` if a faster
access-token invalidation guarantee is required.

---

### 4.2 [MEDIUM] Weak and unused password policy

| | |
|---|---|
| **Location** | `internal/dto/auth.go:5,11,20,29`; `internal/apperror/errors.go:53-58` |
| **Severity** | Medium |
| **Code-actionable** | Yes |

**Description**
Password enforcement is only `validate:"required,min=8"` — so passwords like
`aaaaaaaa` or `password` are accepted. Meanwhile, the complexity sentinel
errors are **defined but never used anywhere in the codebase** (dead code):
`ErrPasswordTooShort`, `ErrPasswordTooLong`, `ErrPasswordNoUpper`,
`ErrPasswordNoLower`, `ErrPasswordNoDigit`, `ErrPasswordNoSpecial`.

Additionally, no maximum length is enforced at the DTO layer, so bcrypt
silently truncates passwords longer than 72 bytes (see finding 4.4).

**Impact**
Accounts can be created / reset with trivially weak passwords, increasing
brute-force susceptibility (mitigated somewhat by rate limiting, but not for
stolen-hash offline attacks).

**Recommended fix**
1. Implement a `validatePassword(pw string) error` function that enforces length
   (8–72) and complexity (upper, lower, digit, special) using the existing
   `apperror` sentinels, and call it in `Register`, `ResetPassword`, and
   `ChangePassword`.
2. Add `max=72` to the password DTO tags (see 4.4).

---

### 4.3 [LOW-MEDIUM] Insecure default S3 credentials baked into config

| | |
|---|---|
| **Location** | `config/config.go:208-209` |
| **Severity** | Low-Medium |
| **Code-actionable** | Yes |

**Description**
```go
AccessKey: envString([]string{"S3_ACCESS_KEY", "MINIO_ACCESS_KEY"}, "minioadmin"),
SecretKey: envString([]string{"S3_SECRET_KEY", "MINIO_SECRET_KEY"}, "minioadmin"),
```
If an operator sets `S3_ENDPOINT`/`S3_BUCKET` but forgets the credentials, the
app silently authenticates to S3/MinIO using the well-known defaults
`minioadmin/minioadmin`.

**Impact**
A misconfigured production instance could store user uploads in a bucket
accessible via widely-known default credentials, or could fail-open in a way
that looks "working" but is insecure.

**Recommended fix**
Default the credentials to `""` and extend `NewS3Storage` to disable storage
when either key is empty (it already returns `nil` when endpoint/bucket are
empty — add the same guard for the keys).

---

### 4.4 [LOW] bcrypt 72-byte silent truncation

| | |
|---|---|
| **Location** | `internal/dto/auth.go` (no `max` tag on password fields) |
| **Severity** | Low |
| **Code-actionable** | Yes |

**Description**
Passwords longer than 72 bytes are silently truncated by bcrypt with no warning.
Two distinct long passwords sharing the same first 72 bytes would authenticate
identically, a surprising and hard-to-detect behavior.

**Impact**
Low — only affects users choosing very long passphrases (>72 bytes). No security
breach, but unexpected equivalence.

**Recommended fix**
Add `max=72` to the password `validate` tags, or enforce the existing
`ErrPasswordTooLong` in the new `validatePassword` function (prefer `max=72` to
match bcrypt's real limit, or explicitly pre-hash/truncate with SHA-256 if long
passphrases must be supported).

---

### 4.5 [LOW] OAuth redirect/callback routes have no rate limiting

| | |
|---|---|
| **Location** | `internal/routes/authRoutes.go:30-31` |
| **Severity** | Low |
| **Code-actionable** | Yes |

**Description**
`/api/auth/oauth/github` (redirect) and `/api/auth/oauth/github/callback` lack
the `FixedWindowRateLimiter` that is applied to every other auth route
(login, register, forgot-password, reset-password, refresh, oauth-exchange).

**Impact**
Unbounded state-cookie issuance and callback processing per IP — enables mild
DoS / log flooding / cookie-set amplification.

**Recommended fix**
Add per-IP limiters, e.g.:
```go
oauthRedirectRateLimit := appmiddleware.FixedWindowRateLimiterWithCache(r.cache, "auth:oauth-redirect", 10, time.Minute)
auth.GET("/oauth/github", r.authHandler.GithubOAuthRedirect, oauthRedirectRateLimit)
auth.GET("/oauth/github/callback", r.authHandler.GithubOAuthCallback, oauthRedirectRateLimit)
```

---

### 4.6 [LOW] GitHub placeholder email for users without a public email

| | |
|---|---|
| **Location** | `internal/service/auth_service.go:407` |
| **Severity** | Low |
| **Code-actionable** | Yes |

**Description**
When a GitHub user has no public email, the service creates an account with
`fmt.Sprintf("%d@github.placeholder", githubUser.ID)` as the email.

**Impact**
These accounts are passwordless (GitHub-only) with a predictable, non-deliverable
placeholder email. The placeholder domain could be enumerated and may cause odd
`FindUserByEmail` matches if a real user later registers that exact string.

**Recommended fix**
Require a verified email before account creation (return an error and redirect
the user to "add a public email on GitHub"), or mark these accounts as
email-unverified and block email-dependent flows until a real email is set.

---

### 4.7 [LOW] CORS defaults to `*`

| | |
|---|---|
| **Location** | `config/config.go:193` → `parseOrigins(envString(..., "*"))` |
| **Severity** | Low |
| **Code-actionable** | Operational (config) |

**Description**
`HTTP_ALLOW_ORIGINS` defaults to `*` (allow all origins).

**Impact**
Low, because authentication is Bearer-token (not cookies), so cross-origin sites
cannot attach credentials. However, a wildcard policy is broader than necessary
and should be tightened in production.

**Recommended fix**
Set `HTTP_ALLOW_ORIGINS` to an explicit comma-separated list of frontend origins
in the production `.env` / Fly.io secrets.

---

### 4.8 [LOW] Database / Redis transport TLS not enforced in code

| | |
|---|---|
| **Location** | `internal/platform/database/setup.go:35` (DSN-driven `sslmode`); `internal/platform/cache/redis.go:42` (`redis://` vs `rediss://`) |
| **Severity** | Low |
| **Code-actionable** | Operational |

**Description**
There is no code-level enforcement that the PostgreSQL DSN uses
`sslmode=require` or that the Redis/Valkey URL uses `rediss://` (TLS). Security
depends entirely on the operator-supplied connection strings.

**Impact**
If a plaintext DSN/URL is supplied, credentials and query data travel
unencrypted over the network.

**Recommended fix**
Require `sslmode=require` (or verify-full) in the production `DATABASE_URL` and
`rediss://` for `REDIS_URL`/`QUEUE_REDIS_URL`. Optionally add a config-time
warning when `sslmode` is absent or `disable`.

---

### 4.9 [INFORMATIONAL] Tokens returned in JSON body (XSS-sensitive storage)

| | |
|---|---|
| **Location** | `internal/handler/auth_handler.go:90-98,176-184` (access + refresh tokens in JSON responses) |
| **Severity** | Informational |
| **Code-actionable** | Architectural (frontend) |

**Description**
This is a stateless-JWT design: both access and refresh tokens are returned in
the JSON body for the frontend to store. The backend cannot control how the
frontend persists them.

**Impact**
If the frontend stores tokens in `localStorage`, an XSS payload can exfiltrate
both tokens. `httpOnly` cookies would mitigate XSS theft but require CSRF
handling (a different tradeoff).

**Recommended fix**
Frontend responsibility (out of scope for this repo). If desired, the backend
could be refactored to set tokens as `httpOnly`, `Secure`, `SameSite` cookies,
and add CSRF protection. Document the expected storage strategy for frontend
implementers.

---

### 4.10 [INFORMATIONAL] GORM debug mode logs SQL with inlined parameters

| | |
|---|---|
| **Location** | `internal/platform/database/setup.go:31` → `ParameterizedQueries: !config.App.Debug` |
| **Severity** | Informational |
| **Code-actionable** | Operational |

**Description**
When `APP_DEBUG=true`, GORM logs full SQL **with inlined values** (which may
include password hashes, emails, reset tokens) at Info level. Debug is opt-in
(default `false`).

**Impact**
If debug is enabled in production, PII/secrets leak into application logs.

**Recommended fix**
Keep `APP_DEBUG=false` in production. Optionally, always set
`ParameterizedQueries: true` and lower the GORM log level to `Warn`/`Error`
regardless of debug, logging query shape rather than values.

---

### 4.11 [INFORMATIONAL] `TrustProxy` + XFF spoofing bypasses per-IP rate limits

| | |
|---|---|
| **Location** | `config/config.go:192` (default `false`); `middleware/setup.go` & `rate_limit.go` use `c.RealIP()` |
| **Severity** | Informational |
| **Code-actionable** | Operational |

**Description**
If `HTTP_TRUST_PROXY=true` is enabled without a trusted reverse proxy that
strips/overwrites the client-controlled `X-Forwarded-For` header, an attacker
can rotate spoofed IPs to bypass all per-IP rate limiting.

**Impact**
Mitigated by the default (`false`). Only relevant if explicitly enabled in a
misconfigured deployment.

**Recommended fix**
Only enable `HTTP_TRUST_PROXY` behind a trusted proxy that sets a single,
authoritative `X-Forwarded-For` value. Document this requirement near the
config key.

---

## 5. Recommended Fix Priority

| Priority | Finding | Effort |
|----------|---------|--------|
| 1 | 4.1 — Session invalidation in `ChangePassword` | Trivial (1 line) |
| 2 | 4.2 — Real password policy (`validatePassword`) | Small (new helper + 3 call sites) |
| 3 | 4.3 — Empty S3 credential defaults | Small |
| 4 | 4.4 — `max=72` password tag | Trivial |
| 5 | 4.5 — Rate limit OAuth routes | Trivial |
| 6 | 4.6 — Require verified email for OAuth | Small |
| 7 | 4.7–4.11 — Operational hardening | Config / deploy |

---

## 6. Appendix — Tool Output

### govulncheck
```
$ go run golang.org/x/vuln/cmd/govulncheck@latest ./...
No vulnerabilities found.
```

### go vet / build / test
```
$ go vet ./...   → exit 0 (clean)
$ go build ./... → exit 0 (clean)
$ go test ./...  → all packages pass
```

### gosec
Not installed in the audit environment. Manual review covered the categories
gosec would flag (hardcoded credentials, weak crypto, SQL injection, error
leakage, file inclusion). The single hardcoded-credentials concern found
(finding 4.3: `minioadmin` defaults) is documented above.

