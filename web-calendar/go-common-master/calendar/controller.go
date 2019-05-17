// controller
package calendar

import (
	"net/http"
	"strconv"
	"time"

	"fmt"

	net "../net"
)

func Start(rw http.ResponseWriter, req *http.Request) {
	//此处可并发记录日志
	req.ParseForm()
	reqMethod := req.Method
	if reqMethod == "GET" {
		// GetCurrentDate(rw, req)
		GetCalendarNoticeInfo(rw, req)
		GetTagInfo(rw, req)
	}
	if reqMethod == "POST" {
		AddNewTag(rw, req)
		DiscardOldTag(rw, req)
	}
}

func GetCurrentDate(rw http.ResponseWriter, req *http.Request) {
	reqContentType := req.Header.Get("Content-Type")
	action := req.Form["action"][0]
	if action == "GetCurrentDate" && reqContentType == net.JsonType {
		result := net.ResponseJson{strconv.Itoa(net.HttpNomalSuccessStatus), "当前Web时间", time.Now().String()}
		net.ReturnResponse(req, rw, result)
	}
}

func GetCalendarNoticeInfo(rw http.ResponseWriter, req *http.Request) {
	action := req.Form["action"][0]
	if action == "GetCalendarNoticeInfo" {
		startDate := req.Form["startDate"]
		startDatetime, err := time.ParseInLocation(time.RFC1123, startDate[0], time.UTC)
		endIndex, errI := strconv.Atoi(req.Form["endIndex"][0])
		if err == nil && !startDatetime.IsZero() && errI == nil && endIndex >= 0 {
			msg, error := GetCalendarNoticeInfoSql(startDatetime, &endIndex)
			if error != nil {
				rw.WriteHeader(net.HttpNomalErrorStatus)
			}
			fmt.Fprint(rw, msg)
		}
	}
}

func GetTagInfo(rw http.ResponseWriter, req *http.Request) {
	action := req.Form["action"][0]
	if action == "GetTagInfo" {
		selectDate := req.Form["selectDate"]
		selectDatetime, err := time.ParseInLocation(time.RFC1123, selectDate[0], time.UTC)
		if err == nil && !selectDatetime.IsZero() {
			msg, error := GetTagInfoSql(selectDatetime)
			if error != nil {
				rw.WriteHeader(net.HttpNomalErrorStatus)
			}
			fmt.Fprint(rw, msg)
		}
	}
}

func AddNewTag(rw http.ResponseWriter, req *http.Request) {
	action := req.PostFormValue("action")
	if action == "AddNewTag" {
		tag := req.PostFormValue("tag")
		content := req.PostFormValue("content")
		tag = string(tag)
		content = string(content)
		if tag == "" {
			tag = "All"
		}
		if content == "" {
			rw.WriteHeader(net.HttpNomalErrorStatus)
			fmt.Fprint(rw, "请输入记录内容")
			return
		}
		selectDate := req.Form["selectDate"]
		selectDatetime, err := time.ParseInLocation(time.RFC1123, selectDate[0], time.UTC)
		if err == nil && !selectDatetime.IsZero() {
			msg, error := AddNewTagSql(&tag, &content, selectDatetime)
			if error != nil {
				rw.WriteHeader(net.HttpNomalErrorStatus)
			}
			fmt.Fprint(rw, msg)
		}
	}
}

func DiscardOldTag(rw http.ResponseWriter, req *http.Request) {
	action := req.PostFormValue("action")
	if action == "DiscardOldTag" {
		cid := req.PostFormValue("cid")
		cid = string(cid)
		if cid == "" {
			rw.WriteHeader(net.HttpNomalErrorStatus)
			fmt.Fprint(rw, "便签无效")
			return
		}
		msg, error := DiscardOldTagSql(&cid)
		if error != nil {
			rw.WriteHeader(net.HttpNomalErrorStatus)
		}
		fmt.Fprint(rw, msg)
	}
}
