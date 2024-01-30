window.onload = function(){
    const regTime = document.getElementById("regTime");
    regTime.textContent = timestampToTime(regTime.textContent)

    const wakeDays = document.getElementById("wakeDays");
    wakeDays.textContent = '已苏醒' + getWakeDays(regTime.textContent) + '天';
}

function timestampToTime(timestamp) {
    timestamp = timestamp ? timestamp : null;
    let date = new Date(timestamp*1000);
    let Y = date.getFullYear() + '-';
    let M = (date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1) + '-';
    let D = (date.getDate() < 10 ? '0' + date.getDate() : date.getDate()) + ' ';
    let h = (date.getHours() < 10 ? '0' + date.getHours() : date.getHours()) + ':';
    let m = (date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()) + ':';
    let s = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds();
    return Y + M + D + h + m + s;
}

function getWakeDays(time) {
    var dateBegin = new Date(time);
    var date = new Date();
    var result = date.getTime() - dateBegin.getTime();
    return Math.floor(result / (24 * 3600 * 1000));
}