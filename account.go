package finn

// Account wraps account data
type Account struct {
	AccoundData *AccoundData `json:"data"`
}

// AccoundData defines user account
type AccoundData struct {
	Type           string         `json:"type"`
	ID             string         `json:"id"`
	OrganisationID string         `json:"organisation_id"`
	Version        int            `json:"version"`
	Attributes     *Attributes    `json:"attributes"`
	Relationships  *Relationships `json:"relationships"`
}

// Attributes defines user account attributes
type Attributes struct {
	Country               string                      `json:"country"`
	BaseCurrency          string                      `json:"base_currency,omitempty"`
	AccountNumber         string                      `json:"account_number,omitempty" `
	BankID                string                      `json:"bank_id,omitempty"`
	BankIDCode            string                      `json:"bank_id_code,omitempty"`
	Bic                   string                      `json:"bic,omitempty"`
	Iban                  string                      `json:"iban,omitempty"`
	Name                  []string                    `json:"name,omitempty"`
	AlternativeNames      []string                    `json:"alternative_names,omitempty"`
	AccountClassification string                      `json:"account_classification,omitempty"`
	JointAccount          bool                        `json:"joint_account,omitempty"`
	AccountMatchingOptOut bool                        `json:"account_matching_opt_out,omitempty"`
	SecondaryID           string                      `json:"secondary_identification,omitempty"`
	Switched              bool                        `json:"switched,omitempty"`
	PrivateID             *PrivateIdentification      `json:"private_identification,omitempty"`
	OrganisationID        *OrganisationIdentification `json:"organisation_identification,omitempty"`
	Status                string                      `json:"status,omitempty"`
}

// PrivateIdentification defines account owner details
type PrivateIdentification struct {
	BirthDate      string `json:"birth_date"`
	BirthCountry   string `json:"birth_country"`
	Identification string `json:"identification"`
	Address        string `json:"address"`
	City           string `json:"city"`
	Country        string `json:"country"`
}

// OrganisationIdentification defines organisation details
type OrganisationIdentification struct {
	Identification string    `json:"identification"`
	Actors         []*Actors `json:"actors"`
	Address        []string  `json:"address"`
	City           string    `json:"city"`
	Country        string    `json:"country"`
}

// Actors declares account actors
type Actors struct {
	Name      []string `json:"name"`
	BirthDate string   `json:"birth_date"`
	Residency string   `json:"residency"`
}

// Relationships holds a slice of AccountEvents
type Relationships struct {
	Master *MasterAccount `json:"master_account"`
	Events *AccountEvents `json:"account_events"`
}

// MasterAccount defines principal account
type MasterAccount struct {
	Data []*Data `json:"data"`
}

// AccountEvents defines associated account events
type AccountEvents struct {
	Data []*Data `json:"data"`
}

// Data detail type on master or events account
type Data struct {
	Type string `json:"type"`
	ID   string `json:"id"`
}

// AccountList returns a list of accounts
type AccountList struct {
	Accounts []*AccoundData `json:"data"`
	Links    *LinkList      `json:"links"`
}

// LinkList enables pagination result access as formal REST
type LinkList struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Self  string `json:"self"`
}
