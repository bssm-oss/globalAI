# globalAI

`globalAI`는 한곳에서 검토하기 어려운 AI 프롬프트와 지침 파일을 모아서 보여주는 Go CLI입니다. 현재 제공되는 첫 번째 워크플로는 `globalai web`이며, `AGENTS.md`, `CLAUDE.md`, `.claude/`, `.cursor/rules/`, `.sisyphus/` 같은 프로젝트 및 전역 지침 파일을 로컬 뷰어에서 확인할 수 있게 해줍니다.

이 프로젝트는 의도적으로 의존성을 매우 적게 유지합니다. CLI는 명령 처리, 탐색, 로컬 HTTP 서버, 임베디드 정적 자산까지 모두 Go 표준 라이브러리만 사용합니다. 덕분에 바이너리를 감사하기 쉽고, 배포하기 쉽고, 이후 기여자가 구조를 이해하기도 쉽습니다.

## 왜 이 프로젝트가 필요한가

많은 AI 도구는 중요한 지침을 서로 다른 위치에 저장합니다. 저장소 파일, 사용자 전역 설정 폴더, 도구별 규칙 디렉터리가 대표적입니다. 도구 수가 늘어날수록 이런 파일은 추적하고 검토하기 어려워집니다. `globalAI`는 프론트엔드 빌드 파이프라인이나 네트워크 의존성 없이, 이런 소스들을 한곳에서 볼 수 있는 로컬 뷰어를 제공합니다.

## 현재 기능

- `globalai web`은 루프백 전용 로컬 뷰어를 시작합니다.
- 뷰어는 전체 파일 시스템을 재귀적으로 스캔하지 않고, 엄선된 allowlist 경로만 노출합니다.
- 프롬프트 본문은 브라우저에서 일반 텍스트로 렌더링되어, 원본 파일에 포함된 HTML을 신뢰하지 않습니다.
- UI 자산은 Go 바이너리 안에 직접 포함되어 배포됩니다.

## 탐색 모델

현재 allowlist는 광범위한 탐색보다 명확성과 안전성을 우선합니다.

`globalAI`는 디스크에 있는 모든 AI CLI를 무차별적으로 수집한다고 주장하지 않습니다. 현재 릴리스는 아래에 명시된 소스 계열만 지원하며, 그보다 넓은 범위는 의도적으로 포함하지 않습니다.

### 프로젝트 로컬 소스

- `AGENTS.md`
- `CLAUDE.md`
- `GEMINI.md`
- `.github/copilot-instructions.md`
- `.claude/**`
- `.cursor/rules/**`
- `.sisyphus/**`

### 전역 소스

- `~/AGENTS.md`
- `~/.claude/**`
- `~/.cursor/rules/**`
- `~/.sisyphus/**`

현재 지원하는 파일 확장자는 `.md`, `.mdc`, `.txt`, `.json`, `.yaml`, `.yml` 입니다. 1MiB보다 큰 파일은 건너뜁니다.

### 지원 소스 매트릭스

| 계열 | 위치 | 테스트 포함 여부 |
| --- | --- | --- |
| 프로젝트 에이전트 파일 | `AGENTS.md` | 예 |
| 프로젝트 Claude 파일 | `CLAUDE.md` | 예 |
| 프로젝트 Gemini 파일 | `GEMINI.md` | 예 |
| 프로젝트 Copilot 지침 | `.github/copilot-instructions.md` | 예 |
| 프로젝트 Claude 규칙 | `.claude/**` | 예 |
| 프로젝트 Cursor 규칙 | `.cursor/rules/**` | 예 |
| 프로젝트 Sisyphus 파일 | `.sisyphus/**` | 예 |
| 전역 에이전트 파일 | `~/AGENTS.md` | 예 |
| 전역 Claude 규칙 | `~/.claude/**` | 예 |
| 전역 Cursor 규칙 | `~/.cursor/rules/**` | 예 |
| 전역 Sisyphus 파일 | `~/.sisyphus/**` | 예 |

