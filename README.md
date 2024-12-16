## Using application
Run:`go run main`

The app will listen to unix socket at: `/tmp/io-rw-app.sock`

You can use socat to interact with socket: `echo -e 'testing' | socat - UNIX-CONNECT:/tmp/io-rw-app.sock`

When a message is sent to the unix socket a random string with be used to create a dirictory `/tmp/io-rw-app/{random string}`. Stdin, stdout and status files will be created at that directory and the value sent to the unix socket will be written to stdin file and copied to the stdout file. Once stdin is successfully written then then status file is written to.
