-- NaengriJeongdon MVP schema (PostgreSQL)
-- Focus:
-- 1) User-defined storage spaces
-- 2) Food master + category fallback shelf-life prediction
-- 3) Recommended storage warnings
-- 4) Pantry item requires purchased_at or expires_at (at least one)

create extension if not exists pgcrypto;
create extension if not exists pg_trgm;

-- -----------------------------
-- Core users
-- -----------------------------
create table app_users (
  id uuid primary key default gen_random_uuid(),
  email text unique not null,
  display_name text,
  is_admin boolean not null default false,
  provider text not null default 'local',
  provider_user_id text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create unique index uq_app_users_provider_user_id
  on app_users (provider, provider_user_id)
  where provider_user_id is not null;

-- -----------------------------
-- Storage model
-- -----------------------------
create table storage_types (
  id smallserial primary key,
  code text not null unique,      -- FRIDGE, FREEZER, KIMCHI, ROOM
  label text not null
);

create table storage_spaces (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references app_users(id) on delete cascade,
  name text not null,             -- user-customizable, e.g. "Main Fridge"
  storage_type_id smallint not null references storage_types(id),
  is_archived boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, name)
);

create index idx_storage_spaces_user on storage_spaces(user_id);

-- -----------------------------
-- Food taxonomy / dictionary
-- -----------------------------
create table food_categories (
  id smallserial primary key,
  code text not null unique,      -- DAIRY, MEAT, VEGETABLE, FRUIT...
  label text not null
);

create table foods (
  id uuid primary key default gen_random_uuid(),
  name text not null unique,      -- canonical name
  category_id smallint references food_categories(id),
  is_active boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table food_aliases (
  id uuid primary key default gen_random_uuid(),
  food_id uuid not null references foods(id) on delete cascade,
  alias text not null,
  unique (food_id, alias),
  unique (alias)
);

-- Foods can recommend one or many storage types.
-- Example: frozen dumpling -> FREEZER only
-- Example: bread -> ROOM or FREEZER
create table food_recommended_storage_types (
  food_id uuid not null references foods(id) on delete cascade,
  storage_type_id smallint not null references storage_types(id),
  priority smallint not null default 1, -- 1=best, bigger=acceptable
  primary key (food_id, storage_type_id)
);

-- -----------------------------
-- Shelf-life rules
-- -----------------------------
-- Exactly one scope must exist:
-- 1) food_id for exact match
-- 2) category_id for fallback
create table shelf_life_rules (
  id uuid primary key default gen_random_uuid(),
  food_id uuid references foods(id) on delete cascade,
  category_id smallint references food_categories(id) on delete cascade,
  storage_type_id smallint references storage_types(id),
  shelf_life_days integer not null check (shelf_life_days > 0),
  rule_source text not null default 'system', -- system, admin, user
  confidence numeric(3,2) not null default 0.80 check (confidence >= 0 and confidence <= 1),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  check (
    (food_id is not null and category_id is null) or
    (food_id is null and category_id is not null)
  )
);

-- Exact food rules are unique per storage type (including null type as generic)
create unique index uq_shelf_life_food_scope
  on shelf_life_rules (food_id, coalesce(storage_type_id, -1))
  where food_id is not null;

-- Category fallback rules are unique per storage type
create unique index uq_shelf_life_category_scope
  on shelf_life_rules (category_id, coalesce(storage_type_id, -1))
  where category_id is not null;

-- -----------------------------
-- Pantry items (actual user inventory)
-- -----------------------------
create table pantry_items (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references app_users(id) on delete cascade,
  storage_space_id uuid not null references storage_spaces(id),
  food_id uuid references foods(id),      -- nullable for custom/unlisted entry
  input_name text not null,               -- raw text user typed
  quantity numeric(10,2),
  quantity_unit text,                     -- ea, g, ml...

  purchased_at date,
  expires_at date,
  is_estimated_expiry boolean not null default false,
  expiry_estimated_by text,               -- FOOD_RULE, CATEGORY_RULE, NONE
  estimation_rule_id uuid references shelf_life_rules(id),

  status_override text,                   -- optional manual status override
  notes text,

  storage_warning_triggered boolean not null default false,
  storage_warning_acknowledged boolean not null default false,
  storage_warning_message text,

  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),

  -- At least one date is required
  check (purchased_at is not null or expires_at is not null),
  -- If both exist, purchase date cannot be later than expiry
  check (purchased_at is null or expires_at is null or purchased_at <= expires_at)
);

create index idx_pantry_items_user on pantry_items(user_id);
create index idx_pantry_items_space on pantry_items(storage_space_id);
create index idx_pantry_items_food on pantry_items(food_id);
create index idx_pantry_items_expires_at on pantry_items(expires_at);
create index idx_pantry_items_input_name on pantry_items using gin (input_name gin_trgm_ops);

-- User-personalized autocomplete ranking (recently used names)
create table user_food_search_history (
  user_id uuid not null references app_users(id) on delete cascade,
  normalized_name text not null,
  last_used_at timestamptz not null default now(),
  use_count integer not null default 1 check (use_count >= 1),
  primary key (user_id, normalized_name)
);

-- -----------------------------
-- Seed candidates (minimal)
-- -----------------------------
insert into storage_types (code, label) values
  ('FRIDGE', '냉장'),
  ('FREEZER', '냉동'),
  ('KIMCHI', '김치냉장'),
  ('ROOM', '실온')
on conflict (code) do nothing;

insert into food_categories (code, label) values
  ('DAIRY', '유제품'),
  ('MEAT', '육류'),
  ('SEAFOOD', '해산물'),
  ('VEGETABLE', '채소'),
  ('FRUIT', '과일'),
  ('KIMCHI_SIDE', '김치/반찬'),
  ('FROZEN_FOOD', '냉동식품'),
  ('SAUCE', '소스/조미료'),
  ('DRINK', '음료'),
  ('ETC', '기타')
on conflict (code) do nothing;
