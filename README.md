# ASN-ROUTE-FETCHER
A script to fetch the routes off an ASN (AS NUMBER)

Usage:

Golang 1.21 or higher

go mod init fetch

go mod tidy

go build fetch.go

./fetch -f asnlist.txt -o routes.txt

-f is used for the file with the list of as numbers the format for this is one per line with "AS000" example:

AS123
AS1234