이 매트릭스 밖의 경로는 현재 범위에 포함되지 않으며, 필요하다면 추정으로 추가하지 말고 명시적으로 확장해야 합니다.

## 설치

### 30초 빠른 시작

가장 빠른 확인 경로는 아래 순서입니다.

```bash
go run github.com/bssm-oss/globalAI/cmd/globalai@latest install
globalai web --no-open
```

설치기가 현재 `PATH` 안에서 쓸 수 있는 사용자 디렉터리를 찾으면 `globalai` 명령이 바로 잡히고, 그렇지 않으면 설치 위치와 셸 재적용 방법을 출력합니다. 정상적으로 시작되면 `globalai viewer ready at http://127.0.0.1:PORT` 형식의 주소가 출력됩니다. 그 주소를 브라우저에 열면 현재 저장소와 홈 디렉터리의 allowlisted 소스를 바로 볼 수 있습니다.

### 처음 설정할 때 가장 쉬운 Go 경로

첫 설치에서는 아래 명령을 권장합니다.

```bash
go run github.com/bssm-oss/globalAI/cmd/globalai@latest install
```

이 명령은 실행 중인 `globalai` 바이너리를 사용자 bin 디렉터리에 복사하고, 이미 `PATH` 에 잡힌 디렉터리가 있으면 그 위치를 우선 사용합니다. 설치 후 `globalai is ready on PATH` 가 나오면 바로 `globalai web` 을 실행하면 됩니다.

### Go로 바로 설치

Go가 이미 설치되어 있고 설치 위치와 `PATH` 를 직접 관리하고 싶다면 아래 명령으로도 설치할 수 있습니다.

```bash
go install github.com/bssm-oss/globalAI/cmd/globalai@latest
```

Go는 바이너리를 `GOBIN` 에 설치하고, `GOBIN` 이 비어 있으면 `$(go env GOPATH)/bin` 에 설치합니다. 즉 설치가 성공해도 그 디렉터리가 현재 셸의 `PATH` 에 없으면 `globalai: command not found` 가 나올 수 있습니다.

가장 먼저 실제 설치 경로를 확인하려면 아래처럼 실행하면 됩니다.

```bash
BIN_DIR="$(go env GOBIN)"
[ -n "$BIN_DIR" ] || BIN_DIR="$(go env GOPATH)/bin"
printf '%s\n' "$BIN_DIR"
```

설치가 끝났다면 아래 명령으로 먼저 동작을 확인할 수 있습니다.

```bash
BIN_DIR="$(go env GOBIN)"
[ -n "$BIN_DIR" ] || BIN_DIR="$(go env GOPATH)/bin"
"$BIN_DIR/globalai" --help
"$BIN_DIR/globalai" web --help
```

위 절대 경로 실행이 성공하면 설치 자체는 정상입니다. 그다음 `globalai` 를 바로 쓰고 싶다면 해당 디렉터리를 `PATH` 에 추가하면 됩니다. 첫 설치에서 이 과정을 자동으로 덜 신경 쓰고 싶다면 위의 `go run ... install` 경로를 쓰는 편이 더 쉽습니다.

#### zsh (macOS 기본 셸)

기본 `GOPATH` 를 쓰는 경우 가장 단순한 설정은 아래와 같습니다.

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

직접 `GOBIN` 을 설정해 쓰고 있다면 그 경로를 대신 넣으면 됩니다.

```bash
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.zshrc
source ~/.zshrc
```

#### bash

