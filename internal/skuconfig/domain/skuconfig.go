package skuconfig

import "errors"

// SKUConfig
// Here we are encapsulating properties by declaring them private
// And by defining a NewSKUConfig function, we are enforcing validation and business rules
// For them

type SKUConfig struct {
	id            string
	packageName   string
	countryCode   string
	percentileMin uint
	percentileMax uint
	sku           string
}

// Getters for SKUConfig object

func (s SKUConfig) ID() string {
	return s.id
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
func NewSKUConfig(id string, packageName string, countryCode string, percentileMin uint, percentileMax uint, sku string) (*SKUConfig, error) {
	if id == "" {
		return nil, errors.New("empty skuconfig id")
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
		id:            id,
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
	id string,
	packageName string,
	countryCode string,
	percentileMin uint,
	percentileMax uint,
	sku string,
) (*SKUConfig, error) {
	skuCnf, err := NewSKUConfig(id, packageName, countryCode, percentileMin, percentileMax, sku)
	if err != nil {
		return nil, err
	}

	return skuCnf, nil
}
