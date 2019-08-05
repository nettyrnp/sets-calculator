# An algebraic calculator over sets of integers

Features:

* Accepts recursive inputs of any depth
* Valid operations: Union, Intersection, Diff

## Grammar
* set = file | expression
* expression = [ operator file_1 file_2 file_3 ..  file_N ]
* operator = SUM | INT | DIF

## Installation
Install the application in the Terminal:
```shell
go get github.com/nettyrnp/sets-calculator
```

## Running
Before running, all sets of integers should be stored in .txt files in some folder (e.g. "/files").
The flag '-folder' can be used to point to a folder other that the default one "/files".

Go to the folder of the application:
```shell
cd sets-calculator
```

Building:
```shell
make build
```

Running (pass your own expression to the '-task' flag):
```shell
./scalc -task='[ SUM [ DIF a.txt b.txt c.txt ] [ INT b.txt c.txt ] ]'
```

Sample output:
1
3
4
