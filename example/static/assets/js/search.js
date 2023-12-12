$(function(){
    $('.sc-searchbox>input').keypress(function(ev){
        if (ev.which !== 13) {
            return
        }
        search()
    }) 
    
    function search() {
        var content = $('.sc-searchbox>input').val();
        alert(content);
    }
})