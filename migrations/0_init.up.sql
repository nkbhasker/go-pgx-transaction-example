CREATE TABLE "public"."users" (
  "id" uuid NOT NULL,
  "first_name" text NULL,
  "last_name" text NULL,
  "email" text NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "public"."teams" (
  "id" uuid NOT NULL,
  "name" text NULL,
  "workspace_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "public"."team_members" (
  "id" uuid NOT NULL,
  "team_id" uuid NOT NULL,
  "user_id" uuid NOT NULL,
  "workspace_id" uuid NOT NULL,
  PRIMARY KEY ("id")
);

CREATE TABLE "public"."workspaces" (
  "id" uuid NOT NULL,
  "name" text NOT NULL,
  "owner" uuid NOT NULL,
  PRIMARY KEY ("id")
);