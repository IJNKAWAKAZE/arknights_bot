// 公共工具函数，供各页面脚本复用
function ready(fn) {
    if (document.readyState !== 'loading') {
        fn();
    } else {
        document.addEventListener('DOMContentLoaded', fn);
    }
}

function pad(n) {
    return n < 10 ? '0' + n : '' + n;
}

// 时间戳(秒)转 "YYYY-MM-DD HH:mm:ss"；withTime 传 false 时只到日期
function timestampToTime(timestamp, withTime) {
    if (!timestamp) return '';
    var date = new Date(parseInt(timestamp, 10) * 1000);
    if (isNaN(date.getTime())) return '';
    var result = date.getFullYear() + '-' + pad(date.getMonth() + 1) + '-' + pad(date.getDate());
    if (withTime === false) return result;
    return result + ' ' + pad(date.getHours()) + ':' + pad(date.getMinutes()) + ':' + pad(date.getSeconds());
}

// 距离某日期已过去的天数
function getWakeDays(time) {
    var dateBegin = new Date(time);
    if (isNaN(dateBegin.getTime())) return 0;
    return Math.floor((Date.now() - dateBegin.getTime()) / (24 * 3600 * 1000));
}

// 距离未来时间戳(秒)的天数，向上取整
function daysUntil(ts) {
    var times = parseInt(ts, 10) - Math.floor(Date.now() / 1000);
    if (isNaN(times)) return 0;
    return Math.ceil(times / 86400);
}

// 秒数转 "HH:MM:SS"
function formatTime(seconds) {
    seconds = parseInt(seconds, 10) || 0;
    return pad(Math.floor(seconds / 3600)) + ':' + pad(Math.floor(seconds % 3600 / 60)) + ':' + pad(seconds % 60);
}
