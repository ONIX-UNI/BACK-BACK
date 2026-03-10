CREATE TABLE IF NOT EXISTS sicou.case_file_otp (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),

  citizen_id      uuid NOT NULL 
                  REFERENCES sicou.citizen(id) 
                  ON DELETE CASCADE,

  case_file_id    uuid 
                  REFERENCES sicou.case_file(id) 
                  ON DELETE CASCADE,

  otp_hash        text NOT NULL,

  attempts        smallint NOT NULL DEFAULT 0,
  max_attempts    smallint NOT NULL DEFAULT 5,

  expires_at      timestamptz NOT NULL,
  used_at         timestamptz,

  created_at      timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_attempts_positive CHECK (attempts >= 0),
  CONSTRAINT ck_max_attempts_positive CHECK (max_attempts > 0)
);

CREATE INDEX idx_case_file_otp_citizen 
ON sicou.case_file_otp(citizen_id);

CREATE INDEX idx_case_file_otp_expires 
ON sicou.case_file_otp(expires_at);

CREATE INDEX idx_case_file_otp_active 
ON sicou.case_file_otp(citizen_id, used_at);

-- CREATE UNIQUE INDEX unique_active_otp
-- ON sicou.case_file_otp (citizen_id)
-- WHERE used_at IS NULL;