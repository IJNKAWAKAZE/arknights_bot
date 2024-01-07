window.onload = function(){
    const complete = document.getElementById("complete");
    var completeTime = complete.textContent;
    var date = convertTimestamp(completeTime)
    complete.textContent = date.hour + "小时" + date.minute + "分后恢复";
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