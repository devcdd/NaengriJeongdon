insert into app_users (id, email, display_name)
values (
  '00000000-0000-0000-0000-000000000001',
  'test-user@naengrijeongdon.local',
  '테스트유저'
);

insert into storage_spaces (id, user_id, name, storage_type_id)
values (
  '00000000-0000-0000-0000-000000000005',
  '00000000-0000-0000-0000-000000000001',
  '테스트 메인 냉장고',
  (select id from storage_types where code = 'FRIDGE')
);

insert into foods (id, name, category_id)
values (
  '00000000-0000-0000-0000-000000000002',
  '테스트우유',
  (select id from food_categories where code = 'DAIRY')
);

insert into food_aliases (id, food_id, alias)
values (
  '00000000-0000-0000-0000-000000000003',
  '00000000-0000-0000-0000-000000000002',
  '우유테스트'
);

insert into food_recommended_storage_types (food_id, storage_type_id, priority)
values (
  '00000000-0000-0000-0000-000000000002',
  (select id from storage_types where code = 'FRIDGE'),
  1
);

insert into shelf_life_rules (id, food_id, storage_type_id, shelf_life_days, rule_source, confidence)
values (
  '00000000-0000-0000-0000-000000000004',
  '00000000-0000-0000-0000-000000000002',
  (select id from storage_types where code = 'FRIDGE'),
  7,
  'system',
  0.95
);

insert into pantry_items (
  id,
  user_id,
  storage_space_id,
  food_id,
  input_name,
  quantity,
  quantity_unit,
  purchased_at,
  expires_at,
  is_estimated_expiry,
  expiry_estimated_by,
  estimation_rule_id,
  storage_warning_triggered,
  storage_warning_acknowledged
)
values (
  '00000000-0000-0000-0000-000000000006',
  '00000000-0000-0000-0000-000000000001',
  '00000000-0000-0000-0000-000000000005',
  '00000000-0000-0000-0000-000000000002',
  '테스트우유',
  1,
  'ea',
  current_date,
  current_date + 7,
  true,
  'FOOD_RULE',
  '00000000-0000-0000-0000-000000000004',
  false,
  false
);

insert into user_food_search_history (user_id, normalized_name, last_used_at, use_count)
values (
  '00000000-0000-0000-0000-000000000001',
  '테스트우유',
  now(),
  1
);
