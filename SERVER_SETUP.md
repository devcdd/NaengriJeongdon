# 서버 실행 가이드 (Go + Swagger + PostgreSQL)

## 1) 사전 준비

- Go (최신 안정 버전)
- Docker + Docker Compose

## 2) PostgreSQL 시작 (최신 버전)

```bash
cd /
make db-up
```

- PostgreSQL image: `postgres:latest`
- 기본 DB/사용자/비밀번호: `naengrijeongdon`

## 3) 환경 변수 설정

```bash
cd /apps/server
cp .env.example .env
```

기본 `DATABASE_URL`:

```env
DATABASE_URL=postgres://naengrijeongdon:naengrijeongdon@localhost:5433/naengrijeongdon?sslmode=disable

KAKAO_REST_API_KEY=your-kakao-rest-api-key
KAKAO_CLIENT_SECRET=your-kakao-client-secret
KAKAO_ADMIN_KEY=your-kakao-admin-key
KAKAO_REDIRECT_URI=http://localhost:8080/api/v1/auth/kakao/callback

AUTH_JWT_SECRET=your-strong-random-secret
AUTH_ACCESS_TOKEN_TTL_MINUTES=60
AUTH_REFRESH_TOKEN_TTL_DAYS=30
AUTH_SIGNUP_TOKEN_TTL_MINUTES=30
```

## 4) 의존성 설치/업데이트 (최신 버전)

```bash
cd /
make deps
```

## 5) Swagger 문서 생성

```bash
cd /
make swagger
```

## 6) DB 마이그레이션 실행

```bash
cd /
make migrate-up
```

유용한 명령:

```bash
make migrate-down
make migrate-force VERSION=1
make migrate-create NAME=add_some_column
```

## 7) 서버 실행

```bash
cd /
make run
```

- API 기본 경로: `http://localhost:8080/api/v1`
- Swagger UI 경로: `http://localhost:8080/swagger/index.html`

## 8) 헬스 체크

```bash
curl http://localhost:8080/api/v1/health/live
curl http://localhost:8080/api/v1/health/ready
```

## 9) 카카오 OAuth 확인용 엔드포인트

```bash
curl "http://localhost:8080/api/v1/auth/kakao/login-url"
```

Swagger에서 아래 흐름을 확인할 수 있습니다.
- `GET /api/v1/auth/kakao/login-url`
- `GET /api/v1/auth/kakao/callback`
- `POST /api/v1/auth/register`
- `POST /api/v1/auth/refresh`
- `POST /api/v1/auth/logout`
- `GET /api/v1/auth/session`

## 10) 데이터베이스 중지

```bash
cd /
make db-down
```

## 11) 현재 DB 백업(복사)하기

```bash
cd /
make db-dump
```

생성 결과:
- `/backups/naengrijeongdon_YYYYMMDD_HHMMSS.sql`
- 예: `/backups/naengrijeongdon_20260218_235500.sql`

## 12) 백업 파일 DB에 다시 넣기(복원)

기존 데이터 유지한 채 복원:

```bash
cd /
make db-restore FILE=backups/naengrijeongdon_20260218_235500.sql
```

`FILE`를 생략하면 `backups/` 목록에서 번호로 선택할 수 있습니다.

```bash
cd /
make db-restore
```

기존 데이터 전부 비우고 복원(초기화 후 복원):

```bash
cd /
make db-restore-reset FILE=backups/naengrijeongdon_20260218_235500.sql
```

`FILE`를 생략하면 `backups/` 목록에서 번호로 선택할 수 있습니다.

```bash
cd /
make db-restore-reset
```

주의:
- `db-restore-reset`은 현재 DB 데이터를 삭제하므로 테스트/복구 상황에서만 사용하세요.
