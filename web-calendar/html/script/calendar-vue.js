var weekText = [
    {name: '日', isActive: false},
    {name: '一', isActive: false},
    {name: '二', isActive: false},
    {name: '三', isActive: false},
    {name: '四', isActive: false},
    {name: '五', isActive: false},
    {name: '六', isActive: false}
];
var vue = new Vue({
    el: '#pc',
    data: {
        //calendar
        selectYear: 1,
        selectMonth: 1,
        selectDate: new Date(0),
        nowServerDate: new Date(0),
        weekText: weekText,
        showCalendar: [],
        //note
        isEdit: false,
        addTag: "All",
        addContent: "",
        centerModalContentData: "",
        showTagInfos: [],
    },
    methods: {
        setSelectDate(date) {
            if (date !== '') {
                var selectDate = new Date(this.selectDate);
                selectDate.setDate(date);
                this.selectDate = selectDate;
            }
        },
        reverseIsEdit() {
            this.isEdit = !this.isEdit;
        },
        GetCalendarNoticeInfo(startDate, endIndex) {
            this.$http.get(
                    '/api/calendar/controller',
                    {params: {action: "GetCalendarNoticeInfo", startDate: startDate.toUTCString(), endIndex: endIndex}}
            ).then(function (res) {
                if (res.status === 200) {
                    calendarAreaSetNoticeInfo(res.body);
                }
            }, function () {
                //todo no deal
            });
        },
        GetTagInfo() {
            this.$http.get(
                    '/api/calendar/controller',
                    {params: {action: "GetTagInfo", selectDate: this.selectDate.toUTCString()}}
            ).then(function (res) {
                if (res.status === 200) {
                    setTagInfo(res.body);
                }
            }, function () {
                //todo no deal
            });
        },
        AddNewTag() {
            if (this.addContent === '') {
                vue.centerModalContent = "请输入记录内容";
                return;
            }
            this.$http.post(
                    '/api/calendar/controller',
                    {action: "AddNewTag", tag: this.addTag, content: this.addContent, selectDate: this.selectDate.toUTCString()},
                    {emulateJSON: true}
            ).then(function (res) {//200
                vue.addTag = "All";
                vue.addContent = "";
                vue.GetTagInfo();
                vue.centerModalContent = res.body;
            }, function (res) {
                vue.centerModalContent = res.body;
            });
        },
        DiscardOldTag(cid) {
            this.$http.post(
                    '/api/calendar/controller',
                    {action: "DiscardOldTag", cid: cid},
                    {emulateJSON: true}
            ).then(function (res) {
                vue.GetTagInfo();
                vue.centerModalContent = res.body;
            }, function (res) {
                vue.centerModalContent = res.body;
            });
        },
    },
    computed: {
        centerModalContent: {
            get() {
                return this.centerModalContentData;
            },
            set(nval) {
                if (nval === '')
                    nval = "无效操作";
                this.centerModalContentData = nval;
                $('#exampleModalCenter').modal('show');
            }
        }
    },
});

function calendarAreaSetNoticeInfo(noticeInfo) {
    if (noticeInfo instanceof Array) {
        for (var noticeIndex = 0; noticeIndex < noticeInfo.length; noticeIndex++) {
            var oneNotice = noticeInfo[noticeIndex];
            if (typeof oneNotice !== 'boolean') {
                return;
            }
        }
        var noticeIndex = 0;
        for (var calendarIndex = 0; calendarIndex < vue.$data.showCalendar.length; calendarIndex++) {
            var areaCalendar = vue.$data.showCalendar[calendarIndex];
            for (var areaIndex = 0; areaIndex < areaCalendar.length; areaIndex++) {
                var calendar = areaCalendar[areaIndex];
                if (!calendar.isEmpty) {
                    vue.$data.showCalendar[calendarIndex][areaIndex].haveNotice = noticeInfo[noticeIndex];
                    noticeIndex++;
                }
            }
        }
        vue.$data.showCalendar.push({});
        vue.$data.showCalendar.pop();
    }
}

