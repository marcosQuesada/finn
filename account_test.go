package finn

import (
	"encoding/json"
	"testing"
)

var raw = `
{
  "data": {
    "type": "accounts",
    "id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
    "organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
    "version": 0,
    "attributes": {
      "country": "GB",
      "base_currency": "GBP",
      "account_number": "41426819",
      "bank_id": "400300",
      "bank_id_code": "GBDSC",
      "bic": "NWBKGB22",
      "iban": "GB11NWBK40030041426819",
      "name": [
        "Samantha Holder"
      ],
      "alternative_names": [
        "Sam Holder"
      ],
      "account_classification": "Personal",
      "joint_account": false,
      "account_matching_opt_out": false,
      "secondary_identification": "A1B2C3D4",
      "switched": false,
      "private_identification": {
        "birth_date": "2017-07-23",
        "birth_country": "GB",
        "identification": "13YH458762",
        "address": "[10 Avenue des Champs]",
        "city": "London",
        "country": "GB"
      },
      "organisation_identification": {
        "identification": "123654",
        "actors": [
          {
            "name": [
              "Jeff Page"
            ],
            "birth_date": "1970-01-01",
            "residency": "GB"
          }
        ],
        "address": ["10 Avenue des Champs"],
        "city": "London",
        "country": "GB"
      },
      "status": "confirmed"
    },
    "relationships": {
      "master_account": {
        "data": [{
          "type": "accounts",
          "id": "a52d13a4-f435-4c00-cfad-f5e7ac5972df"
        }]
      },
      "account_events": {
        "data": [
          {
            "type": "account_events",
            "id": "c1023677-70ee-417a-9a6a-e211241f1e9c"
          },
          {
            "type": "account_events",
            "id": "437284fa-62a6-4f1d-893d-2959c9780288"
          }
        ]
     }
    }
  }
}
`

func TestUnMarshallFullResponseDoesNotThrowError(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}
}

func TestUnMarshallFullResponseCreatesValidAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData == nil {
		t.Fatal("account data is nil")
	}

	if got, want := acc.AccoundData.Type, "accounts"; got != want {
		t.Errorf("account Type does not match, expected %s got %s", want, got)
	}

	if got, want := acc.AccoundData.ID, "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"; got != want {
		t.Errorf("account ID does not match, expected %s got %s", want, got)
	}

	if got, want := acc.AccoundData.OrganisationID, "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c"; got != want {
		t.Errorf("account OrganisationID does not match, expected %s got %s", want, got)
	}

	version := 0
	if acc.AccoundData.Version != version {
		t.Errorf("account version does not match, expected %d got %d", version, acc.AccoundData.Version)
	}
}

