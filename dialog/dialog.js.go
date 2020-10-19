package dialog

const _dialogJs = `
dialog = (function() {
    const directory = function(config) {
        return _dialogMux('directory', config)
    }
    const file = function(config) {
        return _dialogMux('file', config)
    }

    return {
        directory,
        file
    }
})()
`
