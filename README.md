# sets-calculator
A recursive calculator over sets of values

## Grammar
// set = file | expression
// expression = [ operator file_1 file_2 file_3 ..  file_N ]
// operator = SUM | INT | DIF

## Sample input:
go run *.go -task='[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]'

## Sample output:
1
3
4
