create extension if not exists pgcrypto;
create extension if not exists pg_trgm;

create table app_users (
  id uuid primary key default gen_random_uuid(),
  email text unique not null,
  display_name text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table storage_types (
  id smallserial primary key,
  code text not null unique,
  label text not null
);

create table storage_spaces (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references app_users(id) on delete cascade,
  name text not null,
  storage_type_id smallint not null references storage_types(id),
  is_archived boolean not null default false,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  unique (user_id, name)
);

create index idx_storage_spaces_user on storage_spaces(user_id);

create table food_categories (
  id smallserial primary key,
  code text not null unique,
  label text not null
);

create table foods (
  id uuid primary key default gen_random_uuid(),
  name text not null unique,
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

create table food_recommended_storage_types (
  food_id uuid not null references foods(id) on delete cascade,
  storage_type_id smallint not null references storage_types(id),
  priority smallint not null default 1,
  primary key (food_id, storage_type_id)
);

create table shelf_life_rules (
  id uuid primary key default gen_random_uuid(),
  food_id uuid references foods(id) on delete cascade,
  category_id smallint references food_categories(id) on delete cascade,
  storage_type_id smallint references storage_types(id),
  shelf_life_days integer not null check (shelf_life_days > 0),
  rule_source text not null default 'system',
  confidence numeric(3,2) not null default 0.80 check (confidence >= 0 and confidence <= 1),
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  check (
    (food_id is not null and category_id is null) or
    (food_id is null and category_id is not null)
  )
);

create unique index uq_shelf_life_food_scope
  on shelf_life_rules (food_id, coalesce(storage_type_id, -1))
  where food_id is not null;

create unique index uq_shelf_life_category_scope
  on shelf_life_rules (category_id, coalesce(storage_type_id, -1))
  where category_id is not null;

create table pantry_items (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references app_users(id) on delete cascade,
  storage_space_id uuid not null references storage_spaces(id),
  food_id uuid references foods(id),
  input_name text not null,
  quantity numeric(10,2),
  quantity_unit text,
  purchased_at date,
  expires_at date,
  is_estimated_expiry boolean not null default false,
  expiry_estimated_by text,
  estimation_rule_id uuid references shelf_life_rules(id),
  status_override text,
  notes text,
  storage_warning_triggered boolean not null default false,
  storage_warning_acknowledged boolean not null default false,
  storage_warning_message text,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now(),
  check (purchased_at is not null or expires_at is not null),
  check (purchased_at is null or expires_at is null or purchased_at <= expires_at)
);

create index idx_pantry_items_user on pantry_items(user_id);
create index idx_pantry_items_space on pantry_items(storage_space_id);
create index idx_pantry_items_food on pantry_items(food_id);
create index idx_pantry_items_expires_at on pantry_items(expires_at);
create index idx_pantry_items_input_name on pantry_items using gin (input_name gin_trgm_ops);

create table user_food_search_history (
  user_id uuid not null references app_users(id) on delete cascade,
  normalized_name text not null,
  last_used_at timestamptz not null default now(),
  use_count integer not null default 1 check (use_count >= 1),
  primary key (user_id, normalized_name)
);

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
