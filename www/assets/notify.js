var notify = (function(){
    function success(message) {
        alert(message)
    }
    function error(message) {
        alert(message)
    }
    return { success, error }
}())
