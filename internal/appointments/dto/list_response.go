package dto

type PreturnoSupport struct {
	ID          string `json:"id"`
	Nombre      string `json:"nombre"`
	URL         string `json:"url"`
	ViewURL     string `json:"viewUrl,omitempty"`
	DownloadURL string `json:"downloadUrl,omitempty"`
	Tipo        string `json:"tipo,omitempty"`
	Campo       string `json:"campo,omitempty"`
	Tamano      int64  `json:"tamano,omitempty"`
}

type PreturnoTimelineItem struct {
	ID      string `json:"id"`
	Fecha   string `json:"fecha"`
	Titulo  string `json:"titulo"`
	Detalle string `json:"detalle"`
	Usuario string `json:"usuario"`
	Source  string `json:"source,omitempty"`
}

type PreturnoItem struct {
	ID                   string                 `json:"id"`
	Radicado             string                 `json:"radicado"`
	Turno                string                 `json:"turno"`
	Estado               string                 `json:"estado"`
	AutorizacionDatos    bool                   `json:"autorizacionDatos"`
	FechaConsulta        string                 `json:"fechaConsulta"`
	NombreCompleto       string                 `json:"nombreCompleto"`
	TipoDocumento        string                 `json:"tipoDocumento"`
	NumeroDocumento      string                 `json:"numeroDocumento"`
	FechaNacimiento      string                 `json:"fechaNacimiento"`
	Edad                 int                    `json:"edad"`
	EstadoCivil          string                 `json:"estadoCivil"`
	Genero               string                 `json:"genero"`
	Direccion            string                 `json:"direccion"`
	TipoVivienda         string                 `json:"tipoVivienda"`
	Estrato              *int                   `json:"estrato"`
	SisbenCategoria      string                 `json:"sisbenCategoria"`
	Telefono             string                 `json:"telefono"`
	CorreoElectronico    string                 `json:"correoElectronico"`
	TipoPoblacion        string                 `json:"tipoPoblacion"`
	CabezaHogar          *bool                  `json:"cabezaHogar"`
	Ocupacion            string                 `json:"ocupacion"`
	NivelEstudio         string                 `json:"nivelEstudio"`
	Relato               string                 `json:"relato"`
	Soportes             []PreturnoSupport      `json:"soportes"`
	AutorizaNotificacion bool                   `json:"autorizaNotificacion"`
	Timeline             []PreturnoTimelineItem `json:"timeline"`
	Citizen              string                 `json:"citizen"`
	Cedula               string                 `json:"cedula"`
}

type ListPreturnosResponse struct {
	Items []PreturnoItem `json:"items"`
	Total int64          `json:"total"`
	Page  int            `json:"page"`
	Limit int            `json:"limit"`
}
