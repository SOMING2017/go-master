// sql
package calendar

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"

	"../baseSql"
)

var tableName = "calendar_tag"

type TagInfo struct {
	Cid        string    `json:"cid"`
	Tag        string    `json:"tag"`
	Content    string    `json:"content"`
	CreateDate time.Time `json:"createDate"`
}

func GetCalendarNoticeInfoSql(startDatetime time.Time, endIndex *int) (string, error) {
	errorIndex := 0
	db, err := sql.Open(baseSql.DriverName, baseSql.Dsn)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer db.Close()
	endDatetime := time.Date(startDatetime.Year(), startDatetime.Month(), startDatetime.Day()+*endIndex-1,
		startDatetime.Hour(), startDatetime.Minute(), startDatetime.Second(), startDatetime.Nanosecond(),
		startDatetime.Location())
	sqlCommand := "SELECT createFormatDate From " + tableName + " WHERE isDiscard = 0 AND createFormatDate BETWEEN ? AND ? GROUP BY createFormatDate ORDER BY createFormatDate"
	stmt, err := db.Prepare(sqlCommand)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer stmt.Close()
	rows, err := stmt.Query(startDatetime, endDatetime)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer rows.Close()
	noticeInfos := make([]string, *endIndex)
	for key, _ := range noticeInfos {
		noticeInfos[key] = strconv.FormatBool(false)
	}
	noticeDateTimes := make([]time.Time, 0)
	errorIndex++
	for rows.Next() {
		var createFormatDate time.Time
		err := rows.Scan(&createFormatDate)
		if err != nil {
			fmt.Println("断点" + strconv.Itoa(errorIndex))
			fmt.Println(err)
			return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
		}
		noticeDateTimes = append(noticeDateTimes, createFormatDate)
	}
	for _, value := range noticeDateTimes {
		offsetDay := value.YearDay() - startDatetime.YearDay()
		noticeInfos[offsetDay] = strconv.FormatBool(true)
	}
	result := "[" + strings.Join(noticeInfos, ",") + "]"
	return result, nil

}

func GetTagInfoSql(selectDatetime time.Time) (string, error) {
	errorIndex := 0
	db, err := sql.Open(baseSql.DriverName, baseSql.Dsn)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer db.Close()
	sqlCommand := "SELECT cid,tag,content,createDate From " + tableName + " WHERE isDiscard = 0 AND createFormatDate = ? ORDER BY createDate DESC"
	stmt, err := db.Prepare(sqlCommand)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer stmt.Close()
	rows, err := stmt.Query(selectDatetime)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer rows.Close()
	tagInfos := make([]string, 0)
	errorIndex++
	for rows.Next() {
		var tagInfo TagInfo
		err := rows.Scan(&tagInfo.Cid, &tagInfo.Tag, &tagInfo.Content, &tagInfo.CreateDate)
		if err != nil {
			fmt.Println("断点" + strconv.Itoa(errorIndex))
			fmt.Println(err)
			return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
		}
		tagInfoJson, err := json.Marshal(tagInfo)
		if err != nil {
			fmt.Println("断点" + strconv.Itoa(errorIndex))
			fmt.Println(err)
			return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
		}
		tagInfos = append(tagInfos, string(tagInfoJson))
	}
	err = rows.Err()
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	result := "[" + strings.Join(tagInfos, ",") + "]"
	return result, nil
}

func AddNewTagSql(tag *string, content *string, selectDatetime time.Time) (string, error) {
	errorIndex := 0
	db, err := sql.Open(baseSql.DriverName, baseSql.Dsn)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "新增失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer db.Close()
	cid, _ := uuid.NewV4()
	createDate := time.Now()
	// createFormatDate := time.Date(createDate.Year(), createDate.Month(), createDate.Day(), 0, 0, 0, 0, createDate.Location())
	createFormatDate := selectDatetime
	sqlCommand := "INSERT INTO `" + tableName +
		"`(`cid`,`tag`,`content`,`createDate`,`createFormatDate`)" +
		" VALUES(?, ?, ?, ?, ?)"
	stmt, err := db.Prepare(sqlCommand)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "获取失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer stmt.Close()
	rows, err := stmt.Query(cid.String(), tag, content, createDate, createFormatDate)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "新增失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer rows.Close()
	return "新增成功", nil
}

func DiscardOldTagSql(cid *string) (string, error) {
	errorIndex := 0
	db, err := sql.Open(baseSql.DriverName, baseSql.Dsn)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "废除便签失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer db.Close()
	sqlCommand := "UPDATE " + tableName + " SET isDiscard = ? WHERE cid = ?"
	stmt, err := db.Prepare(sqlCommand)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "废除便签失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer stmt.Close()
	rows, err := stmt.Query(1, cid)
	errorIndex++
	if err != nil {
		fmt.Println("断点" + strconv.Itoa(errorIndex))
		fmt.Println(err)
		return "废除便签失败", fmt.Errorf("error" + strconv.Itoa(errorIndex))
	}
	defer rows.Close()
	return "废除便签成功", nil
}
