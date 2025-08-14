-- Добавляем новые типы действий
INSERT INTO action_types (name) VALUES
  ('click_cta_top'),
  ('click_cta_bottom'),
  ('scroll_depth'),
  ('scroll_section_view'),
  ('external_link_garage_raid'),
  ('external_link_raids'),
  ('external_link_details'),
  ('gallery_scroll_right'),
  ('gallery_scroll_left'),
  ('external_link_mentor_page'),
  ('external_link_social'),
  ('faq_open_answer'),
  ('click_links_buy_access'),
  ('click_links_telegram'),
  ('click_links_youtube_entertainment'),
  ('click_links_youtube_streams'),
  ('click_links_instagram'),
  ('click_links_tiktok'),
  ('click_links_artstation'),
  ('click_links_3d_guide');

-- Добавляем поле meta для дополнительных данных
ALTER TABLE actions
ADD COLUMN meta JSONB NULL;