window.onload = function(){
    const complete = document.getElementById("complete");
    var completeTime = complete.textContent;
    var date = convertTimestamp(completeTime)
    complete.textContent = date.hour + "小时" + date.minute + "分后恢复";

    const lowerItemTermTime = document.getElementById("lower_item_term_time");
    const higherItemTermTime = document.getElementById("higher_item_term_time");
    var termTime = lowerItemTermTime.textContent;
    var d = convertTimestamp(termTime)
    lowerItemTermTime.textContent = d.day + "天"
    higherItemTermTime.textContent = d.day + "天"

    const remainSecs = document.getElementById("remain_secs");
    var hours = parseInt(remainSecs.textContent / 3600)
    var minutes = parseInt((remainSecs.textContent % 3600) / 60)
    remainSecs.textContent = hours + ":" + minutes + "h"
}

function convertTimestamp(timestamp) {
    var day, hour, minute;
    minute = Math.floor(timestamp / 60000);
    hour = Math.floor(minute / 60);
    minute = minute % 60;
    day = Math.floor(hour / 24);
    hour = hour % 24;
    return {
        day: day,
        hour: hour,
        minute: minute
    };
}