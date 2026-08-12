ready(function () {
    var complete = document.getElementById("complete");
    if (complete && !isNaN(Number(complete.textContent))) {
        var date = completeRecoveryTime(complete.textContent);
        complete.textContent = date.hour + '时' + date.minute + '分后恢复';
    }

    var lowerItemTermTime = document.getElementById("lower_item_term_time");
    var higherItemTermTime = document.getElementById("higher_item_term_time");
    if (lowerItemTermTime && higherItemTermTime) {
        var d = daysUntil(lowerItemTermTime.textContent);
        lowerItemTermTime.textContent = d + '天';
        higherItemTermTime.textContent = d + '天';
    }

    var campaignRecoverTime = document.getElementById("campaign_recover_time_item");
    if (campaignRecoverTime) {
        campaignRecoverTime.textContent = daysUntil(campaignRecoverTime.textContent) + '天';
    }

    var remainSecs = document.getElementById("remain_secs_item");
    if (remainSecs) {
        remainSecs.textContent = formatTime(remainSecs.textContent);
    }
});

function completeRecoveryTime(time) {
    var times = parseInt(time, 10) - Math.floor(Date.now() / 1000);
    return {
        day: Math.floor(times / 86400),
        hour: Math.floor(times / 3600 % 24),
        minute: Math.floor(times / 60 % 60)
    };
}
