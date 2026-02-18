# 냉장고 MVP 스키마 가이드

이 문서는 `schema.sql`을 빠르게 이해하기 위한 요약 가이드입니다.

## 1. 목표

- 냉장고/냉동고를 고정 개념이 아닌 **보관 공간**으로 관리
- 식품 등록 시 **유통기한 자동 예측**
  - 1순위: 식품 개별 규칙
  - 2순위: 카테고리 규칙(대체 규칙)
- 식품의 **권장 보관 방식**과 사용자가 고른 공간 타입이 다르면 **경고**
- 등록 시 날짜는 `구매일` 또는 `유통기한` 중 **최소 1개 필수**
- 목록 화면에서 **유통기한 상태 필터링**으로 빠른 확인 지원

---

## 2. 핵심 엔티티

### `app_users`
- 사용자 기본 정보

### `storage_types`
- 공간 타입 마스터
- 예: `FRIDGE`, `FREEZER`, `KIMCHI`, `ROOM`

### `storage_spaces`
- 사용자가 직접 만드는 보관 공간
- 예: `우리집 메인 냉장고`, `김치냉장고`, `서브 냉동고`
- `name`은 사용자별 유니크

### `food_categories`
- 식품 카테고리 마스터
- 예: `DAIRY`, `MEAT`, `VEGETABLE`

### `foods`
- 식품 사전(정규 식품명)
- 카테고리 연결 가능

### `food_aliases`
- 자동완성/검색용 별칭
- 예: `방토` -> `방울토마토`

### `food_recommended_storage_types`
- 식품별 권장 보관 타입(복수 허용)
- 예: 만두는 FREEZER, 빵은 ROOM/FREEZER 둘 다 가능

### `shelf_life_rules`
- 유통기한 예측 규칙
- 규칙 스코프:
  1. `food_id` 기반(정밀)
  2. `category_id` 기반(대체 규칙)
- `storage_type_id`로 타입별 규칙 분기 가능

### `pantry_items`
- 실제 사용자가 등록한 식재료 인벤토리
- 핵심 제약:
  - `purchased_at` 또는 `expires_at` 최소 1개 필수
  - 둘 다 있으면 `purchased_at <= expires_at`
- 예측/경고 상태 저장:
  - `is_estimated_expiry`
  - `expiry_estimated_by`
  - `estimation_rule_id`
  - `storage_warning_triggered`
  - `storage_warning_acknowledged`

### `user_food_search_history`
- 사용자별 최근 입력 식품명 히스토리(자동완성 랭킹용)

---

## 2-1. 테이블 구성 요약 (질문 기준)

질문한 방향대로, MVP는 아래 3축으로 보면 가장 이해가 쉽습니다.

1. 사용자 축
- `app_users`

2. 음식 매핑 축
- `food_categories`
- `foods`
- `food_aliases`
- `food_recommended_storage_types`
- `shelf_life_rules`

3. 사용자 공간/재고 축
- `storage_spaces`
- `pantry_items`
- `user_food_search_history` (자동완성 품질 향상용)

정리:
- 유저 테이블: `app_users`
- 음식(매핑용) 테이블: `foods` (+ `food_aliases`, `food_categories`)
- 유저의 공간/음식 관리 테이블: `storage_spaces`, `pantry_items`

---

## 3. 예측 로직 (MVP)

식품 등록 시 유통기한 계산 우선순위:

1. `food_id + storage_type_id` 규칙 조회
2. `food_id + null(storage_type)` 규칙 조회
3. `category_id + storage_type_id` 규칙 조회
4. `category_id + null(storage_type)` 규칙 조회
5. 없으면 사용자 직접 입력 유도

계산식:
- `예상 유통기한 = purchased_at + shelf_life_days`

---

## 4. 경고 로직 (MVP)

1. 식품의 권장 보관 타입(`food_recommended_storage_types`) 조회
2. 사용자가 선택한 공간의 타입(`storage_spaces.storage_type_id`) 비교
3. 불일치면 경고 표시(저장 차단 아님, 소프트 경고)
4. 사용자가 확인 후 저장하면 `storage_warning_acknowledged = true`

---

## 5. 자동완성 로직 (MVP)

추천 소스:
1. `foods.name`, `food_aliases.alias`
2. `user_food_search_history.normalized_name`

검색 품질:
- `pantry_items.input_name` trigram(3-gram) 인덱스(`pg_trgm`) 사용 가능

---

## 6. 유통기한 필터링 (MVP 화면)

목록 화면에서 아래 필터를 기본 제공:

- `전체`
- `만료` : `expires_at < today`
- `오늘 마감` : `expires_at = today`
- `임박(3일)` : `expires_at between today and today+3`
- `임박(7일)` : `expires_at between today and today+7`
- `여유` : `expires_at > today+7`
- `미입력` : `expires_at is null` (구매일만 있고 예측 실패한 케이스 포함)

추가 필터:

- 보관 공간(`storage_space_id`)
- 정렬: `expires_at asc` (기본), `최근 등록순`

참고:
- 현재 스키마의 `pantry_items.expires_at` 인덱스로 MVP 트래픽에서는 충분
- 기본 상태값은 서버 계산식으로 처리하고, 필요하면 나중에 머티리얼라이즈드 뷰/상태 컬럼으로 최적화

---

## 7. 데이터 흐름 예시

1. 사용자가 `우유` 입력
2. 자동완성에서 `foods` 매칭
3. 공간으로 `김치냉장고(KIMCHI)` 선택
4. 우유 권장 타입이 `FRIDGE`라면 경고 노출
5. 구매일이 있으면 규칙으로 예상 유통기한 계산
6. `pantry_items`에 저장 + 예측/경고 상태 기록

---

## 8. 현재 파일

- 스키마 SQL: `/schema.sql`
- 설명 문서(현재 파일): `/MVP_SCHEMA_GUIDE.md`
