os = (function() {
    return {
        exit: exitCode => {
            return _osMux('exit', exitCode)
        },
    }
})()
