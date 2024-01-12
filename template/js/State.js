window.onload = function(){
    const complete = document.getElementById("complete");
    var completeTime = complete.textContent;
    if(!isNaN(Number(completeTime,10))){
        var date = completeRecoveryTime(completeTime);
        complete.textContent = date.hour + '时' + date.minute + '分后恢复';
    }

    const lowerItemTermTime = document.getElementById("lower_item_term_time");
    const higherItemTermTime = document.getElementById("higher_item_term_time");
    var termTime = lowerItemTermTime.textContent;
    var d = recoverTime(termTime);
    lowerItemTermTime.textContent = d + '天';
    higherItemTermTime.textContent = d + '天';

    const campaignRecoverTime = document.getElementById("campaign_recover_time_item");
    var campaignRecover = recoverTime(campaignRecoverTime.textContent);
    campaignRecoverTime.textContent = campaignRecover + '天';

    const remainSecs = document.getElementById("remain_secs_item");
    remainSecs.textContent = formatTime(remainSecs.textContent)
}

function completeRecoveryTime(time) {
    var nowTime = parseInt(new Date().getTime() / 1000);
    var times = (time - nowTime);
    var d = parseInt(times / 60 / 60 / 24);
    var h = parseInt(times / 60 / 60 % 24);
    var m = parseInt(times / 60 % 60);
    var s = parseInt(times % 60);
    return {
        day: d,
        hour: h,
        minute: m
    };
}

function recoverTime(time) {
    var nowTime = parseInt(new Date().getTime() / 1000);
    var times = (time - nowTime);
    var d = Math.ceil(times / 60 / 60 / 24);
    return d;
}

const formatTime = (seconds)=>{
    let result = [];
    let count = 2;
    while(count >= 0) {
        let current = Math.floor(seconds / (60 ** count));
        result.push(current);
        seconds -= current * (60 ** count);
        --count;
    }
    return result.map(item => item <= 9 ? `0${item}` : item).join(":");
}