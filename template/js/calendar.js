var doc = window.document;
var Calendar = function(options) {
    "use strict";
    //默认参数
    var defaults = {
        //中文格式内容
        monthNames: ['一月', '二月', '三月', '四月', '五月', '六月', '七月', '八月' , '九月' , '十月', '十一月', '十二月'],
        dayNames : ['周日', '周一', '周二', '周三', '周四', '周五', '周六'],
        dayLongNames : ['星期日', '星期一', '星期二', '星期三', '星期四', '星期五', '星期六'],
        holiday : {},
        firstDay: 1,                // 从周一开始,计算
        weekendDays: [0, 6],        // 休息日为：周六, 周日
        dateFormat: 'yyyy-mm-dd',   // 打印格式, formatDate 对应
        weekHandler : "dayThead",   // 星期title内容所在类
        monthContainer : "dateUl",  // 日期内容所在容器类
    };
    //参数调整
    options = options || {};
    for (var prop in defaults) {
        if (typeof options[prop] === 'undefined') {
            options[prop] = defaults[prop];
        }
    }
    this.attrs = options;
    //初始化内容
    this.init();
};
Calendar.prototype = {
    constructor : Calendar,
    init : function() {
        //初始布局页面
        this.layout();
        //节点操作
        this.timeNowEle = doc.querySelector(".timeNow");
        //初始化 aside页面
        this.setAside();
    },
    /**
     * 设置aside侧边栏内容
     * @param date 指定日期
     */
    setAside : function(date) {
        var resource = document.getElementById("resource")
        var chip = document.getElementById("chip")
        date = date || new Date();
        var year = date.getFullYear(),
            month = date.getMonth(),
            day = date.getDate();
        var week = date.getDay();
        this.timeNowEle.innerHTML =   "<span>" + year + "年" + (month + 1) + "月" + day + "日" + "</span>"
                                    + "<span>" + this.attrs.dayLongNames[week] + "</span>";
        if(week === 0) {
            resource.textContent = "经验书、技能书、钱、h红票、";
            chip.textContent = "近卫、特种、医疗、重装、辅助、先锋";
        }
        if(week === 1) {
            resource.textContent = "经验书、红票、碳";
            chip.textContent = "术士、狙击、医疗、重装";
        }
        if(week === 2) {
            resource.textContent = "经验书、技能书、钱";
            chip.textContent = "术士、狙击、近卫、特种";
        }
        if(week === 3) {
            resource.textContent = "经验书、技能书、碳";
            chip.textContent = "近卫、特种、辅助、先锋";
        }
        if(week === 4) {
            resource.textContent = "经验书、钱、红票";
            chip.textContent = "医疗、重装、辅助、先锋";
        }
        if(week === 5) {
            resource.textContent = "经验书、技能书、碳";
            chip.textContent = "术士、狙击、医疗、重装";
        }
        if(week === 6) {
            resource.textContent = "经验书、钱、红票";
            chip.textContent = "术士、狙击、近卫、特种、辅助、先锋";
        }
    },
    /**
     * 布局函数
     */
    layout : function() {
        var layoutDate = this.value ? this.value: new Date().setHours(0,0,0,0);
        this.value = layoutDate;

        //三个月的 HTML信息
        var prevMonthHTML = this.monthHTML(layoutDate, 'prev');
        var currentMonthHTML = this.monthHTML(layoutDate);
        var nextMonthHTML = this.monthHTML(layoutDate, 'next');
        var monthHTML =   '<div class="dateLi prev-month-html"><div class="dayTbody">' + prevMonthHTML + "</div></div>"
                        + '<div class="dateLi current-month-html"><div class="dayTbody">' + currentMonthHTML + "</div></div>"
                        + '<div class="dateLi next-month-html"><div class="dayTbody">' + nextMonthHTML + "</div></div>"

        if(!this._initEd) {
            //初次渲染的时候使用, 否则不使用
            //渲染 星期头部
            var weekHeaderHTML = [];
            for(var i = 0; i < 7; i++) {
                var weekDayIndex = (i + this.attrs.firstDay > 6) ? (i - 7 + this.attrs.firstDay) : (i + this.attrs.firstDay);
                var dayName = this.attrs.dayNames[weekDayIndex];
                if(this.attrs.weekendDays.indexOf(weekDayIndex) !== -1) {
                    //休息日样式
                    weekHeaderHTML.push('<div class="dayTd active">' + dayName + '</div>');
                }else {
                    weekHeaderHTML.push('<div class="dayTd">' + dayName + '</div>');
                }
            }

            var yearSelectHTML = [],
                monthSelectHTML = [];
            for(var i = 1900; i <= 2050; i++) {
                var str = "<li class='list-year' data-year='" + i + "'>" + i + "年</li>";
                yearSelectHTML.push(str);
            }
            for(var i = 1; i <= 12; i++) {
                var str = "<li class='list-month' data-month='" + (i - 1) + "'>" + i + "月</li>";
                monthSelectHTML.push(str);
            }

            doc.querySelector("." + this.attrs.weekHandler).innerHTML = weekHeaderHTML.join("");
        }
        doc.querySelector("." + this.attrs.monthContainer).innerHTML = monthHTML;
    },
    /**
     * 获取当月总天数
     * @param date 日期
     */
    totalDaysInMonth : function(date) {
        var d = new Date(date);
        return new Date(d.getFullYear(), d.getMonth() + 1, 0).getDate();
    },
    /**
     * 日期HTML设置
     * @param date   设置时间
     * @param offset 设置偏移 ["next" -> "下个月"  "prev" -> "上个月"]
     * @returns {string} 返回拼接好的HTML字符串
     */
    monthHTML : function(date, offset) {
        var date = new Date(date);
        var year = date.getFullYear(),
            month = date.getMonth(),
            day = date.getDate();

        //下个月
        if (offset === 'next') {
            if (month === 11) {
                date = new Date(year + 1, 0);
            }else {
                date = new Date(year, month + 1, 1);
            }
        }
        //上个月
        if (offset === 'prev') {
            if (month === 0) {
                date = new Date(year - 1, 11);
            }else {
                date = new Date(year, month - 1, 1);
            }
        }
        //调整时间
        if (offset === 'next' || offset === 'prev') {
            month = date.getMonth();
            year = date.getFullYear();
        }

        //上月 | 本月 总天数, 本月第一天索引
        var daysInPrevMonth = this.totalDaysInMonth(new Date(date.getFullYear(), date.getMonth()).getTime() - 10 * 24 * 60 * 60 * 1000),
            daysInMonth = this.totalDaysInMonth(date),
            firstDayOfMonthIndex = new Date(date.getFullYear(), date.getMonth()).getDay();

        if (firstDayOfMonthIndex === 0) {
            firstDayOfMonthIndex = 7;
        }

        var dayDate,
            rows = 6,
            cols = 7,
            monthHTML = [],
            dayIndex = 0 + (this.attrs.firstDay - 1),
            today = new Date().setHours(0,0,0,0);

        for(var i = 1; i <= rows; i++) {
            var row = i;
            var tdList = [];
            for(var j = 1; j <= cols; j++) {
                var col = j;
                dayIndex ++;

                var dayNumber = dayIndex - firstDayOfMonthIndex;
                //要添加的类名, 默认dayTd
                var classNames = ["dayTd"];
                if(dayNumber < 0) {
                    //上个月日期
                    classNames.push("date-prev");
                    dayNumber = daysInPrevMonth + dayNumber + 1;
                    dayDate = new Date(month - 1 < 0 ? year - 1 : year, month - 1 < 0 ? 11 : month - 1, dayNumber).getTime();
                }else {
                    dayNumber = dayNumber + 1;
                    if (dayNumber > daysInMonth) {
                        //下个月日期
                        classNames.push("date-next");
                        dayNumber = dayNumber - daysInMonth;
                        dayDate = new Date(month + 1 > 11 ? year + 1 : year, month + 1 > 11 ? 0 : month + 1, dayNumber).getTime();
                    }else {
                        dayDate = new Date(year, month, dayNumber).getTime();
                    }
                }
                //今天 date-current
                if (dayDate === today) {
                    classNames.push("date-current");
                }
                //周六日 date-reset
                if(dayIndex % 7 === this.attrs.weekendDays[0] || dayIndex % 7 === this.attrs.weekendDays[1]) {
                    classNames.push("date-reset");
                }
                dayDate = new Date(dayDate);

                var dayYear = dayDate.getFullYear();
                var dayMonth = dayDate.getMonth() + 1;
                if(dayMonth < 10) {
                    dayMonth = "0" + dayMonth;
                }
                if(dayNumber < 10) {
                    dayNumber = "0" + dayNumber
                }
                var holiday = this.attrs.holiday[dayYear + "-" + dayMonth + "-" + dayNumber];
                //活动显示
                if(holiday) {
                    var alamanac = holiday;
                    classNames.push("date-holiday");
                }else {
                    var alamanac = "";
                }
                tdList.push(
                      '<div class="' + classNames.join(" ") + '" data-year=' + dayYear + ' data-month=' + dayMonth + ' data-day=' + dayNumber + '>'
                    +   '<span class="dayNumber">' + dayNumber + "</span>"
                    +   '<span class="almanac"><ul>' + alamanac + "</ul></span>"
                    + '</div>'
                )
            }
            monthHTML.push(
                '<div class="dayTr">' + tdList.join("") + '</div>'
            )
        }
        return monthHTML.join("");
    },
    /**
     * 格式化日期内容
     * @param date
     * @returns {string|XML}
     */
    formatDate : function(date) {
        date = new Date(date);
        var year = date.getFullYear();
        var month = date.getMonth();
        var month1 = month + 1;
        var day = date.getDate();
        var weekDay = date.getDay();
        return this.attrs.dateFormat
            .replace(/yyyy/g, year)
            .replace(/yy/g, (year + '').substring(2))
            .replace(/mm/g, month1 < 10 ? '0' + month1 : month1)
            .replace(/m/g, month1)
            .replace(/MM/g, this.attrs.monthNames[month])
            .replace(/dd/g, day < 10 ? '0' + day : day)
            .replace(/d/g, day)
            .replace(/DD/g, this.attrs.dayNames[weekDay])
    }
};
var Util = {
    /**
     * css前缀判断，仅执行一次
     */
    prefix : (function() {
        var div = document.createElement('div');
        var cssText = '-webkit-transition:all .1s; -moz-transition:all .1s; -o-transition:all .1s; -ms-transition:all .1s; transition:all .1s;';
        div.style.cssText = cssText;
        var style = div.style;
        if(style.transition) {
            return '';
        }
        if (style.webkitTransition) {
            return '-webkit-';
        }
        if (style.MozTransition) {
            return '-moz-';
        }
        if (style.oTransition) {
            return '-o-';
        }
        if (style.msTransition) {
            return '-ms-';
        }
    })()
};
