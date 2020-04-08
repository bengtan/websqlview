dialog = (function() {
    const Directory = function(config) {
        return _dialogMux('directory', config)
    }
    const File = function(config) {
        return _dialogMux('file', config)
    }

    return {
        Directory,
        File
    }
})()
