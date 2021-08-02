package dvpnconf

// dVPNConfig http request/response structure
type dVPNConfig struct {
	Chain     Chain     `toml:"chain"`
	Handshake Handshake `toml:"handshake"`
	Keyring   Keyring   `toml:"keyring"`
	Node      Node      `toml:"node"`
	Qos       Qos       `toml:"qos"`
}
type Chain struct {
	GasAdjustment      float64 `toml:"gas_adjustment"`
	Gas                int     `toml:"gas"`
	GasPrices          string  `toml:"gas_prices"`
	ID                 string  `toml:"id"`
	RPCAddress         string  `toml:"rpc_address"`
	SimulateAndExecute bool    `toml:"simulate_and_execute"`
}
type Handshake struct {
	Enable bool `toml:"enable"`
	Peers  int  `toml:"peers"`
}
type Keyring struct {
	Backend string `toml:"backend"`
	From    string `toml:"from"`
}
type Node struct {
	IntervalSetSessions    string `toml:"interval_set_sessions"`
	IntervalUpdateSessions string `toml:"interval_update_sessions"`
	IntervalUpdateStatus   string `toml:"interval_update_status"`
	ListenOn               string `toml:"listen_on"`
	Moniker                string `toml:"moniker"`
	Price                  string `toml:"price"`
	Provider               string `toml:"provider"`
	RemoteURL              string `toml:"remote_url"`
}
type Qos struct {
}
