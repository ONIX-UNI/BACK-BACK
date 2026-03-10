package mail

type SMTPClient struct {
	cfg SMTPConfig
}

func NewSMTPClient(cfg SMTPConfig) (*SMTPClient, error) {
	cfg = withDefaults(cfg)
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return &SMTPClient{cfg: cfg}, nil
}
