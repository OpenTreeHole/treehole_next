-- Run after the application has created user_hidden_division via AutoMigrate.
-- The production division named 圆桌 has ID 2. No other users are affected.
INSERT INTO user_hidden_division (user_id, division_id)
SELECT 60802, id FROM division WHERE id = 2 AND name = '圆桌'
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);