```bash
echo 'export PATH="$(go env GOPATH)/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

#### fish

```fish
set -Ux fish_user_paths (go env GOPATH)/bin $fish_user_paths
```

셸 설정을 바꾼 뒤에는 새 터미널을 열거나 설정 파일을 다시 읽은 다음 아래처럼 확인합니다.

```bash
command -v globalai
globalai --help
```

### 설치 직후 `globalai` 가 안 보일 때

아래 순서대로 보면 됩니다.

1. `go install github.com/bssm-oss/globalAI/cmd/globalai@latest` 가 에러 없이 끝났는지 확인합니다.
2. `BIN_DIR="$(go env GOBIN)"; [ -n "$BIN_DIR" ] || BIN_DIR="$(go env GOPATH)/bin"` 로 설치 위치를 확인합니다.
3. `"$BIN_DIR/globalai" --help` 가 되면 설치는 정상이고, 문제는 `PATH` 설정입니다.
4. 해당 디렉터리를 셸의 `PATH` 에 추가한 뒤 `globalai --help` 를 다시 실행합니다.

### GitHub Releases에서 바이너리 다운로드

Go를 따로 설치하지 않고 쓰고 싶다면 GitHub Releases에서 운영체제에 맞는 압축 파일을 내려받아 압축을 풀고 `globalai` 바이너리를 실행하면 됩니다.

일반적인 확인 순서는 아래와 같습니다.

1. 운영체제에 맞는 압축 파일을 내려받아 풉니다.
2. 압축을 푼 디렉터리에서 macOS/Linux 는 `./globalai --help`, Windows 는 `./globalai.exe --help` 로 실행 가능 여부를 확인합니다.
3. 필요하면 바이너리를 `PATH` 에 잡힌 디렉터리로 옮긴 뒤 `globalai web` 을 실행합니다.

이 바이너리 자산은 유지보수자가 `v*` 형식의 버전 태그를 푸시할 때 GitHub Releases에 자동으로 게시됩니다.

- macOS: `globalai_darwin_amd64.tar.gz`, `globalai_darwin_arm64.tar.gz`
- Linux: `globalai_linux_amd64.tar.gz`, `globalai_linux_arm64.tar.gz`
- Windows: `globalai_windows_amd64.zip`, `globalai_windows_arm64.zip`

릴리스 페이지: `https://github.com/bssm-oss/globalAI/releases`

## 시작하기

### 요구 사항

- Go 1.25+ (`go install` 이나 로컬 빌드 사용 시)

### 로컬 빌드

```bash
go build ./cmd/globalai
```

### 뷰어 실행

```bash
./globalai web
```

설치된 바이너리를 사용 중이라면 아래처럼 실행하면 됩니다.

```bash
globalai web
```

특정 저장소 루트를 명시적으로 지정해서 실행할 수도 있습니다.

```bash
./globalai web --root /path/to/repository
./globalai web /path/to/repository
```

### 첫 실행에서 보게 되는 것

- `globalai web` 은 루프백 주소(`127.0.0.1`)에만 로컬 서버를 엽니다.
- 기본 주소가 `127.0.0.1:0` 이므로 실행할 때마다 사용 가능한 임시 포트를 자동으로 고릅니다.
- 기본 동작은 브라우저를 자동으로 여는 것이며, `--no-open` 을 주면 브라우저를 열지 않습니다.
- 뷰어가 준비되면 `globalai viewer ready at http://127.0.0.1:PORT` 형식의 로컬 URL 이 출력됩니다.
- 왼쪽에는 발견된 소스 목록이, 오른쪽에는 선택한 파일의 본문이 표시됩니다.

브라우저가 자동으로 열리지 않더라도 출력된 URL 이 기준입니다. 원격 세션, SSH, 헤드리스 환경에서는 `--no-open` 으로 실행한 뒤 출력된 주소를 직접 열면 됩니다.

뷰어 프로세스는 실행된 동안 계속 살아 있으며, 종료하려면 터미널에서 `Ctrl+C` 를 누르면 됩니다.

### 뷰어에서 확인할 수 있는 것

- 현재 저장소 루트에서 발견된 프로젝트 로컬 소스
- 현재 사용자 홈 디렉터리에서 발견된 전역 소스
- 선택한 파일의 계열, 상대 경로, 크기, 원문 텍스트

즉, 하나의 저장소를 열어도 프로젝트 파일만이 아니라 `~/AGENTS.md`, `~/.claude/**`, `~/.cursor/rules/**`, `~/.sisyphus/**` 같은 전역 위치도 함께 보일 수 있습니다.

