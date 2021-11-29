package vysagota_libs

import "os"

var JD_ADDRESS string = "http://" + os.Getenv("JD_ADDRESS")