function setTagInfo(tagInfo) {
    vue.$data.showTagInfos.splice(0, vue.$data.showTagInfos.length);
    if (tagInfo instanceof Array) {
        for (var tagIndex = 0; tagIndex < tagInfo.length; tagIndex++) {
            var oneTag = tagInfo[tagIndex];
            if (typeof oneTag === 'object') {
                vue.$data.showTagInfos.push(oneTag);
            }
        }
    }
}

vue.$watch('selectYear', function (nval) {
    if (typeof nval !== 'number' || nval < 1)
        return;
    if (nval !== vue.$data.selectDate.getFullYear()) {
        var selectDate = new Date(vue.$data.selectDate);
        selectDate.setFullYear(nval);
        vue.$data.selectDate = selectDate;
    }
});
vue.$watch('selectMonth', function (nval) {
    if (typeof nval !== 'number')
        return;
    if (nval < 1 || nval > 12)
        return;
    nval--;
    if (nval !== vue.$data.selectDate.getMonth()) {
        var selectDate = new Date(vue.$data.selectDate);
        selectDate.setMonth(nval);
        vue.$data.selectDate = selectDate;
    }
});
vue.$watch('selectDate', function (nval) {
    vue.$data.addTag = "All";
    vue.$data.addContent = "";
    vue.GetTagInfo();
    var newNval = new Date(nval);
    if (vue.$data.selectYear !== newNval.getFullYear()) {
        vue.$data.selectYear = newNval.getFullYear();
    }
    var currentMonth = newNval.getMonth() + 1;
    if (vue.$data.selectMonth !== currentMonth) {
        vue.$data.selectMonth = currentMonth;
    }
    vue.$data.showCalendar.splice(0, vue.$data.showCalendar.length);
    var area = [];
    for (var w = 0; w < weekText.length; w++) {
        if (w === nval.getDay()) {
            vue.$data.weekText[w].isActive = true;
        } else {
            vue.$data.weekText[w].isActive = false;
        }
    }
    var nowMonthOneDate = new Date(nval);
    nowMonthOneDate.setDate(1);
    for (var i = 0; i < nowMonthOneDate.getDay(); i++) {
        area.push({isEmpty: true, date: "", comment: "", isActive: false});
        addAreaCalendar(area);
    }
    var thisMonthAllDay = getDaysInMonth(nowMonthOneDate.getFullYear(), nowMonthOneDate.getMonth());
    for (var j = 1; j <= thisMonthAllDay; j++) {
        var comment = "";
        if (j === 22)
            comment = "";
        area.push({isEmpty: false, date: j, comment: comment, isActive: j === nval.getDate()});
        addAreaCalendar(area);
    }
    if (area.length < 7) {
        var indexLast = 7 - area.length;
        for (var k = 0; k < indexLast; k++) {
            area.push({isEmpty: true, date: "", comment: "", isActive: false});
            addAreaCalendar(area);
        }
    }
    vue.GetCalendarNoticeInfo(nowMonthOneDate, thisMonthAllDay);
});

function addAreaCalendar(area) {
    if (area.length === 7) {
        vue.$data.showCalendar.push(area.slice(0));
        area.splice(0, area.length);
    }
}

function updateCalendar() {
    var nowClientDate = new Date();
    nowClientDate.setHours(0, 0, 0, 0);
    if (nowClientDate.getTime() !== vue.$data.nowServerDate.getTime()) {
        vue.$data.nowServerDate = new Date(nowClientDate);
        vue.$data.selectDate = new Date(nowClientDate);
    }
}
//定时调用更新
setInterval(function () {
    updateCalendar();
}, 10000);
updateCalendar();
//计算某年某月有多少天
function isLeapYear(a) {
    return a % 400 === 0 || a % 4 === 0 && a % 100 !== 0;
}
function getDaysInMonth(a, b) {
    return [31, null, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31][b] || (isLeapYear(a) ? 29 : 28);
}
//https://www.cnblogs.com/miumiu316/p/7210731.html