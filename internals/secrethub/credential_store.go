package secrethub

import (
	"time"

	"github.com/secrethub/secrethub-cli/internals/cli/ui"
	"github.com/secrethub/secrethub-go/pkg/secrethub/configdir"
	"github.com/secrethub/secrethub-go/pkg/secrethub/credentials"
	"github.com/spf13/cobra"
)

// Errors
var (
	ErrCredentialNotExist = errMain.Code("credential_not_exist").Error("could not find credential file. Run `secrethub signup` to create an account.")
)

// CredentialConfig handles the configuration necessary for local credentials.
type CredentialConfig interface {
	IsPassphraseSet() bool
	Provider() credentials.Provider
	Import() (credentials.Key, error)
	ConfigDir() configdir.Dir
	PassphraseReader() credentials.Reader

	//Register(FlagRegisterer)
}

// NewCredentialConfig creates a new CredentialConfig.
func NewCredentialConfig(io ui.IO) CredentialConfig {
	dir, _ := configdir.Default()
	c := ConfigDir{Dir: *dir}
	return &credentialConfig{
		configDir:                    c,
		CredentialPassphraseCacheTTL: 5 * time.Minute,
		io:                           io,
	}
}

type credentialConfig struct {
	configDir                    ConfigDir
	AccountCredential            string
	credentialPassphrase         string
	CredentialPassphraseCacheTTL time.Duration
	io                           ui.IO
}

func (store *credentialConfig) ConfigDir() configdir.Dir {
	return store.configDir.Dir
}

func (store *credentialConfig) IsPassphraseSet() bool {
	return store.credentialPassphrase != ""
}

func (store *credentialConfig) Register(c *cobra.Command) {
	c.Flags().StringVar(&store.AccountCredential, "credential", "", "Use a specific account credential to authenticate to the API. This overrides the credential stored in the configuration directory.")
	c.Flags().StringVar(&store.credentialPassphrase, "credential-passphrase", "", "The passphrase to unlock your credential file. When set, it will not prompt for the passphrase, nor cache it in the OS keyring. Please only use this if you know what you're doing and ensure your passphrase doesn't end up in bash history.")
	c.Flags().DurationVar(&store.CredentialPassphraseCacheTTL, "credential-passphrase-cache-ttl", 5*time.Minute, "Cache the credential passphrase in the OS keyring for this duration. The cache is automatically cleared after the timer runs out. Each time the passphrase is read from the cache the timer is reset. Passphrase caching is turned on by default for 5 minutes. Turn it off by setting the duration to 0.")
}

// Provider retrieves a credential from the store.
// When a credential is set, that credential is returned,
// otherwise the credential is read from the configured file.
func (store *credentialConfig) Provider() credentials.Provider {
	return credentials.UseKey(store.getCredentialReader()).Passphrase(store.PassphraseReader())
}

func (store *credentialConfig) Import() (credentials.Key, error) {
	return credentials.ImportKey(store.getCredentialReader(), store.PassphraseReader())
}

func (store *credentialConfig) getCredentialReader() credentials.Reader {
	if store.AccountCredential != "" {
		return credentials.FromString(store.AccountCredential)
	}
	return store.configDir.Credential()
}

// PassphraseReader returns a PassphraseReader configured by the flags.
func (store *credentialConfig) PassphraseReader() credentials.Reader {
	return NewPassphraseReader(store.io, store.credentialPassphrase, store.CredentialPassphraseCacheTTL)
}
