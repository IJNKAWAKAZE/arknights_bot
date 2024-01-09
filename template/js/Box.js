window.onload = function(){
    const items = document.getElementsByClassName("rarity");
    for (let i = 0; i < items.length; i++) {
        var item = items[i]
        var data = item.getAttribute("data")
        if(data == 4) {
            item.style.width = "34px";
        }
        if(data == 3) {
            item.style.width = "30px";
        }
        if(data == 2) {
            item.style.width = "25px";
        }
        if(data == 1) {
            item.style.width = "20px";
        }
        if(data == 0) {
            item.style.width = "15px";
        }
    }
}