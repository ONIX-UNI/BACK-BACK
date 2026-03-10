CREATE TABLE IF NOT EXISTS sicou.folder (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  name          text NOT NULL,
  description   text,
  created_at    timestamptz NOT NULL DEFAULT now(),
  deleted_at    timestamptz
);

CREATE TABLE IF NOT EXISTS sicou.folder_document (
  folder_id   uuid NOT NULL REFERENCES sicou.folder(id) ON DELETE CASCADE,
  document_id uuid NOT NULL REFERENCES sicou.document(id) ON DELETE CASCADE,
  added_at    timestamptz NOT NULL DEFAULT now(),
  added_by    uuid REFERENCES sicou.app_user(id),

  PRIMARY KEY (folder_id, document_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS ux_folder_name
ON sicou.folder (lower(name))
WHERE deleted_at IS NULL;