-- Disable the enforcement of foreign-keys constraints
PRAGMA foreign_keys = off;
-- Create "new_group_invitation_tokens" table
CREATE TABLE `new_group_invitation_tokens` (`id` uuid NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `token` blob NOT NULL, `expires_at` datetime NOT NULL, `uses` integer NOT NULL DEFAULT (0), `email` text NULL, `role` text NOT NULL DEFAULT ('staff'), `group_invitation_tokens` uuid NULL, PRIMARY KEY (`id`), CONSTRAINT `group_invitation_tokens_groups_invitation_tokens` FOREIGN KEY (`group_invitation_tokens`) REFERENCES `groups` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "group_invitation_tokens" to new temporary table "new_group_invitation_tokens"
INSERT INTO `new_group_invitation_tokens` (`id`, `created_at`, `updated_at`, `token`, `expires_at`, `uses`, `group_invitation_tokens`) SELECT `id`, `created_at`, `updated_at`, `token`, `expires_at`, `uses`, `group_invitation_tokens` FROM `group_invitation_tokens`;
-- Drop "group_invitation_tokens" table after copying rows
DROP TABLE `group_invitation_tokens`;
-- Rename temporary table "new_group_invitation_tokens" to "group_invitation_tokens"
ALTER TABLE `new_group_invitation_tokens` RENAME TO `group_invitation_tokens`;
-- Create index "group_invitation_tokens_token_key" to table: "group_invitation_tokens"
CREATE UNIQUE INDEX `group_invitation_tokens_token_key` ON `group_invitation_tokens` (`token`);
-- Create "new_users" table
CREATE TABLE `new_users` (`id` uuid NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `name` text NOT NULL, `email` text NOT NULL, `password` text NOT NULL, `is_superuser` bool NOT NULL DEFAULT (false), `superuser` bool NOT NULL DEFAULT (false), `role` text NOT NULL DEFAULT ('staff'), `disabled_at` datetime NULL, `activated_on` datetime NULL, `group_users` uuid NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `users_groups_users` FOREIGN KEY (`group_users`) REFERENCES `groups` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "users" to new temporary table "new_users"
INSERT INTO `new_users` (`id`, `created_at`, `updated_at`, `name`, `email`, `password`, `is_superuser`, `superuser`, `role`, `activated_on`, `group_users`) SELECT `id`, `created_at`, `updated_at`, `name`, `email`, `password`, `is_superuser`, `superuser`, CASE WHEN `role` = 'owner' THEN 'owner' ELSE 'staff' END AS `role`, `activated_on`, `group_users` FROM `users`;
-- Drop "users" table after copying rows
DROP TABLE `users`;
-- Rename temporary table "new_users" to "users"
ALTER TABLE `new_users` RENAME TO `users`;
-- Create index "users_email_key" to table: "users"
CREATE UNIQUE INDEX `users_email_key` ON `users` (`email`);
-- Create "new_item_tags" table
CREATE TABLE `new_item_tags` (`id` uuid NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `public_token` text NOT NULL, `revoked_at` datetime NULL, `verified_at` datetime NULL, `last_scanned_at` datetime NULL, `item_tags` uuid NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `item_tags_items_tags` FOREIGN KEY (`item_tags`) REFERENCES `items` (`id`) ON DELETE CASCADE);
-- Copy rows from old table "item_tags" to new temporary table "new_item_tags"
INSERT INTO `new_item_tags` (`id`, `created_at`, `updated_at`, `public_token`, `revoked_at`, `verified_at`, `last_scanned_at`, `item_tags`) SELECT `id`, `created_at`, `updated_at`, `public_token`, `revoked_at`, `verified_at`, `last_scanned_at`, `item_tags` FROM `item_tags`;
-- Drop "item_tags" table after copying rows
DROP TABLE `item_tags`;
-- Rename temporary table "new_item_tags" to "item_tags"
ALTER TABLE `new_item_tags` RENAME TO `item_tags`;
-- Create index "item_tags_public_token_key" to table: "item_tags"
CREATE UNIQUE INDEX `item_tags_public_token_key` ON `item_tags` (`public_token`);
-- Create index "itemtag_revoked_at" to table: "item_tags"
CREATE INDEX `itemtag_revoked_at` ON `item_tags` (`revoked_at`);
-- Create index "itemtag_item_tags" to table: "item_tags"
CREATE UNIQUE INDEX `itemtag_item_tags` ON `item_tags` (`item_tags`) WHERE revoked_at IS NULL;
-- Enable back the enforcement of foreign-keys constraints
PRAGMA foreign_keys = on;
