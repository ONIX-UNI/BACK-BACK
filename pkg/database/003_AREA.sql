ALTER TABLE sicou.service_type
ADD COLUMN IF NOT EXISTS legal_area_id smallint
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

INSERT INTO sicou.legal_area (code, name) VALUES
('laboral', 'Laboral'),
('familia', 'Familia'),
('civil', 'Civil'),
('penal', 'Penal'),
('responsabilidad_fiscal', 'Responsabilidad fiscal'),
('responsabilidad_disciplinaria', 'Responsabilidad disciplinaria')
ON CONFLICT (code) DO NOTHING;

SELECT id, code FROM sicou.legal_area;

INSERT INTO sicou.service_type (code, name) VALUES
('asesorias_juridicas', 'Asesorias juridicas'),
('conciliacion_extrajudicial', 'Conciliacion extrajudicial'),
('representacion_judicial', 'Representacion judicial')
ON CONFLICT (code) DO NOTHING;

SELECT id, code FROM sicou.service_type;

UPDATE sicou.service_type
SET legal_area_id = 1
WHERE code = 'asesorias_juridicas';

UPDATE sicou.service_type
SET legal_area_id = 1
WHERE code = 'conciliacion_extrajudicial';

SELECT id, code, legal_area_id FROM sicou.service_type;

INSERT INTO sicou.service_modality (code, name) VALUES
('INMEDIATA', 'Inmediata'),
('MEDIATA', 'Mediata')
ON CONFLICT (code) DO NOTHING;

SELECT id, code FROM sicou.service_modality;

INSERT INTO sicou.service_variant
(service_type_id, modality_id, code, name)
VALUES
(1, 2, 'oficios', 'Oficios'),
(1, 2, 'memoriales', 'Memoriales'),
(1, 2, 'acciones_tutela', 'Acciones de tutela'),
(1, 2, 'recurso_informes', 'Recurso e informes'),
(1, 2, 'liquidacion_laboral', 'Liquidacion laboral'),
(1, 2, 'derechos_peticion', 'Derechos de peticion'),
(1, 2, 'liquidacion_familia', 'Liquidacion de familia')
ON CONFLICT (service_type_id, modality_id, code) DO NOTHING;
