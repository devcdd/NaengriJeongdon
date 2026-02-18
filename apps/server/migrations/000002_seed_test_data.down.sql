delete from user_food_search_history
where user_id = '00000000-0000-0000-0000-000000000001'
  and normalized_name = '테스트우유';

delete from pantry_items
where id = '00000000-0000-0000-0000-000000000006';

delete from shelf_life_rules
where id = '00000000-0000-0000-0000-000000000004';

delete from food_recommended_storage_types
where food_id = '00000000-0000-0000-0000-000000000002';

delete from food_aliases
where id = '00000000-0000-0000-0000-000000000003';

delete from foods
where id = '00000000-0000-0000-0000-000000000002';

delete from storage_spaces
where id = '00000000-0000-0000-0000-000000000005';

delete from app_users
where id = '00000000-0000-0000-0000-000000000001';
