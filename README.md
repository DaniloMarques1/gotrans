# GoTrans

Transfer files from computer A to computer B. The program have a sender and a
receiver, if you choose a sender you'll need to provide the ip address of the
receiver and the file you wish to transfer. In case you pick a receiver, you
will need to specify where you would like to save the received file.

Both machines needs to be in the same network, when specifying the ip address
it cannot be the loopback address (`localhost`, `127.0.0.1`) it needs to be
something like `192.168.1.5`.

### TODO

* [ ] Only workds with text files, we need to support binary files too
