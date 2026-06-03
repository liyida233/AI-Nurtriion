USE ai_nutrition;

INSERT INTO exercises (id, name, category, muscle_group, equipment, intensity_level, created_at, updated_at)
VALUES
  ('11111111-1111-1111-1111-111111111111', 'Squat', 'strength', 'lower_body', 'barbell', 'moderate', NOW(6), NOW(6)),
  ('22222222-2222-2222-2222-222222222222', 'Bench Press', 'strength', 'chest', 'barbell', 'moderate', NOW(6), NOW(6)),
  ('33333333-3333-3333-3333-333333333333', 'Lat Pulldown', 'strength', 'back', 'machine', 'moderate', NOW(6), NOW(6)),
  ('44444444-4444-4444-4444-444444444444', 'Running', 'cardio', 'full_body', 'none', 'variable', NOW(6), NOW(6))
ON DUPLICATE KEY UPDATE
  category = VALUES(category),
  muscle_group = VALUES(muscle_group),
  equipment = VALUES(equipment),
  intensity_level = VALUES(intensity_level),
  updated_at = NOW(6);

INSERT INTO food_items (id, name, serving_size, calories, protein, carbohydrates, fat, sugar, sodium, created_at, updated_at)
VALUES
  ('aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa', 'Chicken Breast', '100g', 165, 31, 0, 3.6, 0, 74, NOW(6), NOW(6)),
  ('bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb', 'Cooked Rice', '100g', 130, 2.7, 28, 0.3, 0.1, 1, NOW(6), NOW(6)),
  ('cccccccc-cccc-cccc-cccc-cccccccccccc', 'Egg', '1 large', 72, 6.3, 0.4, 4.8, 0.2, 71, NOW(6), NOW(6)),
  ('dddddddd-dddd-dddd-dddd-dddddddddddd', 'Banana', '1 medium', 105, 1.3, 27, 0.4, 14.4, 1, NOW(6), NOW(6))
ON DUPLICATE KEY UPDATE
  serving_size = VALUES(serving_size),
  calories = VALUES(calories),
  protein = VALUES(protein),
  carbohydrates = VALUES(carbohydrates),
  fat = VALUES(fat),
  sugar = VALUES(sugar),
  sodium = VALUES(sodium),
  updated_at = NOW(6);
