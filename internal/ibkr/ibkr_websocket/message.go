package ibkr_websocket

type Message struct {
	Topic string      `json:"topic"`
	Args  interface{} `json:"args,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type StsMessage struct {
	Topic string `json:"topic"`
	Args  struct {
		Authenticated bool   `json:"authenticated"`
		Competing     bool   `json:"competing"`
		Message       string `json:"message"`
		Fail          string `json:"fail"`
		ServerName    string `json:"serverName"`
		ServerVersion string `json:"serverVersion"`
		Username      string `json:"username"`
	} `json:"args"`
}

type SbdMessage struct {
	Topic string `json:"topic"`
	Data  []struct {
		Row   int64   `json:"row"`
		Focus int64   `json:"focus"`
		Price float64 `json:"price"`
		Ask   int64   `json:"ask,omitempty"`
		Bid   int64   `json:"bid,omitempty"`
	} `json:"data"`
}

type SorArgs struct {
	Acct               string  `json:"acct"`
	Conidex            string  `json:"conidex"`
	Conid              int64   `json:"conid"`
	Account            string  `json:"account"`
	OrderID            int64   `json:"orderId"`
	CashCcy            string  `json:"cashCcy"`
	SizeAndFills       string  `json:"sizeAndFills"`
	OrderDesc          string  `json:"orderDesc"`
	Description1       string  `json:"description1"`
	Ticker             string  `json:"ticker"`
	SecType            string  `json:"secType"`
	ListingExchange    string  `json:"listingExchange"`
	RemainingQuantity  float64 `json:"remainingQuantity"`
	FilledQuantity     float64 `json:"filledQuantity"`
	TotalSize          float64 `json:"totalSize"`
	CompanyName        string  `json:"companyName"`
	Status             string  `json:"status"`
	OrderCCPStatus     string  `json:"order_ccp_status"`
	AvgPrice           string  `json:"avgPrice"`
	OrigOrderType      string  `json:"origOrderType"`
	SupportsTaxOpt     string  `json:"supportsTaxOpt"`
	LastExecutionTime  string  `json:"lastExecutionTime"`
	OrderType          string  `json:"orderType"`
	BgColor            string  `json:"bgColor"`
	FgColor            string  `json:"fgColor"`
	OrderRef           string  `json:"order_ref"`
	IsEventTrading     string  `json:"isEventTrading"`
	Price              string  `json:"price"`
	TimeInForce        string  `json:"timeInForce"`
	LastExecutionTimeR int64   `json:"lastExecutionTime_r"`
	Side               string  `json:"side"`
}

type SorMessage struct {
	Topic string    `json:"topic"`
	Args  []SorArgs `json:"args"`
}

type StrMessage struct {
	Topic string `json:"topic"`
	Args  []struct {
		ExecutionID          string  `json:"execution_id"`
		Symbol               string  `json:"symbol"`
		SupportsTaxOpt       string  `json:"supports_tax_opt"`
		Side                 string  `json:"side"`
		OrderDescription     string  `json:"order_description"`
		TradeTime            string  `json:"trade_time"`
		TradeTimeR           int64   `json:"trade_time_r"`
		Size                 int64   `json:"size"`
		Price                string  `json:"price"`
		OrderRef             string  `json:"order_ref"`
		Submitter            string  `json:"submitter"`
		Exchange             string  `json:"exchange"`
		Commission           string  `json:"commission"`
		NetAmount            float64 `json:"net_amount"`
		Account              string  `json:"account"`
		AccountCode          string  `json:"accountCode"`
		CompanyName          string  `json:"company_name"`
		ContractDescription1 string  `json:"contract_description_1"`
		SecType              string  `json:"sec_type"`
		ListingExchange      string  `json:"listing_exchange"`
		Conid                string  `json:"conid"`
		Conidex              string  `json:"conidex"`
		ClearingID           string  `json:"clearing_id"`
		ClearingName         string  `json:"clearing_name"`
		LiquidationTrade     string  `json:"liquidation_trade"`
	} `json:"args"`
}

type SplMessage struct {
	Topic string `json:"topic"`
	Args  struct {
		DU1234Core struct {
			RowType int64   `json:"rowType"`
			Dpl     float64 `json:"dpl"`
			Upl     float64 `json:"upl"`
		} `json:"DU1234.Core"`
	} `json:"args"`
}

type ActMessage struct {
	Topic string `json:"topic"`
	Args  struct {
		Accounts  []string `json:"accounts"`
		AcctProps struct {
			All struct {
				HasChildAccounts  bool `json:"hasChildAccounts"`
				SupportsCashQty   bool `json:"supportsCashQty"`
				NoFXConv          bool `json:"noFXConv"`
				IsProp            bool `json:"isProp"`
				SupportsFractions bool `json:"supportsFractions"`
				AllowCustomerTime bool `json:"allowCustomerTime"`
			} `json:"All"`
			DU1234567 struct {
				HasChildAccounts  bool `json:"hasChildAccounts"`
				SupportsCashQty   bool `json:"supportsCashQty"`
				NoFXConv          bool `json:"noFXConv"`
				IsProp            bool `json:"isProp"`
				SupportsFractions bool `json:"supportsFractions"`
				AllowCustomerTime bool `json:"allowCustomerTime"`
			} `json:"DU1234567"`
		} `json:"acctProps"`
		Aliases struct {
			All       string `json:"All"`
			DU1234567 string `json:"DU1234567"`
		} `json:"aliases"`
		AllowFeatures struct {
			ShowGFIS               bool   `json:"showGFIS"`
			ShowEUCostReport       bool   `json:"showEUCostReport"`
			AllowEventContract     bool   `json:"allowEventContract"`
			AllowFXConv            bool   `json:"allowFXConv"`
			AllowFinancialLens     bool   `json:"allowFinancialLens"`
			AllowMTA               bool   `json:"allowMTA"`
			AllowTypeAhead         bool   `json:"allowTypeAhead"`
			AllowEventTrading      bool   `json:"allowEventTrading"`
			SnapshotRefreshTimeout int64  `json:"snapshotRefreshTimeout"`
			LiteUser               bool   `json:"liteUser"`
			ShowWebNews            bool   `json:"showWebNews"`
			Research               bool   `json:"research"`
			DebugPnl               bool   `json:"debugPnl"`
			ShowTaxOpt             bool   `json:"showTaxOpt"`
			ShowImpactDashboard    bool   `json:"showImpactDashboard"`
			AllowDynAccount        bool   `json:"allowDynAccount"`
			AllowCrypto            bool   `json:"allowCrypto"`
			AllowedAssetTypes      string `json:"allowedAssetTypes"`
		} `json:"allowFeatures"`
		ChartPeriods struct {
			STK    []string `json:"STK"`
			CFD    []string `json:"CFD"`
			OPT    []string `json:"OPT"`
			FOP    []string `json:"FOP"`
			WAR    []string `json:"WAR"`
			IOPT   []string `json:"IOPT"`
			FUT    []string `json:"FUT"`
			CASH   []string `json:"CASH"`
			IND    []string `json:"IND"`
			BOND   []string `json:"BOND"`
			FUND   []string `json:"FUND"`
			CMDTY  []string `json:"CMDTY"`
			PHYSS  []string `json:"PHYSS"`
			CRYPTO []string `json:"CRYPTO"`
		} `json:"chartPeriods"`
		Groups          []string `json:"groups"`
		Profiles        []string `json:"profiles"`
		SelectedAccount string   `json:"selectedAccount"`
		ServerInfo      struct {
			ServerName    string `json:"serverName"`
			ServerVersion string `json:"serverVersion"`
		} `json:"serverInfo"`
		SessionID string `json:"sessionId"`
		IsFT      bool   `json:"isFT"`
		IsPaper   bool   `json:"isPaper"`
	} `json:"args"`
}
