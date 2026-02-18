# 냉리정돈 (NaengriJeongdon)

자취/가정용 냉장고 식재료를 데이터로 관리하는 프로젝트입니다.  
현재는 **Go + Swagger + PostgreSQL** 기반 서버가 구성되어 있습니다.

## 기술 스택

- Go
- Gin
- Swagger (`swaggo`)
- PostgreSQL (`postgres:latest`)
- golang-migrate

## 프로젝트 구조

```text
/
├── apps/server              # Go API 서버
├── backups                  # DB 백업 파일 생성 위치(로컬)
├── schema.sql               # 전체 스키마 참고 파일
├── Makefile                 # 실행/마이그레이션/백업 명령
└── SERVER_SETUP.md          # 상세 실행 가이드
```

## 빠른 시작

```bash
cd /
make db-up
cd /apps/server
cp .env.example .env
cd /
make migrate-up
make run
```

- API base: `http://localhost:8080/api/v1`
- Swagger: `http://localhost:8080/swagger`

## 주요 명령

```bash
# DB
make db-up
make db-down

# Migration
make migrate-up
make migrate-down

# Swagger
make swagger

# DB backup / restore
make db-dump
make db-restore FILE=backups/naengrijeongdon_YYYYMMDD_HHMMSS.sql
make db-restore-reset FILE=backups/naengrijeongdon_YYYYMMDD_HHMMSS.sql
```

## 인증(Auth)

- 카카오 OAuth 기반 로그인/회원가입 플로우를 사용합니다.
- 가입 필요 응답(`signupRequired=true`)에서는 `signupToken`을 발급하고,  
  `POST /api/v1/auth/register` 성공 시 `accessToken`, `refreshToken`을 발급합니다.

## 보안 주의

- 실제 키/시크릿은 `apps/server/.env`에만 두고 커밋하지 마세요.
- 이미 노출된 키가 있다면 즉시 카카오 콘솔에서 재발급/회전하세요.
- DB 백업(`backups/*.sql`)에는 개인정보가 포함될 수 있으니 외부 공유에 주의하세요.

