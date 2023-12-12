$(function(){
    (function(){
        setInterval(() => {
            loadClock($('#footer-timer'), $('#footer-timer-ts'),$('#footer-timer-year'),$('#footer-timer-month'),$('#footer-timer-day'))
        }, 1000);
    })()

    $('#footer-timer').click(function(e){
        copy2Clipboard(e.target);
        addToast('message','已复制到剪贴板')
    })

    $('#footer-timer-ts').click(function(e){
        copy2Clipboard(e.target);
        addToast('message','已复制到剪贴板')
    })

    function copy2Clipboard(e) {
        navigator.clipboard.writeText($(e).html());
    }

    function loadClock(e, f, y, m, d) {
        var curTime = new Date();

        var curTs = parseInt(curTime.getTime() / 1000);

        var curYear = curTime.getFullYear();
        var curMonth = curTime.getMonth() + 1;
        var curDay = curTime.getDate();
        if (curMonth < 10) {
            curMonth = "0" + curMonth;
        }
        if (curDay < 10) {
            curDay = "0" + curDay;
        }

        var curHour = curTime.getHours();
        var curMinute = curTime.getMinutes();
        var curSecond = curTime.getSeconds();
        if (curMinute < 10) {
            curMinute = "0" + curMinute;
        }
        if (curSecond < 10) {
            curSecond = "0" + curSecond;
        }
        var display = curHour + ":" + curMinute + ":" + curSecond
        $(e).html(display)
        $(f).html(curTs)
        $(y).html(curYear)
        $(m).html(curMonth)
        $(d).html(curDay)
    }
})