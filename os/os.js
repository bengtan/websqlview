os = (function() {
    return {
        exit: exitCode => {
            return _osMux('exit', exitCode)
        },
        remove: name => {
            return _osMux('remove', name)
        },
    }
})()
