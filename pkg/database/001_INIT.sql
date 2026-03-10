-- =========================================================
-- SICOU / Consultorio Jurídico UNIMETA - DB Fase 1 (MVP)
-- PostgreSQL 14+ (recomendado 16+)
-- =========================================================

BEGIN;

-- ---------- Extensions ----------
CREATE EXTENSION IF NOT EXISTS pgcrypto;
CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS unaccent;

-- ---------- Schema ----------
CREATE SCHEMA IF NOT EXISTS sicou;
SET search_path TO sicou, public;

-- ---------- Common columns helpers ----------
CREATE OR REPLACE FUNCTION sicou.fn_set_updated_at()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  NEW.updated_at = now();
  RETURN NEW;
END;
$$;

-- Turno consecutivo por fecha (concurrencia segura)
CREATE TABLE IF NOT EXISTS sicou.turn_counter (
  turn_date date PRIMARY KEY,
  last_seq  integer NOT NULL DEFAULT 0
);

CREATE OR REPLACE FUNCTION sicou.fn_next_turn_consecutive(p_date date)
RETURNS integer
LANGUAGE plpgsql
AS $$
DECLARE v_next integer;
BEGIN
  INSERT INTO sicou.turn_counter(turn_date, last_seq)
  VALUES (p_date, 1)
  ON CONFLICT (turn_date)
  DO UPDATE SET last_seq = sicou.turn_counter.last_seq + 1
  RETURNING last_seq INTO v_next;

  RETURN v_next;
END;
$$;

-- Calendario de no laborables (para días hábiles)
CREATE TABLE IF NOT EXISTS sicou.calendar_holiday (
  holiday_date date PRIMARY KEY,
  description  text
);

CREATE OR REPLACE FUNCTION sicou.fn_add_business_days(p_start date, p_days int)
RETURNS date
LANGUAGE plpgsql
AS $$
DECLARE
  d date := p_start;
  added int := 0;
BEGIN
  IF p_days IS NULL OR p_days <= 0 THEN
    RETURN d;
  END IF;

  WHILE added < p_days LOOP
    d := d + INTERVAL '1 day';
    -- 0=domingo, 6=sábado
    IF EXTRACT(DOW FROM d) NOT IN (0,6)
       AND NOT EXISTS (SELECT 1 FROM sicou.calendar_holiday h WHERE h.holiday_date = d) THEN
      added := added + 1;
    END IF;
  END LOOP;

  RETURN d;
END;
$$;

