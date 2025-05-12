package model

type Product struct {
	ID           string
	Name         string
	PackageSizes []int
}

type Package struct {
	PackageUnits []PackageUnit
}

type PackageUnit struct {
	Amount int
	Size   int
}
