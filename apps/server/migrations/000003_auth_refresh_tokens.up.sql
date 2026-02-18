create table auth_refresh_tokens (
  id uuid primary key default gen_random_uuid(),
  user_id uuid not null references app_users(id) on delete cascade,
  token_hash text not null unique,
  expires_at timestamptz not null,
  revoked_at timestamptz,
  created_at timestamptz not null default now()
);

create index idx_auth_refresh_tokens_user_id on auth_refresh_tokens(user_id);
create index idx_auth_refresh_tokens_expires_at on auth_refresh_tokens(expires_at);
