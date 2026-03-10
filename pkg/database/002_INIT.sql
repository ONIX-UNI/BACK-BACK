-- =========================================================
-- A) ESCRITORIOS
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.escritorio (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  code          text NOT NULL UNIQUE,         
  name          text NOT NULL,                
  location      text,                         
  is_active     boolean NOT NULL DEFAULT true,

  -- puede ser NULL (escritorio vacío)
  current_user_id uuid REFERENCES sicou.app_user(id) ON DELETE SET NULL,
  assigned_at     timestamptz,
  assigned_by     uuid REFERENCES sicou.app_user(id) ON DELETE SET NULL,

  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_escritorio_current_user ON sicou.escritorio(current_user_id);

DROP TRIGGER IF EXISTS trg_escritorio_updated_at ON sicou.escritorio;
CREATE TRIGGER trg_escritorio_updated_at
BEFORE UPDATE ON sicou.escritorio
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- =========================================================
-- B) RELACIÓN TURNO 
-- =========================================================

ALTER TABLE sicou.turn
  ADD COLUMN IF NOT EXISTS escritorio_id uuid REFERENCES sicou.escritorio(id) ON DELETE SET NULL;

CREATE INDEX IF NOT EXISTS idx_turn_escritorio ON sicou.turn(escritorio_id);

-- =========================================================
-- C) AUDITORÍA DE TURNO (timeline de acciones)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.turn_audit (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  turn_id          uuid NOT NULL REFERENCES sicou.turn(id) ON DELETE CASCADE,

  event_type       text NOT NULL,
 
  title            text,
  notes            text,
  payload          jsonb,

  actor_user_id    uuid REFERENCES sicou.app_user(id) ON DELETE SET NULL,
  occurred_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_turn_audit_turn ON sicou.turn_audit(turn_id);
CREATE INDEX IF NOT EXISTS idx_turn_audit_type ON sicou.turn_audit(event_type);

-- Helper: obtener actor desde variable de sesión (opcional)
CREATE OR REPLACE FUNCTION sicou.fn_get_actor_user()
RETURNS uuid
LANGUAGE sql
AS $$
  SELECT NULLIF(current_setting('sicou.actor_user_id', true), '')::uuid;
$$;

-- =========================================================
-- D) FUNCIÓN: crear escritorio vacío
-- =========================================================
CREATE OR REPLACE FUNCTION sicou.fn_escritorio_create_empty(
  p_code text,
  p_name text,
  p_location text DEFAULT NULL
)
RETURNS uuid
LANGUAGE plpgsql
AS $$
DECLARE
  v_id uuid;
BEGIN
  INSERT INTO sicou.escritorio(code, name, location, current_user_id, assigned_at, assigned_by)
  VALUES (p_code, p_name, p_location, NULL, NULL, NULL)
  RETURNING id INTO v_id;

  RETURN v_id;
END;
$$;

-- =========================================================
-- E) FUNCIÓN: asignar usuario a escritorio (si estaba vacío o reasignación simple)
-- =========================================================
CREATE OR REPLACE FUNCTION sicou.fn_escritorio_assign_user(
  p_escritorio_id uuid,
  p_user_id uuid,
  p_reason text DEFAULT NULL
)
RETURNS void
LANGUAGE plpgsql
AS $$
DECLARE
  v_actor uuid := sicou.fn_get_actor_user();
  v_old_user uuid;
BEGIN
  SELECT current_user_id INTO v_old_user
  FROM sicou.escritorio
  WHERE id = p_escritorio_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Escritorio no existe: %', p_escritorio_id;
  END IF;

  UPDATE sicou.escritorio
  SET current_user_id = p_user_id,
      assigned_at = now(),
      assigned_by = v_actor
  WHERE id = p_escritorio_id;

END;
$$;

-- =========================================================
-- F) FUNCIÓN: reasignar escritorio a otro usuario (registra auditoría en turn_audit
--     si hay un turno EN_ATENCION/LLAMADO asociado al escritorio)
-- =========================================================
CREATE OR REPLACE FUNCTION sicou.fn_escritorio_reassign_user(
  p_escritorio_id uuid,
  p_new_user_id uuid,
  p_reason text DEFAULT NULL
)
RETURNS void
LANGUAGE plpgsql
AS $$
DECLARE
  v_actor uuid := sicou.fn_get_actor_user();
  v_old_user uuid;
  v_turn_id uuid;
BEGIN
  SELECT current_user_id INTO v_old_user
  FROM sicou.escritorio
  WHERE id = p_escritorio_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Escritorio no existe: %', p_escritorio_id;
  END IF;

  IF v_old_user IS NOT NULL AND v_old_user = p_new_user_id THEN
    RETURN;
  END IF;

  UPDATE sicou.escritorio
  SET current_user_id = p_new_user_id,
      assigned_at = now(),
      assigned_by = v_actor
  WHERE id = p_escritorio_id;

  -- Si existe un turno en proceso ligado a este escritorio, registra auditoría del turno.
  SELECT t.id INTO v_turn_id
  FROM sicou.turn t
  WHERE t.escritorio_id = p_escritorio_id
    AND t.status IN ('LLAMADO','EN_ATENCION')
  ORDER BY t.created_at DESC
  LIMIT 1;

  IF v_turn_id IS NOT NULL THEN
    INSERT INTO sicou.turn_audit(turn_id, event_type, title, notes, payload, actor_user_id)
    VALUES (
      v_turn_id,
      'DESK_USER_REASSIGNED',
      'Reasignación de escritorio',
      p_reason,
      jsonb_build_object(
        'escritorio_id', p_escritorio_id,
        'old_user_id', v_old_user,
        'new_user_id', p_new_user_id
      ),
      v_actor
    );
  END IF;
END;
$$;

-- =========================================================
-- G) FUNCIÓN: asignar escritorio a un turno (y auditar)
-- =========================================================
CREATE OR REPLACE FUNCTION sicou.fn_turn_set_escritorio(
  p_turn_id uuid,
  p_escritorio_id uuid,
  p_reason text DEFAULT NULL
)
RETURNS void
LANGUAGE plpgsql
AS $$
DECLARE
  v_actor uuid := sicou.fn_get_actor_user();
  v_old_escritorio uuid;
BEGIN
  SELECT escritorio_id INTO v_old_escritorio
  FROM sicou.turn
  WHERE id = p_turn_id
  FOR UPDATE;

  IF NOT FOUND THEN
    RAISE EXCEPTION 'Turno no existe: %', p_turn_id;
  END IF;

  UPDATE sicou.turn
  SET escritorio_id = p_escritorio_id
  WHERE id = p_turn_id;

  INSERT INTO sicou.turn_audit(turn_id, event_type, title, notes, payload, actor_user_id)
  VALUES (
    p_turn_id,
    CASE WHEN v_old_escritorio IS NULL THEN 'DESK_ASSIGNED' ELSE 'DESK_CHANGED' END,
    'Asignación/cambio de escritorio',
    p_reason,
    jsonb_build_object(
      'old_escritorio_id', v_old_escritorio,
      'new_escritorio_id', p_escritorio_id
    ),
    v_actor
  );
END;
$$;
