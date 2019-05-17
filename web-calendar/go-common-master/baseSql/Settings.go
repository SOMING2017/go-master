// Settings
package baseSql

var rootName = "root"
var rootPassword = "mysql2019"
var dbName = "go_master"

var DriverName = "mysql"
var Dsn = rootName + ":" + rootPassword + "@tcp(localhost:3306)/" + dbName + "?charset=utf8&parseTime=true"