func TestUnMarshallFullResponseValidAttributesFromAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData.Attributes == nil {
		t.Fatal("nil attributes")
	}

	attr := acc.AccoundData.Attributes

	if got, want := attr.Country, "GB"; got != want {
		t.Errorf("country attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.BaseCurrency, "GBP"; got != want {
		t.Errorf("base currency attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.AccountNumber, "41426819"; got != want {
		t.Errorf("account numer attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.BankID, "400300"; got != want {
		t.Errorf("bank ID attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.BankIDCode, "GBDSC"; got != want {
		t.Errorf("bank ID code attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.Bic, "NWBKGB22"; got != want {
		t.Errorf("BIC code attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.Iban, "GB11NWBK40030041426819"; got != want {
		t.Errorf("IBAN code attribute does not match, expected %s got %s", want, got)
	}

	if got, want := len(attr.Name), 1; got != want {
		t.Fatalf("Unexpected attribute name size, expected %d got %d", want, got)
	}

	if got, want := attr.Name[0], "Samantha Holder"; got != want {
		t.Errorf("Name attribute does not match, expected %s got %s", want, got)
	}

	if got, want := len(attr.AlternativeNames), 1; got != want {
		t.Fatalf("Unexpected alternative name size, expected %d got %d", want, got)
	}

	if got, want := attr.AlternativeNames[0], "Sam Holder"; got != want {
		t.Errorf("Alternative Name attribute does not match, expected %s got %s", want, got)
	}

	if got, want := attr.AccountClassification, "Personal"; got != want {
		t.Errorf("Accound classification attribute does not match, expected %s got %s", want, got)
	}

	if attr.JointAccount {
		t.Errorf("Unexpected join account, expected false")
	}

	if attr.AccountMatchingOptOut {
		t.Errorf("Unexpected matching opt out account, expected false")
	}

	if got, want := attr.SecondaryID, "A1B2C3D4"; got != want {
		t.Errorf("secondary ID attribute does not match, expected %s got %s", want, got)
	}

	if attr.Switched {
		t.Errorf("Unexpected switched, expected false")
	}

	if got, want := attr.Status, "confirmed"; got != want {
		t.Errorf("status does not match, expected %s got %s", want, got)
	}
}

func TestUnMarshallRawValidPrivateIdentificationAttributesFromAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData.Attributes.PrivateID == nil {
		t.Fatal("nil private ID")
	}

	privateID := acc.AccoundData.Attributes.PrivateID

	if got, want := privateID.BirthDate, "2017-07-23"; got != want {
		t.Errorf("birthday attribute does not match, expected %s got %s", want, got)
	}

	if got, want := privateID.BirthCountry, "GB"; got != want {
		t.Errorf("birth country attribute does not match, expected %s got %s", want, got)
	}

	if got, want := privateID.Identification, "13YH458762"; got != want {
		t.Errorf("identification attribute does not match, expected %s got %s", want, got)
	}

	if got, want := privateID.Address, "[10 Avenue des Champs]"; got != want {
		t.Errorf("address attribute does not match, expected %s got %s", want, got)
	}

	if got, want := privateID.City, "London"; got != want {
		t.Errorf("city attribute does not match, expected %s got %s", want, got)
	}

	if got, want := privateID.Country, "GB"; got != want {
		t.Errorf("country attribute does not match, expected %s got %s", want, got)
	}
}

func TestUnMarshalRawValidOrganisationIDttributesFromAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData.Attributes.OrganisationID == nil {
		t.Fatal("nil organisation ID")
	}

	org := acc.AccoundData.Attributes.OrganisationID

	if got, want := org.Identification, "123654"; got != want {
		t.Errorf("identification attribute does not match, expected %s got %s", want, got)
	}

	if got, want := len(org.Actors), 1; got != want {
		t.Fatalf("unexpected actors size, expected %d got %d", want, got)
	}

	if got, want := len(org.Actors[0].Name), 1; got != want {
		t.Fatalf("unexpected actor name size, expected %d got %d", want, got)
	}

	if got, want := org.Actors[0].Name[0], "Jeff Page"; got != want {
		t.Errorf("actor name attribute does not match, expected %s got %s", want, got)
	}

	if got, want := org.Actors[0].BirthDate, "1970-01-01"; got != want {
		t.Errorf("actor birthday attribute does not match, expected %s got %s", want, got)
	}

	if got, want := org.Actors[0].Residency, "GB"; got != want {
		t.Errorf("actor residency attribute does not match, expected %s got %s", want, got)
	}

	if got, want := len(org.Address), 1; got != want {
		t.Fatalf("unexpected address size, expected %d got %d", want, got)
	}

	if got, want := org.Address[0], "10 Avenue des Champs"; got != want {
		t.Errorf("organisation address id attribute does not match, expected %s got %s", want, got)
	}

	if got, want := org.City, "London"; got != want {
		t.Errorf("organisation city attribute does not match, expected %s got %s", want, got)
	}

	if got, want := org.Country, "GB"; got != want {
		t.Errorf("organisation country attribute does not match, expected %s got %s", want, got)
	}
}

func TestUnMarshalRawValidMasterRelationshipsFromAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData.Relationships == nil {
		t.Fatal("nil relationships")
	}

	if acc.AccoundData.Relationships.Master == nil {
		t.Fatalf("nil master relationship")
	}

	master := acc.AccoundData.Relationships.Master

	if got, want := len(master.Data), 1; got != want {
		t.Errorf("master data lenght relationship does not match, expected %d got %d", want, got)
	}

	if got, want := master.Data[0].Type, "accounts"; got != want {
		t.Errorf("master data type does not match, expected %s got %s", want, got)
	}
	if got, want := master.Data[0].ID, "a52d13a4-f435-4c00-cfad-f5e7ac5972df"; got != want {
		t.Errorf("master data id does not match, expected %s got %s", want, got)
	}
}

func TestUnMarshalRawValidEventsRelationshipsFromAccountData(t *testing.T) {
	acc := &Account{}
	err := json.Unmarshal([]byte(raw), acc)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if acc.AccoundData.Relationships == nil {
		t.Fatal("nil relationships")
	}

	if acc.AccoundData.Relationships.Events == nil {
		t.Fatalf("nil events relationship")
	}

	events := acc.AccoundData.Relationships.Events

	if got, want := len(events.Data), 2; got != want {
		t.Errorf("event data lenght relationship does not match, expected %d got %d", want, got)
	}

	if got, want := events.Data[0].Type, "account_events"; got != want {
		t.Errorf("events data type does not match, expected %s got %s", want, got)
	}
	if got, want := events.Data[0].ID, "c1023677-70ee-417a-9a6a-e211241f1e9c"; got != want {
		t.Errorf("master data id does not match, expected %s got %s", want, got)
	}
}

func TestUnMarshalRawListResponseToAccountsList(t *testing.T) {
	accs := &AccountList{}
	err := json.Unmarshal([]byte(listResponse), accs)
	if err != nil {
		t.Fatalf("unexepected error unmarshalling raw data, error %v", err)
	}

	if got, want := len(accs.Accounts), 2; want != got {
		t.Fatal("unexpected account size")
	}

	if accs.Links == nil || accs.Links.First == "" || accs.Links.Last == "" || accs.Links.Self == "" {
		t.Error("unexpected link content")
	}
}

var listResponse = `
{
	"data": [{
		"attributes": {
			"alternative_bank_account_names": null,
			"base_currency": "EUR",
			"country": "ES"
		},
		"created_on": "2020-05-10T15:48:14.164Z",
		"id": "79f5f1f1-646a-44b9-b73f-4032f94a5bab",
		"modified_on": "2020-05-10T15:48:14.164Z",
		"organisation_id": "d1d31e14-5bce-4a83-9921-9ef65467c253",
		"type": "accounts",
		"version": 0
	}, {
		"attributes": {
			"alternative_bank_account_names": null,
			"base_currency": "EUR",
			"country": "ES"
		},
		"created_on": "2020-05-10T15:48:46.161Z",
		"id": "8c06c8d9-d050-4c19-9d18-8b350cf42619",
		"modified_on": "2020-05-10T15:48:46.161Z",
		"organisation_id": "f0943dc5-94f1-41c4-84bf-a1a27c25b942",
		"type": "accounts",
		"version": 0
	}],
	"links": {
		"first": "/v1/organisation/accounts?page%5Bnumber%5D=first\u0026page%5Bsize%5D=10",
		"last": "/v1/organisation/accounts?page%5Bnumber%5D=last\u0026page%5Bsize%5D=10",
		"self": "/v1/organisation/accounts?page%5Bnumber%5D=0\u0026page%5Bsize%5D=10"
	}
}`
