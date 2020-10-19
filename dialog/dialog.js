dialog = (function() {
    const alert = function(message, title = 'Alert') {
        return _dialogMux('message', {type: 'info', message, title})
    }
    const error = function(message, title = 'Error') {
        return _dialogMux('message', {type: 'error', message, title})
    }
    const confirm = function(message, title = 'Confirm?') {
        return _dialogMux('message', {type: 'confirm', message, title})
    }
    const directory = function(config) {
        return _dialogMux('directory', config)
    }
    const file = function(config) {
        return _dialogMux('file', config)
    }

    return {
        alert,
        error,
        confirm,
        directory,
        file
    }
})()
