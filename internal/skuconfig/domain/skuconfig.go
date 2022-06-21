package skuconfig

import "errors"

// SKUConfig
// Here we are encapsulating properties by declaring them private
// And by defining a NewSKUConfig function, we are enforcing validation and business rules
// For them

type SKUConfig struct {
	uuid          string
	packageName   string
	countryCode   string
	percentileMin uint
	percentileMax uint
	sku           string
}

// Getters for SKUConfig object

func (s SKUConfig) UUID() string {
	return s.uuid
}

func (s SKUConfig) PackageName() string {
	return s.packageName
}

func (s SKUConfig) CountryCode() string {
	return s.countryCode
}

func (s SKUConfig) PercentileMin() uint {
	return s.percentileMin
}

func (s SKUConfig) PercentileMax() uint {
	return s.percentileMax
}

func (s SKUConfig) SKU() string {
	return s.sku
}

// NewSKUConfig returns a new sku config object by applying validation for it's members
func NewSKUConfig(uuid string, packageName string, countryCode string, percentileMin uint, percentileMax uint, sku string) (*SKUConfig, error) {
	if uuid == "" {
		return nil, errors.New("empty skuconfig uuid")
	}

	if packageName == "" {
		return nil, errors.New("empty package name")
	}

	if countryCode == "" {
		return nil, errors.New("empty country code")
	}

	if percentileMin < 0 {
		return nil, errors.New("negative percentile min")
	}

	if percentileMax > 100 {
		return nil, errors.New("percentile max is more than 100")
	}

	if sku == "" {
		return nil, errors.New("empty sku")
	}

	return &SKUConfig{
		uuid:          uuid,
		packageName:   packageName,
		countryCode:   countryCode,
		percentileMin: percentileMin,
		percentileMax: percentileMax,
		sku:           sku,
	}, nil
}

// UnmarshalSKUConfigFromDatabase unmarshals SKUConfig from the database.
// This is intended to use only for unmarshalling and should not be used
// as a constructor to avoid inconsistent state
func UnmarshalSKUConfigFromDatabase(
	uuid string,
	packageName string,
	countryCode string,
	percentileMin uint,
	percentileMax uint,
	sku string,
) (*SKUConfig, error) {
	skuCnf, err := NewSKUConfig(uuid, packageName, countryCode, percentileMin, percentileMax, sku)
	if err != nil {
		return nil, err
	}

	return skuCnf, nil
}
