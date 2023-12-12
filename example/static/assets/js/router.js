$(function(){
    $('.nav-item i').click(function(e){
        router = $(e.target).parent('span').parent('div').data('router')
        if (router === undefined) {
            return false
        }
        window.location = '/' + router
    });

    (function() {
        initNavRouter();
    })()

    function initNavRouter() {
        let router = window.location.pathname;
        $('.nav-item').removeClass('nav-item-active')
        $('.nav-item').each((index,ele) => {
            if ('/' + $(ele).data('router') === router) {
                $(ele).addClass('nav-item-active')
            }
        });
    }
})