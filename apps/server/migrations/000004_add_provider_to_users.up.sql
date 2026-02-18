alter table app_users
  add column provider text not null default 'local',
  add column provider_user_id text;

create unique index uq_app_users_provider_user_id
  on app_users (provider, provider_user_id)
  where provider_user_id is not null;
