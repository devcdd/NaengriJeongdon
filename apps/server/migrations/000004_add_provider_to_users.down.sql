drop index if exists uq_app_users_provider_user_id;

alter table app_users
  drop column if exists provider_user_id,
  drop column if exists provider;
