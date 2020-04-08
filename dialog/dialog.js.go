package dialog

const _dialogJs = `
dialog = (function() {
    const File = function(config) {
        return _dialogMux('file', config)
    }

    return {
        File
    }
})()
`
