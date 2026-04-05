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

## 시작하기

### 요구 사항

- Go 1.25+

### 빌드

```bash
go build ./cmd/globalai
```

### 뷰어 실행

```bash
./globalai web
```

특정 저장소 루트를 명시적으로 지정해서 실행할 수도 있습니다.

```bash
./globalai web --root /path/to/repository
./globalai web /path/to/repository
```

### 유용한 플래그

- `--addr`: 리슨 주소를 덮어씁니다. 기본값은 `127.0.0.1:0` 입니다.
- `--open`: 시작 후 브라우저를 강제로 엽니다.
- `--no-open`: 브라우저를 열지 않습니다. CI나 스모크 테스트에 유용합니다.

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
- macOS, Linux, Windows용 패키지 릴리스
- 프롬프트 상태를 감사하기 위한 더 나은 내보내기 또는 스냅샷 워크플로