### 아무 파일도 보이지 않을 때

allowlist 에 포함된 파일이 하나도 없으면 뷰어는 빈 목록 상태로 열립니다. 이 경우 아래 위치에 파일이 있는지 먼저 확인하면 됩니다.

- 프로젝트 루트: `AGENTS.md`, `CLAUDE.md`, `GEMINI.md`
- 프로젝트 설정 위치: `.github/copilot-instructions.md`, `.claude/**`, `.cursor/rules/**`, `.sisyphus/**`
- 전역 위치: `~/AGENTS.md`, `~/.claude/**`, `~/.cursor/rules/**`, `~/.sisyphus/**`

### 유용한 플래그

- `--addr`: 리슨 주소를 덮어씁니다. 기본값은 `127.0.0.1:0` 입니다.
- `--open`: 시작 후 브라우저를 강제로 엽니다.
- `--no-open`: 브라우저를 열지 않습니다. CI나 스모크 테스트에 유용합니다.

### 로컬 API

뷰어는 정적 UI와 함께 로컬 전용 JSON 엔드포인트도 함께 노출합니다.

```text
GET /api/sources
```

이 엔드포인트는 현재 수집된 소스 목록과 메타데이터를 반환하므로, 스크립트나 디버깅 용도로 사용할 수 있습니다.

예를 들어 브라우저를 열지 않고 상태만 확인하고 싶다면 아래처럼 사용할 수 있습니다.

```bash
globalai web --no-open
# 다른 터미널에서
curl http://127.0.0.1:PORT/api/sources
```

## 개발 워크플로

### 포맷, 테스트, 빌드

```bash
gofmt -w $(find . -name '*.go' -not -path './.git/*')
go test ./...
go build ./cmd/globalai
```

### 기능 스모크 테스트

```bash
bash scripts/functional_smoke.sh
```

이 스크립트는 바이너리를 빌드하고, `globalai web --no-open`을 시작한 뒤, 로컬 API를 조회해 반환된 JSON을 검증합니다.

저장소에는 현재 지원하는 모든 소스 계열, zero-source 응답, 심볼릭 링크 거부, 대용량 파일 거부를 검증하는 단위 테스트도 포함되어 있습니다.

## CI

이 저장소는 다음 검증을 수행하는 GitHub Actions 워크플로를 포함합니다.

- `gofmt -s` 기반 포맷 검증
- `go vet`
- `go test -race ./...`
- `go build ./cmd/globalai`
- 기능 스모크 테스트 스크립트
- 설치 스모크 테스트 스크립트

또한 버전 태그(`v*`)를 푸시하면 macOS, Linux, Windows용 압축 바이너리를 GitHub Release 자산으로 자동 게시합니다.

최근 GitHub Actions 검증 결과:

- PR 실행: PR #1에서 `ci` 성공
- 메인 브랜치 실행: `main` 머지 후 `ci` 성공

## 프로젝트 구조

```text
cmd/globalai/           CLI 진입점
internal/cli/           명령 파싱과 오케스트레이션
internal/browser/       OS 브라우저 실행기
internal/source/        allowlist 기반 프롬프트 소스 탐색
internal/viewer/        HTTP 서버와 임베디드 정적 자산
docs/                   아키텍처 메모와 변경 기록
scripts/                반복 가능한 검증 도구
```

## 기여자와 AI 에이전트를 위한 문서

- `AGENTS.md`는 저장소 규칙, 아키텍처 경계, 이후 자동화된 변경 시 기대사항을 설명합니다.
- `docs/architecture.md`는 현재 시스템 구조를 설명합니다.
- `docs/changes/2026-04-06-initial-bootstrap.md`는 첫 부트스트랩 변경을 장기 보관 가능한 형식으로 기록합니다.

## 로드맵

- 더 많은 AI 도구 소스 계열 지원
- 뷰어 내부 검색과 필터링
- 프롬프트 상태를 감사하기 위한 더 나은 내보내기 또는 스냅샷 워크플로
