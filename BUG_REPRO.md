# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	emergency-claim-code/cmd/claimctl	0.002s
ok  	emergency-claim-code/internal/httpapi	0.010s
ok  	emergency-claim-code/internal/model	0.002s
?   	emergency-claim-code/internal/policy	[no test files]
ok  	emergency-claim-code/internal/query	0.002s
ok  	emergency-claim-code/internal/repository	0.018s
--- FAIL: TestBusiness01Regression (0.01s)
    service_test.go:35: record state conflict: approved to submitted
FAIL
FAIL	emergency-claim-code/internal/service	0.023s
ok  	emergency-claim-code/internal/store	0.018s
ok  	emergency-claim-code/internal/workflow	0.018s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/claimctl): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/claimctl): exit `0`
