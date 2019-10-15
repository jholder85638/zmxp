package common


type GeneralSettings struct{
	RunWithoutUI bool
	RunWithoutUIDescription string

	ZmXMaxThreads int
	ZmXMaxThreadsDescription string

	ZmXRestoreSession bool
	ZmXRestoreSessionDescription string

	ZmXUITheme string
	ZmXUIThemeDescription string

	ZmXInstallPath string
	ZmXInstallPathDescription string

	ZmXSettingsPath string
	ZmXSettingsPathDescription string

	ZmXEncryptSettings bool
	ZmXEncryptSettingsDescription string

	ZmXEncryptCloudRecovery bool
	ZmXEncryptCloudRecoveryDescription string

	UpdateURL string
	UpdateURLDescription string

	KillSwitchURL string
	KillSwitchURLDescription string

	KillSwitchActive bool
	KillSwitchActiveDescription string

	ZmcontrolCloudSyncActive bool
	ZmcontrolCloudSyncActiveDescription string

	ZmcontrolCloudSyncKey string
	ZmcontrolCloudSyncKeyDescription string

	LoggerEnabled bool
	LoggerEnabledDescription string

	LoggerLogToFile bool
	LoggerLogToFileDescription string

	ChangeManagementSaveChangeHistory bool
	ChangeManagementSaveChangeHistoryDescription string

	ChangeManagementSaveBackup bool
	ChangeManagementSaveBackupDescription string
}

type SshSettings struct {
	PrivateKeyFile string
	PrivateKeyFileDescription string

	PrivateKey string
	PrivateKeyDescription string

	AllowRunAsRoot bool
	AllowRunAsRootDescription string

	AllowPasswordlessKey bool
	AllowPasswordlessKeyDescription string
}

type ZcsAdminSettings struct {
	AdminUserName string
	AdminUserNameDescription string

	AuthMailboxServer string
	AuthMailboxServerDescription string

	AuthProtocol string
	AuthProtocolDescription string

	AuthAdminPort int
	AuthAdminPortDescription string

	AllowRunAsZimbra bool
	AllowRunAsZimbraDescription string

	AllowSaveAdminUsername bool
	AllowSaveAdminUsernameDescription string

	AllowSaveAdminPassword bool
	AllowSaveAdminPasswordDescription string
}

type NetworkSettings struct{
	UseSocksProxy bool
	UseSocksProxyDescription string

	Socks5URL string
	Socks5URLDescription string

	SkipSSLChecks bool
	SkipSSLChecksDescription string

	ZmcontrolCloudProxyEnabled bool
	ZmcontrolCloudProxyEnabledDescription string
}
