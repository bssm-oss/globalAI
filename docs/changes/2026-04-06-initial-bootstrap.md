# 2026-04-06 — 초기 부트스트랩

## 요약

이번 변경은 비어 있던 `globalAI` 저장소를, 로컬 웹 뷰어·테스트·문서·CI 자동화를 갖춘 동작하는 Go CLI 프로젝트로 부트스트랩한 작업입니다.

## 무엇이 배포되었는가

- `web` 명령을 포함한 `globalai` Go 바이너리
- 프로젝트 및 전역 AI 지침 소스를 위한 결정적인 탐색 로직
- `/api/sources` 엔드포인트를 포함한 루프백 기반 로컬 임베디드 뷰어
- CLI 동작, 탐색, 브라우저 실행, 뷰어 서빙에 대한 단위 테스트
- 빌드된 바이너리를 끝까지 검증하는 기능 스모크 테스트 스크립트
- `README.md`, `AGENTS.md`, `CONTRIBUTING.md`, `docs/architecture.md` 를 포함한 저장소 문서
- 포맷, vet, 테스트, 빌드, 스모크 검증을 수행하는 GitHub Actions CI

## 머지 전 리뷰 기반 하드닝

- 로컬 웹 서버는 루프백 전용 바인딩을 강제합니다.
- 뷰어 메타데이터는 `innerHTML` 대신 텍스트 안전 DOM 업데이트로 렌더링됩니다.
- 탐색 로직은 심볼릭 링크 탈출과 대용량 파일을 거부합니다.
- `globalai web --help` 는 명시적인 서브커맨드 도움말을 제공합니다.
- zero-source API 응답은 `null` 대신 빈 배열을 반환합니다.

## 검증 근거

- `go test ./...`
- `go test -race ./...`
- `go vet ./...`
- `go build ./cmd/globalai`
- `bash scripts/functional_smoke.sh`
- `globalai web` 에 대한 수동 브라우저/API QA
- PR #1과 머지된 `main` 브랜치 빌드에서 성공한 GitHub Actions 실행

## 제품 결정

- 표준 라이브러리 우선
- 프런트엔드 빌드 파이프라인 없음
- 광범위한 재귀 스캔 대신 allowlist 기반 탐색
- 원본 프롬프트 파일 내용을 HTML로 신뢰하지 않는 텍스트 렌더링

## 후속 아이디어

- 실제 수요가 확인되면 더 많은 AI 도구 소스 계열 추가
- 명령 표면이 안정되면 릴리스 자동화 도입
- 소스 카탈로그가 커지면 뷰어 내 검색과 필터링 추가
