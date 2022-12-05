-- +goose Up
-- +goose StatementBegin
-- Squences
CREATE SEQUENCE IF NOT EXISTS btc_info_id_seq;

-- Table Definition
CREATE TABLE IF NOT EXISTS "public"."btc_info" (
    "id" int4 NOT NULL DEFAULT nextval('btc_info_id_seq'::regclass),
    "amount" numeric,
    "created" timestamptz,
    "offsettz" numeric
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE "btc_info";
DROP SEQUENCE "btc_info_id_seq";
-- +goose StatementEnd
