package common

import (
	cmap "github.com/orcaman/concurrent-map"
)

var AdminProtocols = []string{"http", "https","SSH", "rdp"}
var ServerRoles = []string{"mailbox", "ldap", "mta", "proxy", "archive"}

type ZCSServer struct {
	Name string
	CPUCount int
	Memory int
	IPAddress string
	PreferredConnectionMethod string
	SupportedAdminMethods SupportedAdminMethods
	ZCSInfo ZCSInfo
	OSInfo OSInfo
	ServerRoles []string
	SafeGuards SafeGuards
	VirtualizationInfo VirtualizationInfo
	BackingBlobStorage string
}

type CommandLineFlags struct{
	shortFlag string
	longFlag string
	required bool
}

type SupportedAdminMethods struct{
	ZimbraAdminHTTPS bool
	SSH bool
	RDP bool
}

type SafeGuards struct{
	RecordChanges bool
	BackupChanges bool
}

type VirtualizationInfo struct{
	Virtualized bool
	VirtPlatform string
	VirtVersion string
}

type ZCSInfo struct{
	Version string
	InstallDate string
}

type OSInfo struct{
	Flavor string
	Version string
	InstallDate string
}



type ConnectionServerConfig struct {
	MailboxServers       string
	LdapServers          []string
	ServerTypePreference string
	AdminUsername        string
	AdminPassword        string
	AdminPort            string
	AdminProtocol        string
	SocksServerString    string
	AdminAuthToken       string
	AllServers           cmap.ConcurrentMap
	SkipSSLChecks        bool
	UseSocks5Proxy       bool
	PrintedSSLSkipNotice bool
}

