# Mysql dialect

Implements the Dialect interface for a MySQL database, using
the following sql driver:

	github.com/go-sql-driver/mysql

## Error translation

Because the driver only passes the numeric error codes, we had to redefine
all the error constants in errors.go.
This is done automatically from the mysql headers by running the following
command in the current directory:

    go generate

The script "tools/generrors.go" is doing the actual job of writing the source
file.
It requires the headers to be installed, and was last executed with
libmysqlclient-dev 5.7.16-0ubuntu0.16.04.1 on a ubuntu 16.04.
