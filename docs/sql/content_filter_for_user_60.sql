-- Run only after the new content-filter tables have been created by AutoMigrate.
-- Test account 60 first; remove its test configuration before enabling 60802.
START TRANSACTION;
INSERT INTO user_content_filter (user_id) VALUES (60)
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);
INSERT INTO user_hidden_tag (user_id, name) VALUES
  (60, 'LGBT'), (60, 'LBGT+'), (60, '*LGBT+'), (60, '梯子'),
  (60, 'VPN'), (60, '机场'), (60, 'FQ'), (60, '翻墙'),
  (60, '*bdsm'), (60, '可能bdsm'), (60, 'zz'), (60, 'xxg'),
  (60, '性相关'), (60, '女菩萨'), (60, '男菩萨'), (60, '代课'),
  (60, '折叠'), (60, '为了折叠'), (60, '理性讨论'),
  (60, '可能引起不适'), (60, '不适'), (60, 'spanking')
ON DUPLICATE KEY UPDATE user_id = VALUES(user_id);
COMMIT;
