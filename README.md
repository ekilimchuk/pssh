# pssh
pssh is a parallel ssh client.
## Install
### Build
```
$ git clone https://github.com/ekilimchuk/pssh.git
$ cd pssh
$ make
```
## Usage
### Help
```
$ ./pssh -h
```
### Running
```
$ ssh-add # running once.
$ ./pssh run -P 2222 -c uptime 127.0.0.1 127.0.0.1
```
