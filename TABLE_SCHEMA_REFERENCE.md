# 테이블 스키마 레퍼런스

이 문서는 현재 DB 스키마를 테이블별로 `key / type / nullable / description` 형식으로 정리한 문서입니다.

## 1) `app_users`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 사용자 고유 ID. 기본값 `gen_random_uuid()` |
| `email` | `text` | NOT NULL | 사용자 이메일. UNIQUE |
| `display_name` | `text` | NULL | 사용자 표시 이름(닉네임) |
| `is_admin` | `boolean` | NOT NULL | 관리자 계정 여부 (기본값 `false`) |
| `provider` | `text` | NOT NULL | 로그인 제공자(예: `kakao`, `local`) |
| `provider_user_id` | `text` | NULL | 제공자 측 사용자 ID(예: 카카오 사용자 ID) |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |
| `updated_at` | `timestamptz` | NOT NULL | 수정 시각 |

부연설명:
- 모든 사용자별 데이터(`storage_spaces`, `pantry_items`, `user_food_search_history`)의 기준 테이블입니다.
- `provider + provider_user_id` 조합은 부분 유니크 인덱스로 관리됩니다(`provider_user_id`가 있을 때만).

---

## 2) `storage_types`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `smallserial` | NOT NULL | PK. 공간 타입 ID |
| `code` | `text` | NOT NULL | 공간 타입 코드(예: `FRIDGE`, `FREEZER`) UNIQUE |
| `label` | `text` | NOT NULL | 사용자 노출용 타입 이름(예: 냉장, 냉동) |

부연설명:
- 사용자 커스텀 공간(`storage_spaces`)이 어떤 성격의 공간인지 판단할 때 사용됩니다.

---

## 3) `storage_spaces`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 보관 공간 ID |
| `user_id` | `uuid` | NOT NULL | FK -> `app_users.id`. 공간 소유자 |
| `name` | `text` | NOT NULL | 사용자 커스텀 공간 이름(예: 메인 냉장고) |
| `storage_type_id` | `smallint` | NOT NULL | FK -> `storage_types.id`. 공간 타입 |
| `is_archived` | `boolean` | NOT NULL | 보관 공간 보관처리 여부 |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |
| `updated_at` | `timestamptz` | NOT NULL | 수정 시각 |

부연설명:
- UNIQUE(`user_id`, `name`) 제약으로 사용자 단위 중복 이름을 방지합니다.
- `pantry_items`는 반드시 하나의 `storage_space`에 속합니다.

---

## 4) `food_categories`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `smallserial` | NOT NULL | PK. 카테고리 ID |
| `code` | `text` | NOT NULL | 카테고리 코드(예: `DAIRY`, `MEAT`) UNIQUE |
| `label` | `text` | NOT NULL | 카테고리 이름(예: 유제품, 육류) |

부연설명:
- 식품 직접 규칙이 없을 때 `shelf_life_rules`의 대체 규칙(fallback) 기준이 됩니다.

---

## 5) `foods`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 식품 마스터 ID |
| `name` | `text` | NOT NULL | 표준 식품명(정규 이름). UNIQUE |
| `category_id` | `smallint` | NULL | FK -> `food_categories.id` |
| `is_active` | `boolean` | NOT NULL | 사용 여부(비활성 처리 가능) |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |
| `updated_at` | `timestamptz` | NOT NULL | 수정 시각 |

부연설명:
- 자동완성, 규칙 매칭, 권장 보관 방식 판단의 중심 테이블입니다.

---

## 6) `food_aliases`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 별칭 ID |
| `food_id` | `uuid` | NOT NULL | FK -> `foods.id` |
| `alias` | `text` | NOT NULL | 별칭/검색어(예: 방토) |

부연설명:
- UNIQUE(`alias`)로 별칭 충돌을 방지합니다.
- 자동완성 시 `foods.name`과 함께 검색 대상으로 사용합니다.

---

## 7) `food_recommended_storage_types`

| key | type | nullable | description |
|---|---|---|---|
| `food_id` | `uuid` | NOT NULL | PK(FK). FK -> `foods.id` |
| `storage_type_id` | `smallint` | NOT NULL | PK(FK). FK -> `storage_types.id` |
| `priority` | `smallint` | NOT NULL | 권장 우선순위(1이 최우선) |

부연설명:
- 하나의 식품이 여러 보관 타입을 가질 수 있습니다.
- 사용자가 선택한 공간 타입과 불일치 시 경고(soft warning) 판단 근거가 됩니다.

---

