# 2026-04-08 — Go 설치 부트스트랩 경로 추가

## 요약

이번 변경은 plain `go install` 뒤에 사용자가 직접 `PATH` 를 맞춰야 하는 마찰을 줄이기 위해, `go run github.com/bssm-oss/globalAI/cmd/globalai@latest install` 형태의 Go 기반 설치 부트스트랩 경로를 추가한 작업입니다.

## 무엇이 달라졌는가

- `globalai install` 서브커맨드를 추가했습니다.
- 이 설치기는 현재 실행 중인 `globalai` 바이너리를 사용자 bin 디렉터리에 복사합니다.
- 이미 `PATH` 에 잡힌 사용자 쓰기 가능 디렉터리가 있으면 그 위치를 우선 사용합니다.
- 그렇지 않으면 Go bin 디렉터리로 설치하고, 인식 가능한 셸(zsh, bash, fish)에 대해서는 설정 파일에 PATH 구문을 추가합니다.
- `scripts/install_smoke.sh` 에 `go run ./cmd/globalai install` 시나리오를 추가했습니다.
- README는 첫 설치 기본 경로를 `go run ... install` 로 올리고, 수동 `go install` 은 고급 경로로 유지했습니다.

## 왜 필요한가

- plain `go install` 은 바이너리를 설치해도 사용자의 현재 셸 PATH 를 바꾸지 못합니다.
- 그 결과 처음 쓰는 사용자는 설치가 실패했다고 느끼기 쉽습니다.
- 설치 부트스트랩은 이 간극을 줄여, 가능한 경우 바로 `globalai` 를 쓰게 하고 그렇지 않은 경우에도 정확한 다음 단계를 출력합니다.

## 검증 근거

- `go test ./...`
- `go build ./cmd/globalai`
- `bash scripts/functional_smoke.sh`
- `bash scripts/install_smoke.sh`
