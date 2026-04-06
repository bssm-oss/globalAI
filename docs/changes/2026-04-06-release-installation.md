# 2026-04-06 — 설치 및 GitHub 릴리스 경로 정리

## 요약

이번 변경은 `globalAI`를 처음 접한 사용자가 설치 경로를 바로 이해할 수 있도록 README를 정리하고, 버전 태그를 푸시하면 GitHub Releases에서 즉시 내려받을 수 있는 바이너리 자산을 자동으로 게시하도록 릴리스 워크플로를 추가한 작업입니다.

## 무엇이 달라졌는가

- README에 `go install github.com/bssm-oss/globalAI/cmd/globalai@latest` 기반 설치 경로를 추가했습니다.
- README에 GitHub Releases에서 운영체제별 압축 바이너리를 내려받는 경로를 추가했습니다.
- `.github/workflows/release.yml` 을 추가해 버전 태그(`v*`) 푸시 시 검증 후 macOS, Linux, Windows용 빌드를 만들고 GitHub Release 자산으로 게시하도록 했습니다.

## 왜 필요한가

- 기존 저장소는 로컬 빌드와 CI 검증은 잘 설명하고 있었지만, 사용자가 "어떻게 설치하면 되는지"를 첫 화면에서 바로 이해하기 어려웠습니다.
- `go install` 경로는 이미 모듈 구조상 가능했지만 README에 드러나지 않았습니다.
- GitHub Releases 자동화가 없어서 GitHub에서 바로 바이너리를 내려받는 배포 경험이 실제로 제공되지 않았습니다.

## 검증 근거

- `go test ./...`
- `go build ./cmd/globalai`
- `bash scripts/functional_smoke.sh`
- `go install ./cmd/globalai`
- `go install github.com/bssm-oss/globalAI/cmd/globalai@latest`
- 설치된 `globalai --help` 실행 확인
