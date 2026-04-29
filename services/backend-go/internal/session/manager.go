package session

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/security"
	"github.com/gustavo/5g-energia-fatura/services/backend-go/internal/store"
)

type Manager struct {
	store     store.IntegrationStore
	cipher    *security.Cipher
	bootstrap BootstrapRunner
}

type CredentialInput struct {
	Label      string `json:"label"`
	Documento  string `json:"documento"`
	Senha      string `json:"senha"`
	UF         string `json:"uf"`
	TipoAcesso string `json:"tipo_acesso"`
}

type CredentialView struct {
	ID         string `json:"id"`
	Label      string `json:"label"`
	Documento  string `json:"documento"`
	UF         string `json:"uf"`
	TipoAcesso string `json:"tipo_acesso"`
	CreatedAt  string `json:"created_at"`
}

type SessionView struct {
	ID           string `json:"id"`
	CredentialID string `json:"credential_id"`
	CreatedAt    string `json:"created_at"`
}

type ResolvedSession struct {
	Documento string
	Token     string
}

func NewManager(store store.IntegrationStore, cipher *security.Cipher, bootstrap BootstrapRunner) *Manager {
	return &Manager{store: store, cipher: cipher, bootstrap: bootstrap}
}

func maskDocumento(value string) string {
	digits := strings.TrimSpace(value)
	if len(digits) <= 4 {
		return "***"
	}
	return strings.Repeat("*", len(digits)-4) + digits[len(digits)-4:]
}

func (m *Manager) CreateCredential(ctx context.Context, in CredentialInput) (CredentialView, error) {
	documentoCipher, documentoNonce, err := m.cipher.Encrypt(in.Documento)
	if err != nil {
		return CredentialView{}, err
	}
	senhaCipher, senhaNonce, err := m.cipher.Encrypt(in.Senha)
	if err != nil {
		return CredentialView{}, err
	}
	rec, err := m.store.InsertCredential(store.CredentialRecord{
		Label:           in.Label,
		DocumentoCipher: documentoCipher,
		DocumentoNonce:  documentoNonce,
		SenhaCipher:     senhaCipher,
		SenhaNonce:      senhaNonce,
		UF:              defaultString(in.UF, "BA"),
		TipoAcesso:      defaultString(in.TipoAcesso, "normal"),
		KeyVersion:      "v1",
	})
	if err != nil {
		return CredentialView{}, err
	}
	return CredentialView{
		ID:         rec.ID,
		Label:      rec.Label,
		Documento:  maskDocumento(in.Documento),
		UF:         rec.UF,
		TipoAcesso: rec.TipoAcesso,
		CreatedAt:  rec.CreatedAt,
	}, nil
}

func (m *Manager) CreateSessionFromCredential(ctx context.Context, credentialID string) (SessionView, ResolvedSession, error) {
	cred, err := m.store.GetCredentialByID(credentialID)
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}
	documento, err := m.cipher.Decrypt(cred.DocumentoCipher, cred.DocumentoNonce)
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}
	senha, err := m.cipher.Decrypt(cred.SenhaCipher, cred.SenhaNonce)
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}

	boot, err := m.bootstrap.Run(ctx, BootstrapInput{
		Documento:  documento,
		Senha:      senha,
		UF:         cred.UF,
		TipoAcesso: cred.TipoAcesso,
	})
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}

	tokenCipher, tokenNonce, err := m.cipher.Encrypt(boot.Token)
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}

	sess, err := m.store.InsertSession(store.SessionRecord{
		CredentialID:      cred.ID,
		BearerTokenCipher: tokenCipher,
		BearerTokenNonce:  tokenNonce,
	})
	if err != nil {
		return SessionView{}, ResolvedSession{}, err
	}

	return SessionView{
			ID:           sess.ID,
			CredentialID: sess.CredentialID,
			CreatedAt:    sess.CreatedAt,
		}, ResolvedSession{
			Documento: documento,
			Token:     boot.Token,
		}, nil
}

func (m *Manager) ResolveToken(ctx context.Context, credentialID string) (ResolvedSession, error) {
	cred, err := m.store.GetCredentialByID(credentialID)
	if err != nil {
		return ResolvedSession{}, err
	}
	documento, err := m.cipher.Decrypt(cred.DocumentoCipher, cred.DocumentoNonce)
	if err != nil {
		return ResolvedSession{}, err
	}

	sess, err := m.store.GetLatestSessionByCredentialID(credentialID)
	if err == nil {
		token, decErr := m.cipher.Decrypt(sess.BearerTokenCipher, sess.BearerTokenNonce)
		if decErr == nil {
			return ResolvedSession{Documento: documento, Token: token}, nil
		}
	}
	if err != nil && err != sql.ErrNoRows {
		return ResolvedSession{}, err
	}

	_, resolved, err := m.CreateSessionFromCredential(ctx, credentialID)
	if err != nil {
		return ResolvedSession{}, err
	}
	return resolved, nil
}

func defaultString(value string, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

func (m *Manager) RequireCredentialFields(in CredentialInput) error {
	if strings.TrimSpace(in.Documento) == "" || strings.TrimSpace(in.Senha) == "" {
		return fmt.Errorf("documento e senha são obrigatórios")
	}
	return nil
}
