package main

import (
	"regexp"
)

const DBHOST = "localhost"
const DBPORT = 5432
const DBUSER = "admin"
const DBPASS = "123"
const DBNAME = "xdmedia"

const BASE_URL = "https://localhost"
const SERVER_PORT = 5555
const LOGGER_PORT = 5556
const MIN_ALIAS_LEN = 3
const MAX_ALIAS_LEN = 32
const MIN_USERNAME_LEN = 3
const MAX_USERNAME_LEN = 16
const MIN_PASSWORD_LEN = 8
const MAX_PASSWORD_LEN = 32
const MAX_POST_LEN = 1024

var VALID_CHARS = regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
