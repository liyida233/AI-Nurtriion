CREATE DATABASE IF NOT EXISTS ai_nutrition
  CHARACTER SET utf8mb4
  COLLATE utf8mb4_unicode_ci;

USE ai_nutrition;

CREATE TABLE IF NOT EXISTS users (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(120) NOT NULL,
  email VARCHAR(180) NOT NULL UNIQUE,
  password_hash VARCHAR(255) NOT NULL,
  role VARCHAR(40) NOT NULL DEFAULT 'user',
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS user_profiles (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL UNIQUE,
  age INT NOT NULL DEFAULT 0,
  gender VARCHAR(40),
  height_cm DOUBLE NOT NULL DEFAULT 0,
  weight_kg DOUBLE NOT NULL DEFAULT 0,
  activity_level VARCHAR(60),
  primary_goal VARCHAR(80),
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  CONSTRAINT fk_user_profiles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS exercises (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(140) NOT NULL,
  category VARCHAR(80),
  muscle_group VARCHAR(80),
  equipment VARCHAR(80),
  intensity_level VARCHAR(60),
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_exercises_name (name)
);

CREATE TABLE IF NOT EXISTS workout_sessions (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  workout_date DATETIME(6),
  category VARCHAR(80),
  duration_min INT NOT NULL DEFAULT 0,
  notes TEXT,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_workout_sessions_user_id (user_id),
  INDEX idx_workout_sessions_workout_date (workout_date),
  CONSTRAINT fk_workout_sessions_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS workout_entries (
  id CHAR(36) PRIMARY KEY,
  session_id CHAR(36) NOT NULL,
  exercise_id CHAR(36) NOT NULL,
  sets INT NOT NULL DEFAULT 0,
  reps INT NOT NULL DEFAULT 0,
  weight_kg DOUBLE NOT NULL DEFAULT 0,
  rest_sec INT NOT NULL DEFAULT 0,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_workout_entries_session_id (session_id),
  INDEX idx_workout_entries_exercise_id (exercise_id),
  CONSTRAINT fk_workout_entries_session FOREIGN KEY (session_id) REFERENCES workout_sessions(id) ON DELETE CASCADE,
  CONSTRAINT fk_workout_entries_exercise FOREIGN KEY (exercise_id) REFERENCES exercises(id)
);

CREATE TABLE IF NOT EXISTS food_items (
  id CHAR(36) PRIMARY KEY,
  name VARCHAR(140) NOT NULL,
  serving_size VARCHAR(80),
  calories DOUBLE NOT NULL DEFAULT 0,
  protein DOUBLE NOT NULL DEFAULT 0,
  carbohydrates DOUBLE NOT NULL DEFAULT 0,
  fat DOUBLE NOT NULL DEFAULT 0,
  sugar DOUBLE NOT NULL DEFAULT 0,
  sodium DOUBLE NOT NULL DEFAULT 0,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_food_items_name (name)
);

CREATE TABLE IF NOT EXISTS meal_logs (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  food_item_id CHAR(36) NOT NULL,
  meal_type VARCHAR(60),
  quantity DOUBLE NOT NULL DEFAULT 0,
  meal_time DATETIME(6),
  total_calories DOUBLE NOT NULL DEFAULT 0,
  total_protein DOUBLE NOT NULL DEFAULT 0,
  total_carbs DOUBLE NOT NULL DEFAULT 0,
  total_fat DOUBLE NOT NULL DEFAULT 0,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_meal_logs_user_id (user_id),
  INDEX idx_meal_logs_food_item_id (food_item_id),
  INDEX idx_meal_logs_meal_time (meal_time),
  CONSTRAINT fk_meal_logs_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  CONSTRAINT fk_meal_logs_food_item FOREIGN KEY (food_item_id) REFERENCES food_items(id)
);

CREATE TABLE IF NOT EXISTS body_records (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  record_date DATETIME(6),
  weight_kg DOUBLE NOT NULL DEFAULT 0,
  note TEXT,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_body_records_user_id (user_id),
  INDEX idx_body_records_record_date (record_date),
  CONSTRAINT fk_body_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS goals (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  goal_type VARCHAR(80),
  target_metric VARCHAR(80),
  target_value DOUBLE NOT NULL DEFAULT 0,
  deadline DATETIME(6),
  priority VARCHAR(40),
  status VARCHAR(40) DEFAULT 'active',
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_goals_user_id (user_id),
  CONSTRAINT fk_goals_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS goal_milestones (
  id CHAR(36) PRIMARY KEY,
  goal_id CHAR(36) NOT NULL,
  title VARCHAR(140),
  target_value DOUBLE NOT NULL DEFAULT 0,
  due_date DATETIME(6),
  completed BOOLEAN NOT NULL DEFAULT FALSE,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_goal_milestones_goal_id (goal_id),
  CONSTRAINT fk_goal_milestones_goal FOREIGN KEY (goal_id) REFERENCES goals(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS analytics_snapshots (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  period_type VARCHAR(40),
  start_date DATETIME(6),
  end_date DATETIME(6),
  workout_consistency DOUBLE NOT NULL DEFAULT 0,
  training_volume DOUBLE NOT NULL DEFAULT 0,
  calorie_intake DOUBLE NOT NULL DEFAULT 0,
  calorie_balance DOUBLE NOT NULL DEFAULT 0,
  goal_adherence DOUBLE NOT NULL DEFAULT 0,
  weight_trend VARCHAR(80),
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_analytics_snapshots_user_id (user_id),
  INDEX idx_analytics_snapshots_period_type (period_type),
  INDEX idx_analytics_snapshots_start_date (start_date),
  INDEX idx_analytics_snapshots_end_date (end_date),
  CONSTRAINT fk_analytics_snapshots_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS ai_recommendations (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  type VARCHAR(60),
  prompt_context TEXT,
  content TEXT,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_ai_recommendations_user_id (user_id),
  CONSTRAINT fk_ai_recommendations_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS recommendation_feedbacks (
  id CHAR(36) PRIMARY KEY,
  recommendation_id CHAR(36) NOT NULL,
  rating VARCHAR(40),
  suitability VARCHAR(80),
  comment TEXT,
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_recommendation_feedbacks_recommendation_id (recommendation_id),
  CONSTRAINT fk_recommendation_feedbacks_recommendation FOREIGN KEY (recommendation_id) REFERENCES ai_recommendations(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS progress_reports (
  id CHAR(36) PRIMARY KEY,
  user_id CHAR(36) NOT NULL,
  period_type VARCHAR(40),
  generated_at DATETIME(6),
  summary TEXT,
  file_url VARCHAR(255),
  created_at DATETIME(6) NOT NULL,
  updated_at DATETIME(6) NOT NULL,
  INDEX idx_progress_reports_user_id (user_id),
  CONSTRAINT fk_progress_reports_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
