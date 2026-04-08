# 2026-04-08 — Go install 경로 안내와 설치 스모크 검증 추가

## 요약

이번 변경은 `go install github.com/bssm-oss/globalAI/cmd/globalai@latest` 이후 `globalai` 명령이 바로 잡히지 않을 수 있는 실제 셸 환경 문제를 README에서 명확히 설명하고, 이를 검증하는 설치 스모크 테스트를 추가한 작업입니다.

## 무엇이 달라졌는가

- README의 빠른 시작과 설치 섹션을 절대 경로 기반 검증 흐름으로 바꿨습니다.
- `GOBIN` 또는 `$(go env GOPATH)/bin` 이 `PATH` 에 없을 때 왜 `command not found` 가 나는지 설명을 추가했습니다.
- zsh, bash, fish 에서 Go 바이너리 디렉터리를 `PATH` 에 넣는 예시를 추가했습니다.
- `scripts/install_smoke.sh` 를 추가해 설치 직후 절대 경로 실행은 성공하고, `PATH` 추가 전에는 bare command 가 실패하는 흐름을 검증하도록 했습니다.
- CI와 기여 가이드에 새 설치 스모크 테스트를 포함했습니다.

## 왜 필요한가

- `go install` 자체는 성공해도 셸 `PATH` 설정이 빠져 있으면 사용자는 도구가 설치되지 않았다고 느끼기 쉽습니다.
- 기존 README는 이 점을 짧게만 언급했고, 첫 실행 예제가 너무 빨리 `globalai ...` 호출로 넘어가 실제 실패 사례를 막기에 부족했습니다.

## 검증 근거

- `go test ./...`
- `go build ./cmd/globalai`
- `bash scripts/functional_smoke.sh`
- `bash scripts/install_smoke.sh`
