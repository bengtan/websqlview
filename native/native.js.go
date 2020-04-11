package native

const _nativeJs = `
native = (function() {
    return {
        exit: exitCode => {
            return _nativeMux('exit', exitCode)
        },
        remove: name => {
            return _nativeMux('remove', name)
        },
    }
})()
`
