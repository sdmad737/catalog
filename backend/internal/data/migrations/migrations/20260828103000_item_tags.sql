-- create "item_tags" table
CREATE TABLE `item_tags` (`id` uuid NOT NULL, `created_at` datetime NOT NULL, `updated_at` datetime NOT NULL, `public_token` text NOT NULL, `revoked_at` datetime NULL, `verified_at` datetime NULL, `last_scanned_at` datetime NULL, `item_tags` uuid NOT NULL, PRIMARY KEY (`id`), CONSTRAINT `item_tags_items_tags` FOREIGN KEY (`item_tags`) REFERENCES `items` (`id`) ON DELETE CASCADE);
-- create index "item_tags_public_token_key" to table: "item_tags"
CREATE UNIQUE INDEX `item_tags_public_token_key` ON `item_tags` (`public_token`);
-- create index "itemtag_revoked_at" to table: "item_tags"
CREATE INDEX `itemtag_revoked_at` ON `item_tags` (`revoked_at`);
-- enforce one active tag per item while preserving revoked tag history
CREATE UNIQUE INDEX `itemtag_item_tags` ON `item_tags` (`item_tags`) WHERE `revoked_at` IS NULL;
