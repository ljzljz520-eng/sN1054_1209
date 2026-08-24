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
?   	example.com/familyitinerary/cmd/itinerary	[no test files]
ok  	example.com/familyitinerary/internal/advisor	0.008s
--- FAIL: TestItineraryChatRetainsStatus (0.00s)
    chat_test.go:25: new message status=sent
FAIL
FAIL	example.com/familyitinerary/internal/chat	0.004s
ok  	example.com/familyitinerary/internal/config	0.001s
ok  	example.com/familyitinerary/internal/httpapi	0.004s
ok  	example.com/familyitinerary/internal/itinerary	0.002s
ok  	example.com/familyitinerary/internal/model	0.002s
ok  	example.com/familyitinerary/internal/report	0.004s
ok  	example.com/familyitinerary/internal/store	0.030s
ok  	example.com/familyitinerary/internal/validation	0.001s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/itinerary): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/itinerary): exit `0`
- Frontend build (web): exit `0`
