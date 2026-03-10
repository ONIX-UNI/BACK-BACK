package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func ensurePreturnoTables(ctx context.Context, tx pgx.Tx) error {
	if err := ensurePreturnoCatalogSeeds(ctx, tx); err != nil {
		return err
	}

	_, err := tx.Exec(ctx, `
		ALTER TABLE sicou.preturno
		ADD COLUMN IF NOT EXISTS tracking_seq bigserial,
		ADD COLUMN IF NOT EXISTS assigned_coordinator_id uuid REFERENCES sicou.app_user(id),
		ADD COLUMN IF NOT EXISTS service_type_id smallint REFERENCES sicou.service_type(id),
		ADD COLUMN IF NOT EXISTS assigned_at timestamptz
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		ALTER TABLE sicou.preturno
		ALTER COLUMN stratum_snapshot DROP NOT NULL,
		ALTER COLUMN head_of_household DROP NOT NULL
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE UNIQUE INDEX IF NOT EXISTS uq_preturno_tracking_seq
		ON sicou.preturno(tracking_seq)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_preturno_status_created
		ON sicou.preturno(status, created_at)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_preturno_assigned_coordinator
		ON sicou.preturno(assigned_coordinator_id)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_preturno_service_type
		ON sicou.preturno(service_type_id)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS sicou.preturno_timeline (
			id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
			preturno_id uuid NOT NULL REFERENCES sicou.preturno(id) ON DELETE CASCADE,
			event_type text NOT NULL,
			title text NOT NULL,
			detail text,
			source text NOT NULL DEFAULT 'lista_interna',
			created_by uuid REFERENCES sicou.app_user(id) ON DELETE SET NULL,
			created_at timestamptz NOT NULL DEFAULT now()
		)
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		ALTER TABLE sicou.preturno_timeline
		ADD COLUMN IF NOT EXISTS source text
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		UPDATE sicou.preturno_timeline
		SET source = 'lista_interna'
		WHERE source IS NULL OR BTRIM(source) = ''
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		ALTER TABLE sicou.preturno_timeline
		ALTER COLUMN source SET DEFAULT 'lista_interna',
		ALTER COLUMN source SET NOT NULL
	`)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, `
		CREATE INDEX IF NOT EXISTS idx_preturno_timeline_preturno
		ON sicou.preturno_timeline(preturno_id, created_at)
	`)
	return err
}

func ensurePreturnoCatalogSeeds(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Exec(ctx, `
		INSERT INTO sicou.catalog_civil_status(code, name) VALUES
			('SOLTERO', 'Soltero'),
			('CASADO', 'Casado'),
			('UNION_LIBRE', 'Union libre'),
			('SEPARADO', 'Separado'),
			('DIVORCIADO', 'Divorciado'),
			('VIUDO', 'Viudo'),
			('OTRO', 'Otro')
		ON CONFLICT (code) DO NOTHING;

		INSERT INTO sicou.catalog_gender(code, name) VALUES
			('M', 'Hombre'),
			('F', 'Mujer'),
			('NB', 'No binario'),
			('OTRO', 'Otro')
		ON CONFLICT (code) DO NOTHING;

		INSERT INTO sicou.catalog_housing_type(code, name) VALUES
			('ARRENDADA', 'Arrendada'),
			('PROPIA', 'Propia'),
			('FAMILIAR', 'Familiar'),
			('OTRO', 'Otro')
		ON CONFLICT (code) DO NOTHING;

		INSERT INTO sicou.catalog_population_type(code, name) VALUES
			('VICTIMA_CONFLICTO', 'Victima del conflicto'),
			('VULNERABLE', 'Poblacion vulnerable'),
			('DISCAPACIDAD', 'Persona con discapacidad'),
			('ESPECIAL_PROTECCION', 'Especial proteccion'),
			('OTRO', 'Otro')
		ON CONFLICT (code) DO NOTHING;

		INSERT INTO sicou.catalog_education_level(code, name) VALUES
			('NINGUNO', 'Ninguno'),
			('PRIMARIA', 'Primaria'),
			('SECUNDARIA', 'Secundaria'),
			('TECNICO', 'Tecnico'),
			('TECNOLOGO', 'Tecnologo'),
			('UNIVERSITARIO', 'Universitario'),
			('POSGRADO', 'Posgrado'),
			('OTRO', 'Otro')
		ON CONFLICT (code) DO NOTHING;
	`)
	return err
}
