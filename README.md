# Delta CLI
Command line application for executing source code inside of a container. 

####Build:
```
go build -o ./delta-cli ./main.go
```

####Usage:
```
./delta-cli <path_to_main_file>
```

Optionally, you can set execution time limit:

```
./delta-cli -t <time_in_ms> <path_to_main_file>
```
and language of the source code provided:
```
./delta-cli -l <lanugage> <path_to_main_file>
```

####Example:
Execute "Hello World" example in Java with execution time limit set to 1500ms:
```
./delta-cli -l java -t 1500 ./Main.java
```
Output:
```
Finished :)
STDOUT:
Hello Java

STDERR:

EXIT CODE: 0
Elapsed exec time: 758ms
Elapsed total time: 991ms
```