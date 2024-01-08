window.onload = function(){
    const complete = document.getElementById("complete");
    var completeTime = complete.textContent;
    if(!isNaN(Number(completeTime,10))){
        var date = countDown(completeTime);
        complete.textContent = date.hour + '时' + date.minute + '分后恢复';
    }

    const lowerItemTermTime = document.getElementById("lower_item_term_time");
    const higherItemTermTime = document.getElementById("higher_item_term_time");
    var termTime = lowerItemTermTime.textContent;
    var d = countDown(termTime);
    lowerItemTermTime.textContent = d.day + '天';
    higherItemTermTime.textContent = d.day + '天';

    const remainSecs = document.getElementById("remain_secs");
    remainSecs.textContent = formatTime(remainSecs.textContent)
}

function countDown(time) {
    var nowTime = parseInt(new Date().getTime() / 1000);
    var inputTime  = time;
    var times = (inputTime - nowTime);
    var d = parseInt(times / 60 / 60 / 24);
    d = d < 10 ? '0' + d : d;
    var h = parseInt(times / 60 / 60 % 24);
    h = h < 10 ? '0' + h : h;
    var m = parseInt(times / 60 % 60);
    m = m < 10 ? '0' + m : m;
    var s = parseInt(times % 60);
    s = s < 10 ? '0' + s : s;
    return {
        day: d,
        hour: h,
        minute: m
    };
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