## 8) `shelf_life_rules`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 규칙 ID |
| `food_id` | `uuid` | NULL | FK -> `foods.id`. 식품 직접 매칭용 |
| `category_id` | `smallint` | NULL | FK -> `food_categories.id`. 카테고리 대체 규칙용 |
| `storage_type_id` | `smallint` | NULL | FK -> `storage_types.id`. 타입별 규칙 분기 |
| `shelf_life_days` | `integer` | NOT NULL | 기본 보관일(예측 계산 기준) |
| `rule_source` | `text` | NOT NULL | 규칙 출처(예: `system`, `admin`, `user`) |
| `confidence` | `numeric(3,2)` | NOT NULL | 규칙 신뢰도(0~1) |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |
| `updated_at` | `timestamptz` | NOT NULL | 수정 시각 |

부연설명:
- CHECK 제약으로 `food_id`와 `category_id`는 정확히 하나만 채워야 합니다.
- 예측 우선순위 예시:
  1. `food_id + storage_type_id`
  2. `food_id + null`
  3. `category_id + storage_type_id`
  4. `category_id + null`

---

## 9) `pantry_items`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 실제 보유 식품 항목 ID |
| `user_id` | `uuid` | NOT NULL | FK -> `app_users.id` |
| `storage_space_id` | `uuid` | NOT NULL | FK -> `storage_spaces.id` |
| `food_id` | `uuid` | NULL | FK -> `foods.id`. 목록 미매핑 시 NULL 가능 |
| `input_name` | `text` | NOT NULL | 사용자가 입력한 원문 식품명 |
| `quantity` | `numeric(10,2)` | NULL | 수량 |
| `quantity_unit` | `text` | NULL | 수량 단위(예: ea, g, ml) |
| `purchased_at` | `date` | NULL | 구매일 |
| `expires_at` | `date` | NULL | 유통기한/마감일 |
| `is_estimated_expiry` | `boolean` | NOT NULL | 유통기한이 예측값인지 여부 |
| `expiry_estimated_by` | `text` | NULL | 예측 근거(예: `FOOD_RULE`, `CATEGORY_RULE`) |
| `estimation_rule_id` | `uuid` | NULL | FK -> `shelf_life_rules.id`. 사용된 규칙 |
| `status_override` | `text` | NULL | 상태 수동 오버라이드 |
| `notes` | `text` | NULL | 메모 |
| `storage_warning_triggered` | `boolean` | NOT NULL | 보관 방식 경고 발생 여부 |
| `storage_warning_acknowledged` | `boolean` | NOT NULL | 사용자가 경고 확인했는지 여부 |
| `storage_warning_message` | `text` | NULL | 경고 메시지 원문 |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |
| `updated_at` | `timestamptz` | NOT NULL | 수정 시각 |

부연설명:
- CHECK 제약:
  - `purchased_at` 또는 `expires_at` 중 최소 1개 필수
  - 둘 다 있을 경우 `purchased_at <= expires_at`
- 목록 필터링(`만료`, `임박`, `오늘 마감`)은 주로 `expires_at` 기준으로 계산합니다.

---

## 10) `user_food_search_history`

| key | type | nullable | description |
|---|---|---|---|
| `user_id` | `uuid` | NOT NULL | PK(FK). FK -> `app_users.id` |
| `normalized_name` | `text` | NOT NULL | PK. 사용자 입력 식품명(정규화) |
| `last_used_at` | `timestamptz` | NOT NULL | 마지막 사용 시각 |
| `use_count` | `integer` | NOT NULL | 누적 사용 횟수 |

부연설명:
- 자동완성 추천 순위 개선(최근성/빈도)에 활용합니다.

---

## 11) `auth_refresh_tokens`

| key | type | nullable | description |
|---|---|---|---|
| `id` | `uuid` | NOT NULL | PK. 리프레시 토큰 레코드 ID |
| `user_id` | `uuid` | NOT NULL | FK -> `app_users.id` |
| `token_hash` | `text` | NOT NULL | 리프레시 토큰 SHA-256 해시값(UNIQUE) |
| `expires_at` | `timestamptz` | NOT NULL | 만료 시각 |
| `revoked_at` | `timestamptz` | NULL | 폐기(로그아웃) 처리 시각 |
| `created_at` | `timestamptz` | NOT NULL | 생성 시각 |

부연설명:
- 원문 토큰 대신 해시만 저장합니다.
- 토큰 재발급 시 기존 토큰은 `revoked_at`을 채워 무효화합니다.

---

## 12) `schema_migrations` (golang-migrate 관리 테이블)

| key | type | nullable | description |
|---|---|---|---|
| `version` | `bigint` | NOT NULL | 적용된 최신 마이그레이션 버전 |
| `dirty` | `boolean` | NOT NULL | 실패 중단 상태 여부 |

부연설명:
- `make migrate-up/down/force` 실행 상태를 추적하는 시스템 테이블입니다.
