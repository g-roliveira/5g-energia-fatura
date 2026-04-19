module github.com/gustavo/5g-energia-fatura/packages/normalizer

go 1.25.0

require (
	github.com/gustavo/5g-energia-fatura/packages/calc-engine v0.0.0-00010101000000-000000000000
	github.com/shopspring/decimal v1.4.0
)

replace github.com/gustavo/5g-energia-fatura/packages/calc-engine => ../calc-engine