-- =========================================================
-- 1) CATÁLOGOS (parametrizables)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.catalog_document_type (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- CC, CE, PAS, PPT, OTRO
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_gender (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- M, F, NB, OTRO
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_civil_status (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- SOLTERO, CASADO, UNION_LIBRE, etc
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_housing_type (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- ARRENDADA, PROPIA, FAMILIAR, OTRO
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_population_type (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- VICTIMA, VULNERABLE, DISCAPACIDAD, ESPECIAL_PROTECCION, OTRO
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_education_level (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- PRIMARIA, SECUNDARIA, TECNICO, TECNOLOGO, UNIVERSITARIO, OTRO
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.catalog_channel (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,  -- PRESENCIAL, VIRTUAL
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.legal_area (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.service_type (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE, -- ASESORIA_INMEDIATA, ASESORIA_MEDIATA, REPRESENTACION
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);
-- TODO: INSERTAR EN LA TABLA
ALTER TABLE sicou.service_type
ADD COLUMN legal_area_id smallint
REFERENCES sicou.legal_area(id);

CREATE TABLE IF NOT EXISTS sicou.service_modality (
  id smallserial PRIMARY KEY,
  code text NOT NULL UNIQUE, -- INMEDIATA, MEDIATA
  name text NOT NULL,
  is_active boolean NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS sicou.service_variant (
  id smallserial PRIMARY KEY,
  service_type_id smallint NOT NULL REFERENCES sicou.service_type(id),
  modality_id smallint NOT NULL REFERENCES sicou.service_modality(id),
  code text NOT NULL,
  name text NOT NULL,
  is_active boolean NOT NULL DEFAULT true,
  created_at timestamptz NOT NULL DEFAULT now(),
  updated_at timestamptz NOT NULL DEFAULT now(),

  UNIQUE(service_type_id, modality_id, code)
);

CREATE TABLE IF NOT EXISTS sicou.non_competence_reason (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE IF NOT EXISTS sicou.derivation_entity (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE,
  name        text NOT NULL,
  contact_info text,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Tipos de trámite para SLA (días hábiles parametrizables)
CREATE TABLE IF NOT EXISTS sicou.procedure_type (
  id                 smallserial PRIMARY KEY,
  code               text NOT NULL UNIQUE, -- LIQUIDACION, TUTELA, MEMORIAL, DERECHO_PETICION, OTRO
  name               text NOT NULL,
  sla_business_days  int  NOT NULL DEFAULT 0,
  alert_48h          boolean NOT NULL DEFAULT true,
  alert_24h          boolean NOT NULL DEFAULT true,
  is_active          boolean NOT NULL DEFAULT true,
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ck_sla_nonneg CHECK (sla_business_days >= 0)
);

-- Tipos de documento / metadatos (expediente + preturno)
CREATE TABLE IF NOT EXISTS sicou.document_kind (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE, -- ID_DOC, FACTURA_SERVICIO, FORMATO_ASESORIA, CONSTANCIA_ENVIO, etc.
  name        text NOT NULL,
  is_active   boolean NOT NULL DEFAULT true,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

-- Updated_at triggers for catalogs
DO $$
DECLARE r record;
BEGIN
  FOR r IN
    SELECT unnest(ARRAY[
      'catalog_document_type','catalog_gender','catalog_civil_status','catalog_housing_type',
      'catalog_population_type','catalog_education_level','catalog_channel',
      'legal_area','service_type','non_competence_reason','derivation_entity',
      'procedure_type','document_kind'
    ]) AS t
  LOOP
    EXECUTE format('
      DROP TRIGGER IF EXISTS trg_%1$s_updated_at ON sicou.%1$s;
      CREATE TRIGGER trg_%1$s_updated_at
      BEFORE UPDATE ON sicou.%1$s
      FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();', r.t);
  END LOOP;
END$$;

-- =========================================================
-- 2) USUARIOS INTERNOS + RBAC
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.app_user (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  email         citext NOT NULL UNIQUE,
  display_name  text NOT NULL,
  password_hash text NOT NULL,
  is_active     boolean NOT NULL DEFAULT true,
  created_at    timestamptz NOT NULL DEFAULT now(),
  updated_at    timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_app_user_updated_at
BEFORE UPDATE ON sicou.app_user
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE TABLE IF NOT EXISTS sicou.role (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE, -- SUPER_ADMIN, ADMIN_CONSULTORIO, SECRETARIA, COORDINADOR, ESTUDIANTE, DOCENTE
  name        text NOT NULL,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_role_updated_at
BEFORE UPDATE ON sicou.role
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE TABLE IF NOT EXISTS sicou.permission (
  id          smallserial PRIMARY KEY,
  code        text NOT NULL UNIQUE, -- e.g. CITIZEN_CREATE, PRETURNO_VALIDATE, CASE_ASSIGN, etc.
  description text,
  created_at  timestamptz NOT NULL DEFAULT now(),
  updated_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_permission_updated_at
BEFORE UPDATE ON sicou.permission
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE TABLE IF NOT EXISTS sicou.user_role (
  user_id uuid NOT NULL REFERENCES sicou.app_user(id) ON DELETE CASCADE,
  role_id smallint NOT NULL REFERENCES sicou.role(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE IF NOT EXISTS sicou.role_permission (
  role_id smallint NOT NULL REFERENCES sicou.role(id) ON DELETE CASCADE,
  permission_id smallint NOT NULL REFERENCES sicou.permission(id) ON DELETE CASCADE,
  created_at timestamptz NOT NULL DEFAULT now(),
  PRIMARY KEY (role_id, permission_id)
);

-- =========================================================
-- 3) CIUDADANO + CARACTERIZACIÓN SOCIOECONÓMICA
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.citizen (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  document_type_id   smallint NOT NULL REFERENCES sicou.catalog_document_type(id),
  document_number    text NOT NULL,
  full_name          text NOT NULL,
  birth_date         date,
  phone_mobile       text,
  email              citext,
  address            text,
  created_by         uuid REFERENCES sicou.app_user(id),
  updated_by         uuid REFERENCES sicou.app_user(id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  deleted_at         timestamptz,
  CONSTRAINT uq_citizen_document UNIQUE (document_type_id, document_number)
);

CREATE TRIGGER trg_citizen_updated_at
BEFORE UPDATE ON sicou.citizen
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- Caracterización (estrato / sisben) SIN bloquear pre-turno (se guarda aunque esté pendiente)
CREATE TABLE IF NOT EXISTS sicou.citizen_socioeconomic (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id         uuid NOT NULL UNIQUE REFERENCES sicou.citizen(id) ON DELETE CASCADE,
  housing_type_id    smallint REFERENCES sicou.catalog_housing_type(id),
  stratum            smallint,
  sisben_category    text,
  sisben_score       numeric(6,2),
  verification_status text NOT NULL DEFAULT 'PENDIENTE', -- PENDIENTE, VERIFICADO, NO_APORTA
  observation        text,
  support_document_id uuid, -- opcional (se relaciona a document más abajo si quieres)
  created_by         uuid REFERENCES sicou.app_user(id),
  updated_by         uuid REFERENCES sicou.app_user(id),
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ck_stratum_range CHECK (stratum IS NULL OR (stratum BETWEEN 1 AND 6))
);

CREATE TRIGGER trg_citizen_socioeconomic_updated_at
BEFORE UPDATE ON sicou.citizen_socioeconomic
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- Elegibilidad/beneficiario (criterio + excepción)
CREATE TABLE IF NOT EXISTS sicou.citizen_eligibility (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id      uuid NOT NULL REFERENCES sicou.citizen(id) ON DELETE CASCADE,
  criterion       text NOT NULL, -- vulnerable / especial protección / etc (puedes convertirlo a catálogo si quieres)
  is_eligible     boolean NOT NULL DEFAULT false,
  support_document_id uuid,
  observation     text,
  exception_authorized boolean NOT NULL DEFAULT false,
  exception_authorized_by uuid REFERENCES sicou.app_user(id),
  exception_authorized_at timestamptz,
  created_by      uuid REFERENCES sicou.app_user(id),
  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_citizen_eligibility_updated_at
BEFORE UPDATE ON sicou.citizen_eligibility
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- =========================================================
-- 4) ARCHIVOS + DOCUMENTOS (unificado para preturno/caso)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.file_object (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  storage_key   text NOT NULL UNIQUE,   -- ruta/clave en storage (S3/MinIO/local)
  original_name text NOT NULL,
  mime_type     text,
  size_bytes    bigint,
  sha256        text,
  created_at    timestamptz NOT NULL DEFAULT now()
);

-- Documento: puede pertenecer a un preturno o a un caso (o ambos via evento)
CREATE TABLE IF NOT EXISTS sicou.document (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  file_id         uuid NOT NULL REFERENCES sicou.file_object(id) ON DELETE RESTRICT,
  document_kind_id smallint NOT NULL REFERENCES sicou.document_kind(id),
  preturno_id     uuid, -- FK más abajo
  case_id         uuid, -- FK más abajo
  case_event_id   uuid, -- FK más abajo
  version         int NOT NULL DEFAULT 1,
  is_current      boolean NOT NULL DEFAULT true,
  notes           text,
  uploaded_by     uuid REFERENCES sicou.app_user(id),
  uploaded_at     timestamptz NOT NULL DEFAULT now(),
  deleted_at      timestamptz,
  CONSTRAINT ck_document_owner CHECK (preturno_id IS NOT NULL OR case_id IS NOT NULL)
);

-- =========================================================
-- 5) PRE-TURNO (TRIAJE) + VALIDACIÓN/DERIVACIÓN
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.preturno (
  id                   uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id           uuid NOT NULL REFERENCES sicou.citizen(id) ON DELETE RESTRICT,

  -- Campos obligatorios del formulario inicial (snapshot)
  contact_email        citext NOT NULL,     -- “Correo (campo de contacto)”
  data_consent_text    text NOT NULL,
  data_consent_accepted boolean NOT NULL DEFAULT false,
  consultation_at      timestamptz NOT NULL, -- fecha de consulta

  full_name_snapshot   text NOT NULL,
  document_type_id_snapshot smallint NOT NULL REFERENCES sicou.catalog_document_type(id),
  document_number_snapshot  text NOT NULL,

  birth_date_snapshot  date,
  age_years_snapshot   smallint,
  civil_status_id      smallint REFERENCES sicou.catalog_civil_status(id),
  gender_id            smallint REFERENCES sicou.catalog_gender(id),
  address_snapshot     text NOT NULL,
  housing_type_id      smallint REFERENCES sicou.catalog_housing_type(id),
  stratum_snapshot     smallint NOT NULL,
  sisben_category_snapshot text,
  phone_mobile_snapshot text NOT NULL,
  email_snapshot       citext NOT NULL,      -- correo del ciudadano
  population_type_id   smallint REFERENCES sicou.catalog_population_type(id),
  head_of_household    boolean NOT NULL DEFAULT false,
  occupation_snapshot  text,
  education_level_id   smallint REFERENCES sicou.catalog_education_level(id),
  situation_story      text NOT NULL,        -- relato breve

  notify_by_email_consent boolean NOT NULL DEFAULT false,

  channel_id           smallint NOT NULL REFERENCES sicou.catalog_channel(id), -- presencial/virtual
  status               text NOT NULL DEFAULT 'PENDIENTE_VALIDACION', -- PENDIENTE_VALIDACION, VALIDADO, NO_COMPETENCIA, EN_TURNO, CERRADO
  observations         text,

  created_by           uuid REFERENCES sicou.app_user(id),
  updated_by           uuid REFERENCES sicou.app_user(id),
  created_at           timestamptz NOT NULL DEFAULT now(),
  updated_at           timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_stratum_snapshot CHECK (stratum_snapshot BETWEEN 1 AND 6)
);

CREATE TRIGGER trg_preturno_updated_at
BEFORE UPDATE ON sicou.preturno
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- Ahora sí: relacionamos document.preturno_id a preturno
ALTER TABLE sicou.document
  ADD CONSTRAINT fk_document_preturno
  FOREIGN KEY (preturno_id) REFERENCES sicou.preturno(id) ON DELETE SET NULL;

-- Validación / clasificación por Coordinador
CREATE TABLE IF NOT EXISTS sicou.preturno_review (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  preturno_id     uuid NOT NULL UNIQUE REFERENCES sicou.preturno(id) ON DELETE CASCADE,
  reviewed_by     uuid NOT NULL REFERENCES sicou.app_user(id),
  reviewed_at     timestamptz NOT NULL DEFAULT now(),

  is_competent    boolean NOT NULL,
  service_type_id smallint REFERENCES sicou.service_type(id),
  legal_area_id   smallint REFERENCES sicou.legal_area(id),

  requirements_notes text,
  result_notes    text
);

-- No competencia / derivación
CREATE TABLE IF NOT EXISTS sicou.preturno_derivation (
  id                 uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  preturno_id        uuid NOT NULL UNIQUE REFERENCES sicou.preturno(id) ON DELETE CASCADE,
  reason_id          smallint NOT NULL REFERENCES sicou.non_competence_reason(id),
  entity_id          smallint NOT NULL REFERENCES sicou.derivation_entity(id),
  channel_id         smallint NOT NULL REFERENCES sicou.catalog_channel(id),
  oriented_at        timestamptz NOT NULL DEFAULT now(),
  evidence_document_id uuid REFERENCES sicou.document(id),
  observations       text NOT NULL
);

-- =========================================================
--=============== 6) TURNO + REPARTO (PROVISIONAL / DEFINITIVO)
-- ==========================================

CREATE TABLE IF NOT EXISTS sicou.turn (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  preturno_id   uuid NOT NULL UNIQUE REFERENCES sicou.preturno(id) ON DELETE RESTRICT,

  turn_date     date NOT NULL DEFAULT (now()::date),
  consecutive   int  NOT NULL,
  status        text NOT NULL DEFAULT 'EN_COLA', -- EN_COLA, LLAMADO, EN_ATENCION, FINALIZADO, ANULADO
  priority      smallint NOT NULL DEFAULT 0,

  created_by    uuid REFERENCES sicou.app_user(id),
  created_at    timestamptz NOT NULL DEFAULT now(),

  called_at     timestamptz,
  attended_at   timestamptz,
  finished_at   timestamptz,

  CONSTRAINT uq_turn_date_consecutive UNIQUE (turn_date, consecutive),
  CONSTRAINT ck_turn_priority CHECK (priority >= 0)
);

CREATE OR REPLACE FUNCTION sicou.fn_turn_before_insert()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.turn_date IS NULL THEN
    NEW.turn_date := (now()::date);
  END IF;

  IF NEW.consecutive IS NULL OR NEW.consecutive <= 0 THEN
    NEW.consecutive := sicou.fn_next_turn_consecutive(NEW.turn_date);
  END IF;

  RETURN NEW;
END;
$$;

DROP TRIGGER IF EXISTS trg_turn_before_insert ON sicou.turn;
CREATE TRIGGER trg_turn_before_insert
BEFORE INSERT ON sicou.turn
FOR EACH ROW EXECUTE FUNCTION sicou.fn_turn_before_insert();

-- Repartos / asignaciones (auditoría completa)
CREATE TABLE IF NOT EXISTS sicou.turn_assignment (
  id             uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  turn_id         uuid NOT NULL REFERENCES sicou.turn(id) ON DELETE CASCADE,
  stage           text NOT NULL, -- PROVISIONAL / DEFINITIVO
  method          text NOT NULL, -- AUTO / MANUAL
  assigned_to     uuid NOT NULL REFERENCES sicou.app_user(id),
  assigned_by     uuid NOT NULL REFERENCES sicou.app_user(id),
  reason          text,
  assigned_at     timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ck_stage CHECK (stage IN ('PROVISIONAL','DEFINITIVO')),
  CONSTRAINT ck_method CHECK (method IN ('AUTO','MANUAL'))
);

CREATE INDEX IF NOT EXISTS idx_turn_assignment_turn ON sicou.turn_assignment(turn_id);
CREATE INDEX IF NOT EXISTS idx_turn_assignment_assigned_to ON sicou.turn_assignment(assigned_to);

-- =========================================================
-- 7) CASO / EXPEDIENTE + HISTORIAL
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.case_file (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  citizen_id       uuid NOT NULL REFERENCES sicou.citizen(id) ON DELETE RESTRICT,
  preturno_id      uuid NOT NULL UNIQUE REFERENCES sicou.preturno(id) ON DELETE RESTRICT,
  turn_id          uuid UNIQUE REFERENCES sicou.turn(id) ON DELETE SET NULL,

  service_type_id  smallint NOT NULL REFERENCES sicou.service_type(id),
  legal_area_id    smallint REFERENCES sicou.legal_area(id),

  status           text NOT NULL DEFAULT 'ABIERTO', -- ABIERTO, PENDIENTE_DOCUMENTOS, EN_TRAMITE, DESCARGADO, CERRADO
  current_responsible uuid REFERENCES sicou.app_user(id),
  supervisor_user     uuid REFERENCES sicou.app_user(id),

  opened_at        timestamptz NOT NULL DEFAULT now(),
  closed_at        timestamptz,
  close_notes      text,

  created_by       uuid REFERENCES sicou.app_user(id),
  updated_by       uuid REFERENCES sicou.app_user(id),
  created_at       timestamptz NOT NULL DEFAULT now(),
  updated_at       timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_case_status CHECK (status IN ('ABIERTO','PENDIENTE_DOCUMENTOS','EN_TRAMITE','DESCARGADO','CERRADO'))
);

CREATE TRIGGER trg_case_file_updated_at
BEFORE UPDATE ON sicou.case_file
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

-- Eventos / historial (línea de tiempo)
CREATE TABLE IF NOT EXISTS sicou.case_event (
  id           uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  case_id      uuid NOT NULL REFERENCES sicou.case_file(id) ON DELETE CASCADE,
  event_type   text NOT NULL, -- e.g. CREATED, STATUS_CHANGED, NOTE, DOCUMENT_ADDED, ASSIGNMENT, etc.
  title        text,
  notes        text,
  payload      jsonb,
  created_by   uuid REFERENCES sicou.app_user(id),
  created_at   timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_case_event_case ON sicou.case_event(case_id);
CREATE INDEX IF NOT EXISTS idx_case_event_type ON sicou.case_event(event_type);

-- document.case_id / document.case_event_id
ALTER TABLE sicou.document
  ADD CONSTRAINT fk_document_case
  FOREIGN KEY (case_id) REFERENCES sicou.case_file(id) ON DELETE SET NULL;

ALTER TABLE sicou.document
  ADD CONSTRAINT fk_document_case_event
  FOREIGN KEY (case_event_id) REFERENCES sicou.case_event(id) ON DELETE SET NULL;

-- Un documento actual por tipo en un caso
CREATE UNIQUE INDEX IF NOT EXISTS uq_case_doc_kind_current
ON sicou.document(case_id, document_kind_id)
WHERE case_id IS NOT NULL AND is_current = true AND deleted_at IS NULL;

-- =========================================================
-- 8) CHECKLIST DESCARGO / CIERRE + ENCUESTA
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.checklist_item_def (
  id                 smallserial PRIMARY KEY,
  code               text NOT NULL UNIQUE,
  name               text NOT NULL,
  description        text,
  requires_document_kind_id smallint REFERENCES sicou.document_kind(id),
  required_for_close boolean NOT NULL DEFAULT false,
  is_active          boolean NOT NULL DEFAULT true,
  created_at         timestamptz NOT NULL DEFAULT now(),
  updated_at         timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_checklist_item_def_updated_at
BEFORE UPDATE ON sicou.checklist_item_def
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE TABLE IF NOT EXISTS sicou.service_type_checklist (
  service_type_id smallint NOT NULL REFERENCES sicou.service_type(id) ON DELETE CASCADE,
  checklist_item_id smallint NOT NULL REFERENCES sicou.checklist_item_def(id) ON DELETE CASCADE,
  PRIMARY KEY (service_type_id, checklist_item_id)
);

CREATE TABLE IF NOT EXISTS sicou.case_checklist_item (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  case_id         uuid NOT NULL REFERENCES sicou.case_file(id) ON DELETE CASCADE,
  checklist_item_id smallint NOT NULL REFERENCES sicou.checklist_item_def(id),
  is_done         boolean NOT NULL DEFAULT false,
  done_at         timestamptz,
  done_by         uuid REFERENCES sicou.app_user(id),
  evidence_document_id uuid REFERENCES sicou.document(id),
  notes           text,
  created_at      timestamptz NOT NULL DEFAULT now(),
  UNIQUE (case_id, checklist_item_id)
);

-- Encuesta de satisfacción (condiciona cierre o excepción)
CREATE TABLE IF NOT EXISTS sicou.case_satisfaction_survey (
  case_id      uuid PRIMARY KEY REFERENCES sicou.case_file(id) ON DELETE CASCADE,
  channel_id   smallint REFERENCES sicou.catalog_channel(id),
  score        smallint NOT NULL, -- 1..5
  comments     text,
  created_at   timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT ck_score CHECK (score BETWEEN 1 AND 5)
);

CREATE TABLE IF NOT EXISTS sicou.case_satisfaction_exception (
  case_id        uuid PRIMARY KEY REFERENCES sicou.case_file(id) ON DELETE CASCADE,
  reason         text NOT NULL,
  authorized_by  uuid NOT NULL REFERENCES sicou.app_user(id),
  authorized_at  timestamptz NOT NULL DEFAULT now()
);

-- =========================================================
-- 9) TÉRMINOS / VENCIMIENTOS (SLA días hábiles) + ALERTAS
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.case_term (
  id               uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  case_id           uuid NOT NULL REFERENCES sicou.case_file(id) ON DELETE CASCADE,
  procedure_type_id smallint NOT NULL REFERENCES sicou.procedure_type(id),
  start_date        date NOT NULL DEFAULT (now()::date),
  due_date          date NOT NULL,
  status            text NOT NULL DEFAULT 'ATIEMPO', -- ATIEMPO, PROXIMO, VENCIDO, CUMPLIDO, CANCELADO
  completed_at      timestamptz,
  completed_by      uuid REFERENCES sicou.app_user(id),
  notes             text,

  alert_48_sent_at  timestamptz,
  alert_24_sent_at  timestamptz,

  created_at        timestamptz NOT NULL DEFAULT now(),
  updated_at        timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_term_status CHECK (status IN ('ATIEMPO','PROXIMO','VENCIDO','CUMPLIDO','CANCELADO'))
);

CREATE TRIGGER trg_case_term_updated_at
BEFORE UPDATE ON sicou.case_term
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE INDEX IF NOT EXISTS idx_case_term_case ON sicou.case_term(case_id);
CREATE INDEX IF NOT EXISTS idx_case_term_due ON sicou.case_term(due_date);

-- =========================================================
-- 10) NOTIFICACIONES POR CORREO (plantillas + bitácora / outbox)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.email_template (
  code         text PRIMARY KEY, -- ASIGNACION, VENCIMIENTO, NO_COMPETENCIA, ENTREGA, PQRS, etc.
  subject_tpl  text NOT NULL,
  body_tpl     text NOT NULL,
  cc_emails    text[] NOT NULL DEFAULT '{}',
  is_active    boolean NOT NULL DEFAULT true,
  created_at   timestamptz NOT NULL DEFAULT now(),
  updated_at   timestamptz NOT NULL DEFAULT now()
);

CREATE TRIGGER trg_email_template_updated_at
BEFORE UPDATE ON sicou.email_template
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE TABLE IF NOT EXISTS sicou.email_outbox (
  id            uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  template_code text REFERENCES sicou.email_template(code),
  to_emails      text[] NOT NULL,
  cc_emails      text[] NOT NULL DEFAULT '{}',
  subject        text NOT NULL,
  body           text NOT NULL,
  status         text NOT NULL DEFAULT 'PENDIENTE', -- PENDIENTE, ENVIANDO, ENVIADO, FALLIDO
  attempts       int NOT NULL DEFAULT 0,
  last_error     text,
  scheduled_at   timestamptz NOT NULL DEFAULT now(),
  sent_at        timestamptz,

  related_case_id uuid REFERENCES sicou.case_file(id),
  related_turn_id uuid REFERENCES sicou.turn(id),
  related_term_id uuid REFERENCES sicou.case_term(id),
  related_pqrs_id uuid, -- FK más abajo

  created_at     timestamptz NOT NULL DEFAULT now(),
  updated_at     timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_email_status CHECK (status IN ('PENDIENTE','ENVIANDO','ENVIADO','FALLIDO'))
);

CREATE TRIGGER trg_email_outbox_updated_at
BEFORE UPDATE ON sicou.email_outbox
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

CREATE INDEX IF NOT EXISTS idx_email_outbox_status ON sicou.email_outbox(status, scheduled_at);

-- =========================================================
-- 11) PQRS (correo + registro)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.pqrs (
  id              uuid PRIMARY KEY DEFAULT gen_random_uuid(),
  radicado        bigserial UNIQUE, -- radicado interno
  citizen_id      uuid REFERENCES sicou.citizen(id) ON DELETE SET NULL,
  from_email      citext,
  subject         text NOT NULL,
  body            text NOT NULL,
  received_at     timestamptz NOT NULL DEFAULT now(),
  status          text NOT NULL DEFAULT 'RADICADA', -- RADICADA, EN_GESTION, RESPONDIDA, CERRADA
  assigned_to     uuid REFERENCES sicou.app_user(id),
  response_text   text,
  responded_at    timestamptz,
  evidence_document_id uuid REFERENCES sicou.document(id),

  email_message_id text,   -- id del correo (si lo integras)
  email_thread_id  text,

  created_at      timestamptz NOT NULL DEFAULT now(),
  updated_at      timestamptz NOT NULL DEFAULT now(),

  CONSTRAINT ck_pqrs_status CHECK (status IN ('RADICADA','EN_GESTION','RESPONDIDA','CERRADA'))
);

CREATE TRIGGER trg_pqrs_updated_at
BEFORE UPDATE ON sicou.pqrs
FOR EACH ROW EXECUTE FUNCTION sicou.fn_set_updated_at();

ALTER TABLE sicou.email_outbox
  ADD CONSTRAINT fk_email_outbox_pqrs
  FOREIGN KEY (related_pqrs_id) REFERENCES sicou.pqrs(id) ON DELETE SET NULL;

-- =========================================================
-- 12) AUDITORÍA (bitácora)
-- =========================================================

CREATE TABLE IF NOT EXISTS sicou.audit_log (
  id            bigserial PRIMARY KEY,
  occurred_at   timestamptz NOT NULL DEFAULT now(),
  actor_user_id uuid REFERENCES sicou.app_user(id),
  action        text NOT NULL,            -- INSERT / UPDATE / DELETE
  table_name    text NOT NULL,
  row_id        uuid,
  before_data   jsonb,
  after_data    jsonb,
  ip_address    inet,
  user_agent    text
);

-- Evitar updates/deletes sobre la bitácora
CREATE OR REPLACE FUNCTION sicou.fn_audit_log_immutable()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'audit_log es inmutable';
END;
$$;

DROP TRIGGER IF EXISTS trg_audit_log_no_update ON sicou.audit_log;
CREATE TRIGGER trg_audit_log_no_update
BEFORE UPDATE OR DELETE ON sicou.audit_log
FOR EACH ROW EXECUTE FUNCTION sicou.fn_audit_log_immutable();

CREATE OR REPLACE FUNCTION sicou.fn_audit_row()
RETURNS trigger
LANGUAGE plpgsql
AS $$
DECLARE
  v_row_id uuid;
BEGIN
  IF TG_OP = 'INSERT' THEN
    v_row_id := NEW.id;
    INSERT INTO sicou.audit_log(action, table_name, row_id, after_data)
    VALUES ('INSERT', TG_TABLE_NAME, v_row_id, to_jsonb(NEW));
    RETURN NEW;
  ELSIF TG_OP = 'UPDATE' THEN
    v_row_id := NEW.id;
    INSERT INTO sicou.audit_log(action, table_name, row_id, before_data, after_data)
    VALUES ('UPDATE', TG_TABLE_NAME, v_row_id, to_jsonb(OLD), to_jsonb(NEW));
    RETURN NEW;
  ELSIF TG_OP = 'DELETE' THEN
    v_row_id := OLD.id;
    INSERT INTO sicou.audit_log(action, table_name, row_id, before_data)
    VALUES ('DELETE', TG_TABLE_NAME, v_row_id, to_jsonb(OLD));
    RETURN OLD;
  END IF;
  RETURN NULL;
END;
$$;

-- Activar auditoría en tablas clave (puedes añadir/quitar)
DO $$
DECLARE t text;
BEGIN
  FOREACH t IN ARRAY ARRAY[
    'citizen','citizen_socioeconomic','citizen_eligibility',
    'preturno','preturno_review','preturno_derivation',
    'turn','turn_assignment',
    'case_file','case_event',
    'document','case_term',
    'pqrs','email_outbox'
  ]
  LOOP
    EXECUTE format('DROP TRIGGER IF EXISTS trg_%1$s_audit ON sicou.%1$s;', t);
    EXECUTE format('
      CREATE TRIGGER trg_%1$s_audit
      AFTER INSERT OR UPDATE OR DELETE ON sicou.%1$s
      FOR EACH ROW EXECUTE FUNCTION sicou.fn_audit_row();', t);
  END LOOP;
END$$;

-- =========================================================
-- 13) SEEDS mínimos (roles + catálogos base)
-- =========================================================

INSERT INTO sicou.role(code, name) VALUES
  ('SUPER_ADMIN','Super Administrador (TI)'),
  ('ADMIN_CONSULTORIO','Administrador Consultorio (Jefatura)'),
  ('SECRETARIA','Secretaría / Administrativo'),
  ('COORDINADOR','Coordinador de turno'),
  ('ESTUDIANTE','Estudiante (Asesor jurídico)'),
  ('DOCENTE','Docente supervisor')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_channel(code, name) VALUES
  ('PRESENCIAL','Presencial'),
  ('VIRTUAL','Virtual')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_document_type(code, name) VALUES
  ('CC','Cédula de ciudadanía'),
  ('CE','Cédula de extranjería'),
  ('PAS','Pasaporte'),
  ('PPT','Permiso por Protección Temporal'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.service_type(code, name) VALUES
  ('ASESORIA_INMEDIATA','Asesoría inmediata'),
  ('ASESORIA_MEDIATA','Asesoría mediata'),
  ('REPRESENTACION','Representación')
ON CONFLICT (code) DO NOTHING;

-- SLA ejemplo (parametrizable): 3, 3, 2, 0 según documento
INSERT INTO sicou.procedure_type(code, name, sla_business_days) VALUES
  ('LIQUIDACION','Liquidaciones (laboral/alimentos) u otros cálculos', 3),
  ('TUTELA','Tutelas / impugnación / desacato', 3),
  ('MEMORIAL','Memoriales / recursos / informes', 2),
  ('DERECHO_PETICION','Derecho de petición (mismo día)', 0)
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.document_kind(code, name) VALUES
  ('Soporte','soporte de documento.')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_civil_status(code, name) VALUES
  ('SOLTERO','Soltero'),
  ('CASADO','Casado'),
  ('UNION_LIBRE','Union libre'),
  ('SEPARADO','Separado'),
  ('DIVORCIADO','Divorciado'),
  ('VIUDO','Viudo'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_gender(code, name) VALUES
  ('M','Hombre'),
  ('F','Mujer'),
  ('NB','No binario'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_housing_type(code, name) VALUES
  ('ARRENDADA','Arrendada'),
  ('PROPIA','Propia'),
  ('FAMILIAR','Familiar'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_population_type(code, name) VALUES
  ('VICTIMA_CONFLICTO','Victima del conflicto'),
  ('VULNERABLE','Poblacion vulnerable'),
  ('DISCAPACIDAD','Persona con discapacidad'),
  ('ESPECIAL_PROTECCION','Especial proteccion'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

INSERT INTO sicou.catalog_education_level(code, name) VALUES
  ('NINGUNO','Ninguno'),
  ('PRIMARIA','Primaria'),
  ('SECUNDARIA','Secundaria'),
  ('TECNICO','Tecnico'),
  ('TECNOLOGO','Tecnologo'),
  ('UNIVERSITARIO','Universitario'),
  ('POSGRADO','Posgrado'),
  ('OTRO','Otro')
ON CONFLICT (code) DO NOTHING;

COMMIT;
