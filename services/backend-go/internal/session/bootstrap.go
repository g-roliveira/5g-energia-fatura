package session

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
)

type BootstrapRunner struct {
	PythonBin  string
	ScriptPath string
}

type BootstrapInput struct {
	Documento  string
	Senha      string
	UF         string
	TipoAcesso string
}

type BootstrapOutput struct {
	Token        string            `json:"token"`
	TokenNeSe    map[string]string `json:"token_ne_se"`
	LocalStorage map[string]string `json:"local_storage"`
}

func (r BootstrapRunner) Run(ctx context.Context, in BootstrapInput) (BootstrapOutput, error) {
	cmd := exec.CommandContext(
		ctx,
		r.PythonBin,
		r.ScriptPath,
		"--documento", in.Documento,
		"--senha", in.Senha,
		"--uf", in.UF,
		"--tipo-acesso", in.TipoAcesso,
	)
	raw, err := cmd.Output()
	if err != nil {
		return BootstrapOutput{}, fmt.Errorf("bootstrap runner failed: %w", err)
	}

	var out BootstrapOutput
	if err := json.Unmarshal(raw, &out); err != nil {
		return BootstrapOutput{}, fmt.Errorf("bootstrap output inválido: %w", err)
	}
	if out.Token == "" {
		return BootstrapOutput{}, fmt.Errorf("bootstrap sem bearer token")
	}
	return out, nil
}
